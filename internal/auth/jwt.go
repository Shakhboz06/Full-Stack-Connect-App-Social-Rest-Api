package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	var JWTToken JWTAuthenticator
	JWTToken.aud = aud
	JWTToken.secret = secret
	JWTToken.iss = iss

	return &JWTToken
}

func (tk *JWTAuthenticator) GenerateToken(claims jwt.Claims)(string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tk.secret))

	if err != nil{
		return "", err
	}

	return tokenString, nil

}


func(tk *JWTAuthenticator)ValidateToken(token string)(*jwt.Token, error){
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}


		return []byte(tk.secret),nil
	},
	jwt.WithExpirationRequired(),
	jwt.WithAudience(tk.aud),
	jwt.WithIssuer(tk.iss),
	jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}