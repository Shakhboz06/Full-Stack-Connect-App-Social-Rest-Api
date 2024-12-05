package auth

import (
	// "time"

	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct {
	// GenerateToken(jwt.Claims) (string, error)
	// ValidateToken(string)(*jwt.Token, error)
}

const secret = "123"

// var testClaims = jwt.MapClaims{
// 	"sub": int64(951),
// 	"exp": time.Now().Add(time.Hour * 3).Unix(),
// 	"iat": time.Now().Unix(),
// 	"nbf": time.Now().Unix(),
// 	"iss": "test-aus",
// 	"aud": "test-aud",
// }

func (a *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))


	return tokenString, nil
}

func (a *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
