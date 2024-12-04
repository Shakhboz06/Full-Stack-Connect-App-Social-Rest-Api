package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-project/internal/mailer"
	"go-project/internal/store"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserPostPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type UserToken struct {
	User  *store.Users `json:"user"`
	Token string       `json:"token"`
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

// User Login godoc
//
//	@Summary		Register a user
//	@Description	Register a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UserPostPayload	true	"User Credentials"
//	@Success		201		{object}	UserToken  "User Registered"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/authentication/user [post]
func (app *application) userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload UserPostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &store.Users{
		Username: payload.Username,
		Email:    payload.Email,
		Role: store.Roles{
			Name: "user",
		},
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	token := uuid.New().String()

	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	if err := app.store.Users.CreateAndInvite(ctx, user, hashedToken, app.config.mail.mailExp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequest(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequest(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	usersToken := &UserToken{
		User:  user,
		Token: token,
	}

	isProdEnv := app.config.env == "production"
	activationURl := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, token)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURl,
	}

	status, err := app.mailer.Send(mailer.UserTemp, user.Username, user.Email, vars, !isProdEnv)

	if err != nil {
		app.logger.Errorw("error sending confirmation email", "error", err)

		//rollback in case of email failure
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email has been sent successfully", "status code:", status)

	if err := app.jsonResponse(w, http.StatusCreated, usersToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) userActivationHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	ctx := r.Context()

	err := app.store.Users.Activation(ctx, token)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getUserTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) getUserTokenHandler(w http.ResponseWriter, r *http.Request) {

	userPayload := CreateUserTokenPayload{}

	if err := readJSON(w, r, &userPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(userPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}


	ctx := r.Context()
	user, err := app.store.Users.GetByEmail(ctx, userPayload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.internalServerError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss, 
	}	
	token, err := app.authenticator.GenerateToken(claims)	
	if err != nil{
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}
