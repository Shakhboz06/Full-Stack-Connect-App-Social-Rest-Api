package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-project/internal/auth"
	"go-project/internal/env"
	"go-project/internal/mailer"
	"go-project/internal/ratelimiter"
	"go-project/internal/store"
	"go-project/internal/store/cache"
	"net/http"
	"time"

	"go-project/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config servConfig
	store store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
	authenticator auth.Aunthenticator
	cacheStorage cache.Storage
	ratelimiter ratelimiter.Limiter
}

type servConfig struct {
	addr string
	db dbConfig
	env string
	apiURL string
	mail mailConfig
	frontendURL string
	auth authConfig
	redis redisConfig
	ratelimiter ratelimiter.Config
}

type redisConfig struct{
	addr string
	pass string
	db int
	enabled bool
}
type authConfig struct{
	basic authBasicConfig
	token tokenConfig
}

type tokenConfig struct{
	secret string
	exp time.Duration
	iss string
}

type authBasicConfig struct{
	user string
	pass string
}
type mailConfig struct{
	mailExp time.Duration
	sendGrid sendGrid 
	fromEmail string
}

type sendGrid struct{
	apikey string
}


type dbConfig struct{
	addr string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:3000")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	if app.config.ratelimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}
 
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicMiddlewareAuth()).Get("/health", app.healthCheckHandler)
		r.With(app.BasicMiddlewareAuth()).Get("/metrics", expvar.Handler().ServeHTTP)


		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router){
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPosts)
		
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPosts)
				r.Delete("/", app.RoleBasedAuthMiddleware("moderator", app.deletePost))
				r.Patch("/", app.RoleBasedAuthMiddleware("admin", app.updatePost))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activation/{token}", app.userActivationHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			
			

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
		
		
		//Public Routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.userRegisterHandler)
			r.Post("/token", app.getUserTokenHandler)
		})
	})


	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func ()  {
		
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		app.logger.Infow("signal  caught", "signal", s.String())

		shutdown <-srv.Shutdown(ctx)

	}()

	app.logger.Infow("server is runnning at", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed){
		return err
	}

	err = <-shutdown
	if err != nil{
		return err
	}
	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
