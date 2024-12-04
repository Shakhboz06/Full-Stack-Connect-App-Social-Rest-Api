package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-project/internal/store"

	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicMiddlewareAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.BasicunAuthorizedError(w, r, fmt.Errorf("not Authorized, missing credentials"))
			}

			//parse it
			parse := strings.Split(authHeader, " ")
			if len(parse) != 2 || parse[0] != "Basic" {
				app.BasicunAuthorizedError(w, r, fmt.Errorf("not Authorized, incorrect Credentials"))
				return
			}

			// json decode
			decoded, err := base64.StdEncoding.DecodeString(parse[1])
			if err != nil {
				app.BasicunAuthorizedError(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			password := app.config.auth.basic.pass
			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.BasicunAuthorizedError(w, r, fmt.Errorf("incorrect Username Or Password"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unAuthorizedError(w, r, fmt.Errorf("not Authorized, missing credentials"))
		}

		//parse it
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthorizedError(w, r, fmt.Errorf("not Authorized, malformed header"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)

		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.getUser(ctx, userId)
		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) RoleBasedAuthMiddleware(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserCtx(r)
		post := getPostFromCtx(r)

		fmt.Println(post)
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		//role procedence
		Ispermitted, err := app.checkRoleProcedence(r.Context(), user, role)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		fmt.Println(Ispermitted)
		if !Ispermitted {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRoleProcedence(ctx context.Context, user *store.Users, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)

	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func(app *application) getUser(ctx context.Context, userId int64) (*store.Users, error) {
	
	if !app.config.redis.enabled{
		return app.store.Users.GetUser(ctx, userId)
	}
	

	user, err := app.cacheStorage.Users.Get(ctx, userId)
	
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		user, err := app.store.Users.GetUser(ctx, userId)
		if err != nil {
			return nil, err
		}
		
		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}


func (app *application) RateLimiterMiddleware (next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.ratelimiter.Enabled{
			if permitted, retryAfter := app.ratelimiter.Permit(r.RemoteAddr); !permitted{
				app.rateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}