package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Token ...
type Token struct {
	AccessToken      string `json:"accessToken"`
	AccessExpiresAt  int64  `json:"accessExpiresAt"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresAt int64  `json:"refreshExpiresAt"`
	Username         string `json:"username"`
}

// CreateJWTToken creates a JWTToken and returns an access and refresh token signature.
func (auth *Auth) CreateJWTToken(customClaims map[string]interface{}, exp int64) (*jwt.Token, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Unix() + exp,
	}

	for key, value := range customClaims {
		claims[key] = value
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims), nil
}

// SignToken signs and returns the signed token string.
func (auth *Auth) SignToken(jwtToken *jwt.Token) (string, error) {
	signedString, err := jwtToken.SignedString(auth.hmacSecret)
	if err != nil {
		return signedString, err
	}
	return signedString, nil
}

// ValidateRequestToken validates a JWTToken.
func (auth *Auth) ValidateRequestToken(r *http.Request, hmacSecret []byte) (*jwt.Token, error) {
	if r.Header["Authorization"] == nil {
		return nil, fmt.Errorf("Authorization header is empty")
	}

	if !strings.Contains(r.Header["Authorization"][0], "Bearer ") {
		return nil, fmt.Errorf("missing Bearer in token")
	}

	tokenString := strings.Split(r.Header["Authorization"][0], "Bearer ")[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token not valid")
	}

	return token, nil
}
