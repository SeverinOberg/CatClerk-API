package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"cat-clerk-api/mail"
	"cat-clerk-api/util"

	"github.com/gorilla/mux"
)

func (api *API) sendForgottenPasswordMail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]

	acc, err := api.DB.EmailExists(email)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	claims := map[string]interface{}{
		"username": acc.Username,
	}

	accessExp := time.Now().Unix() + int64(5*time.Minute.Seconds()) // Now + 5 minutes in milliseconds.

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

	data := struct {
		URL string
	}{
		URL: fmt.Sprintf("http://localhost:8080/tabs/%s/reset-password/%s/%d", acc.Username, signedAccessToken, accessExp),
	}

	if err := mail.SendEmailOAUTH2(
		email,
		"Forgotten Password | Cat Clerk",
		data,
		"forgotten-password.gohtml",
	); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}
