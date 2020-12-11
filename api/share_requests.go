package api

import (
	"cat-clerk-api/util"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ShareRequest structure
type ShareRequest struct {
	ToUsername string `json:"toUsername"`
	ShareType  string `json:"shareType"`
	Title      string `json:"title"`
	IDRequest  int    `json:"idRequest"`
}

func (api *API) createShareRequest(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	request := ShareRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	accounts, err := api.DB.GetAccounts()
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusNotFound, w)
		return
	}

	accountExists := false
	for _, a := range accounts {
		if a.Username == request.ToUsername {
			accountExists = true
			break
		}
	}

	if !accountExists {
		util.WriteJSON(util.Error("no such account exists"), http.StatusNotFound, w)
		return
	}

	currentShareRequests, err := api.DB.GetShareRequests(request.ToUsername)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	for _, csr := range currentShareRequests {
		if csr.ToUsername == request.ToUsername && csr.IDRequest == request.IDRequest {
			util.WriteJSON(util.Error("this request already exists"), http.StatusUnprocessableEntity, w)
			return
		}
	}

	if request.ShareType != "storage" && request.ShareType != "shopping_list" {
		util.WriteJSON(util.Error("type must be 'storage' or 'shopping_list'"), http.StatusUnprocessableEntity, w)
		return
	}

	if err := api.DB.CreateShareRequest(username, request.ToUsername, request.ShareType, request.Title, request.IDRequest); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) getShareRequests(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetShareRequests(username)
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

func (api *API) deleteShareRequest(w http.ResponseWriter, r *http.Request) {
	shareIDstring := mux.Vars(r)["share_id"]

	shareID, _ := strconv.Atoi(shareIDstring)

	if err := api.DB.DeleteShareRequest(shareID); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}
