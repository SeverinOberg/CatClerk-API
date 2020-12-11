package auth

import (
	"fmt"
	"net/http"
	"strings"

	"cat-clerk-api/util"

	"golang.org/x/time/rate"

	"github.com/dgrijalva/jwt-go"
)

const path = "/api/v1/"

// Auth ...
type Auth struct {
	handler      http.Handler
	hmacSecret   []byte
	tokenLimiter *rate.Limiter
}

// New returns a new Auth object
func New(handler http.Handler, hmacSecret []byte) *Auth {
	return &Auth{
		handler:      handler,
		hmacSecret:   hmacSecret,
		tokenLimiter: rate.NewLimiter(10, 20),
	}
}

func (auth *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		util.AddCORSHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	whiteList := []string{"ping", "sign-up", "login", "forgotten-password", "email-exists"}
	for _, wl := range whiteList {
		if strings.Contains(r.URL.Path, path+wl) {
			auth.handler.ServeHTTP(w, r)
			return
		}
	}

	if auth.tokenLimiter.Allow() == false {
		util.WriteJSON(util.Error("too many requests"), http.StatusTooManyRequests, w)
		return
	}

	token, err := auth.ValidateRequestToken(r, auth.hmacSecret)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnauthorized, w)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		util.WriteJSON(util.Error("unable to assert token claims type"), http.StatusUnauthorized, w)
		return
	}

	pathPrefixes := []string{path + "accounts/"} // Add more in this array if you need to whitelist more paths
	usernamePath, err := getUsernameFromPathPrefixes(r, pathPrefixes)

	// Serve the HTTP request if the username exists and is the same in both the token and request url path.
	usernameClaim := claims["username"].(string)
	if claims["username"] != nil && len(usernameClaim) > 0 && usernamePath == usernameClaim {
		auth.handler.ServeHTTP(w, r)
		return
	}

	util.WriteJSON(util.Error("not authorized"), http.StatusUnauthorized, w)
}

// getUsernameFromPathPrefixes checks if the requested url is a path prefix to look for a username.
func getUsernameFromPathPrefixes(r *http.Request, pathPrefixes []string) (string, error) {
	for _, p := range pathPrefixes {
		if strings.HasPrefix(r.URL.Path, p) {
			trimPath := strings.TrimPrefix(r.URL.Path, p)
			return strings.Trim(strings.Split(trimPath, "/")[0], "/"), nil
		}
	}
	return "", fmt.Errorf("unauthorized URL path")
}
