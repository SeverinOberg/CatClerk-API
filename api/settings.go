package api

import (
	"cat-clerk-api/util"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (api *API) getNotificatiosSetting(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetNotificationSetting(username)
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

func (api *API) toggleNotificationSetting(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if err := api.DB.ToggleNotificationSetting(
		username,
	); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}
