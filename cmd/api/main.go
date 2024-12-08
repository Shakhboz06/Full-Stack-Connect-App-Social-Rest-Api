package main

import (
	"expvar"
	"go-project/internal/auth"
	"go-project/internal/db"
	"go-project/internal/env"
	"go-project/internal/mailer"
	"go-project/internal/ratelimiter"
	"go-project/internal/store"
	"go-project/internal/store/cache"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const version = "1.0.0"

//	@title			ConnectApp API
//	@description	This is social media ConnectApp API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	
	cfg := servConfig{
		addr: env.GetString("ADDR", ":3000"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30 ),
			maxIdleTime: env.GetString("DB_MAX_TIME_CONNS", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			mailExp: time.Hour * 24 * 2, //2 days
			sendGrid: sendGrid{
				apikey: env.GetString("SENDGRID_API_KEY", ""),
			},
			fromEmail: env.GetString("FROM_EMAIL", ""),
		},
		auth: authConfig{
			basic: authBasicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRETS", ""),
				exp: time.Hour * 24 * 2,
				iss: "ConnectApp Social",
			},
		},
		redis: redisConfig{
			addr: env.GetString("REDIS_ADDR", "localhost:6379"),
			pass: env.GetString("REDIS_PASS", ""),
			db: env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		
		ratelimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetInt("RATELIMITER_REQUEST_COUNT", 20),
			TimeFrame: time.Second * 5,
			Enabled: env.GetBool("RATELIMITER_REQUEST", true),
		},

	}
	
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()


	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)

	if err != nil{
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Info("Database connection established")
	store := store.NewPostgresStorage(db)
	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apikey, cfg.mail.fromEmail)

	var rdb *redis.Client
	if cfg.redis.enabled{
		rdb = cache.NewRedisClient(cfg.redis.addr,cfg.redis.pass, cfg.redis.db)
		logger.Info("redis cache connection established")

		defer rdb.Close()
	}

	cacheStorage := cache.NewRedisStorage(rdb)

	rateLimiter := ratelimiter.NewfixedWindowRateLimiter(
		cfg.ratelimiter.RequestPerTimeFrame,
		cfg.ratelimiter.TimeFrame,
	)

	JWTAuth := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &application{
		config: cfg,
		store: store,	
		logger: logger,
		mailer: mailer,
		authenticator: JWTAuth, 
		cacheStorage: cacheStorage,
		ratelimiter: rateLimiter,
	}

	expvar.NewString("version").Set(version)	
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("go routines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))


	mux := app.mount()
	logger.Fatal(app.run(mux))
}