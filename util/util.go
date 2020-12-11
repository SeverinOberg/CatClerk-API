package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"unicode"
)

// Utility structure
type Utility struct {
	Salt string
}

var util Utility

// Init initializses the utility package
func Init(salt string) {
	util.Salt = salt
}

// RandomStringGenerator generates a random string with random letters provided (provide as many as possible for more variety)
func RandomStringGenerator() (string, error) {
	if util.Salt == "" {
		return util.Salt, fmt.Errorf("salt can not be empty")
	}

	var saltRunes []rune
	for _, s := range util.Salt {
		saltRunes = append(saltRunes, s)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, rand.Intn(7)+13)
	for i := range b {
		b[i] = saltRunes[rand.Intn(len(saltRunes))]
	}
	return string(b), nil

}

// AddCORSHeaders adds the necessary CORS headers
func AddCORSHeaders(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// WriteJSON writes the JSON output for API calls
func WriteJSON(data interface{}, statusCode int, w http.ResponseWriter) {
	AddCORSHeaders(w)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Error handles error messages
func Error(errorResponse string) interface{} {
	return struct {
		Error string `json:"error"`
	}{
		Error: errorResponse,
	}
}

// PasswordStrengthCheck checks if a password meets the password strength requirements
func PasswordStrengthCheck(password string) (bool, error) {
	var hasMinLen, hasNumber, hasUpper, hasLower, hasSpecial bool

	if len(password) >= 8 {
		hasMinLen = true
	} else {
		return false, errors.New("password must be at least 8 characters")
	}

	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char) || char == ' ':
			hasSpecial = true
		default:
			return false, errors.New("something unexpected happened in the password strength check")
		}
		if hasMinLen && hasNumber && hasLower && hasUpper && hasSpecial {
			return true, nil
		}
	}

	switch {
	case !hasNumber:
		return false, errors.New("password must contain atleast 1 number")
	case !hasLower:
		return false, errors.New("password must contain atleast 1 lower character")
	case !hasUpper:
		return false, errors.New("password must contain atleast 1 capitalized character")
	case !hasSpecial:
		return false, errors.New("password must contain atleast 1 symbol")
	default:
		return false, errors.New("something unexpected happened in the errors password strength check")
	}

}

// ConvertDBTimestamp converts a given time to a database compatible string format
func ConvertToDBTimestamp(t time.Time) string {
	var timestampFormat = "2006-01-02T15:04:05"
	return t.Format(timestampFormat)
}
