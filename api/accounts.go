package api

import (
	"cat-clerk-api/auth"
	"cat-clerk-api/util"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// CreateAccount add a new account.
func (api *API) createAccount(w http.ResponseWriter, r *http.Request) {
	request := Login{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	if request.Username == "" || request.Email == "" || request.Password == "" {
		util.WriteJSON(util.Error("input fields can't be empty"), http.StatusNotAcceptable, w)
		return
	}

	if !strings.Contains(request.Email, "@") || !strings.Contains(request.Email, ".") {
		util.WriteJSON(util.Error(("not a valid email")), http.StatusNotAcceptable, w)
		return
	}

	pswCheck, err := util.PasswordStrengthCheck(request.Password)
	if pswCheck == false {
		util.WriteJSON(util.Error(err.Error()), http.StatusNotAcceptable, w)
		return
	}

	salt, err := util.RandomStringGenerator()
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	result, err := api.DB.CreateAccount(request.Username, request.Email, string(hashedPassword), salt)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	if rowsAff == 0 {
		util.WriteJSON(nil, http.StatusNotFound, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

// Login -
type Login struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate account credentials and retrieve an access token upon successful sign-in
func (api *API) login(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		util.WriteJSON(util.Error("request body is empty"), http.StatusUnauthorized, w)
		return
	}

	request := Login{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	account, err := api.DB.CheckAccountCredentials(request.Username, request.Email)
	if err != nil {
		util.WriteJSON(util.Error("wrong login or password"), http.StatusUnauthorized, w)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password+account.Salt)); err != nil {
		util.WriteJSON(util.Error("wrong login or password"), http.StatusUnauthorized, w)
		return
	}

	claims := map[string]interface{}{
		"username": account.Username,
	}

	// Access Token
	accessExp := time.Now().Unix() + int64(60*time.Minute.Seconds()) // Now + 60 minutes in milliseconds.

	jwtAccessToken, err := api.Auth.CreateJWTToken(claims, accessExp)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusBadRequest, w)
		return
	}

	signedAccessToken, err := api.Auth.SignToken(jwtAccessToken)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusBadRequest, w)
		return
	}

	// Refresh Token
	refreshExp := time.Now().Unix() + int64(24*time.Hour.Seconds()) // Now + 24 hours in milliseconds.

	jwtRefreshToken, err := api.Auth.CreateJWTToken(claims, refreshExp)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusBadRequest, w)
		return
	}

	signedRefreshToken, err := api.Auth.SignToken(jwtRefreshToken)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusBadRequest, w)
		return
	}

	util.WriteJSON(auth.Token{
		AccessToken:      signedAccessToken,
		AccessExpiresAt:  accessExp,
		RefreshToken:     signedRefreshToken,
		RefreshExpiresAt: refreshExp,
		Username:         account.Username,
	}, http.StatusOK, w)
}

func (api *API) getAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetAccount(username)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(payload, http.StatusOK, w)
}

func (api *API) getAccountEmail(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetAccountEmail(username)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(payload, http.StatusOK, w)
}

func (api *API) emailExists(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		util.WriteJSON(util.Error("request body is empty"), http.StatusUnauthorized, w)
		return
	}

	request := struct {
		Email string `json:"email"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	_, err := api.DB.EmailExists(request.Email)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(true, http.StatusOK, w)
}

// AccountRequest structure
type AccountRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	DarkTheme     bool   `json:"darkTheme"`
	Notifications bool   `json:"notifications"`
}

func (api *API) updateAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	request := AccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	if request.Username == "" || request.Password == "" || request.Email == "" {
		util.WriteJSON(util.Error("input fields can't be empty"), http.StatusNotAcceptable, w)
		return
	}

	if !strings.Contains(request.Email, "@") || !strings.Contains(request.Email, ".") {
		util.WriteJSON(util.Error(("not a valid email")), http.StatusNotAcceptable, w)
		return
	}

	if pswCheck, err := util.PasswordStrengthCheck(request.Password); pswCheck == false {
		util.WriteJSON(util.Error(err.Error()), http.StatusNotAcceptable, w)
		return
	}

	accounts, err := api.DB.GetAccounts()
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	for _, acc := range accounts {
		if acc.Username == request.Username {
			util.WriteJSON(util.Error("that username is already taken"), http.StatusNotAcceptable, w)
			return
		}
		if acc.Email == request.Email {
			util.WriteJSON(util.Error("that email is already taken"), http.StatusNotAcceptable, w)
			return
		}
	}
	if err := api.DB.UpdateAccount(
		username,
		request.Username,
		request.Password,
		request.Email,
		request.DarkTheme,
		request.Notifications,
	); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) updateAccountEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	email := vars["email"]

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		util.WriteJSON(util.Error(("not a valid email")), http.StatusNotAcceptable, w)
		return
	}

	accounts, err := api.DB.GetAccounts()
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	for _, acc := range accounts {
		if acc.Email == email {
			util.WriteJSON(util.Error("that email is already taken"), http.StatusNotAcceptable, w)
			return
		}
	}

	if err := api.DB.UpdateAccountEmail(
		username,
		email,
	); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) updateAccountUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	username := vars["username"]
	newUsername := vars["new_username"]

	accounts, err := api.DB.GetAccounts()
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	for _, acc := range accounts {
		if acc.Username == newUsername {
			util.WriteJSON(util.Error("that email is already taken"), http.StatusNotAcceptable, w)
			return
		}
	}

	if err := api.DB.UpdateAccountUsername(
		username,
		newUsername,
	); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) updateAccountPassword(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	request := struct {
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	pswCheck, err := util.PasswordStrengthCheck(request.Password)
	if pswCheck == false {
		util.WriteJSON(util.Error(err.Error()), http.StatusNotAcceptable, w)
		return
	}

	salt, err := util.RandomStringGenerator()
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	if err := api.DB.UpdateAccountPassword(
		username,
		string(hashedPassword),
		salt,
	); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) deleteAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if err := api.DB.DeleteAccount(username); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}
