package api

import (
	"cat-clerk-api/util"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// StorageRequest -
type StorageRequest struct {
	Title string `json:"title"`
	Owner bool   `json:"owner"`
}

func (api *API) createStorage(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	request := StorageRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	payload, err := api.DB.CreateStorage(username, request.Title, request.Owner)
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(payload, http.StatusOK, w)
}

func (api *API) getStorages(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetStorages(username)
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

func (api *API) getStoragesCount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetStoragesCount(username)
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

func (api *API) updateStorage(w http.ResponseWriter, r *http.Request) {
	storageIDString := mux.Vars(r)["storage_id"]
	storageID, _ := strconv.Atoi(storageIDString)

	request := StorageRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	if err := api.DB.UpdateStorage(request.Title, storageID); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) deleteStorage(w http.ResponseWriter, r *http.Request) {
	storageIDString := mux.Vars(r)["storage_id"]
	storageID, _ := strconv.Atoi(storageIDString)

	if err := api.DB.DeleteStorage(storageID); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) shareStorage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	storageIDString := vars["storage_id"]
	storageID, _ := strconv.Atoi(storageIDString)

	usernameRequest := vars["username_request"]

	if err := api.DB.ShareStorage(usernameRequest, storageID); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) getStorageOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	storageIDString := vars["storage_id"]
	storageID, _ := strconv.Atoi(storageIDString)

	owner := vars["owner"]

	payload, err := api.DB.GetStorageOwner(owner, storageID)
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

func (api *API) removeShareStorageFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	storageIDString := vars["storage_id"]
	storageID, _ := strconv.Atoi(storageIDString)

	usernameRequest := vars["username_request"]

	if err := api.DB.RemoveShareStorage(usernameRequest, storageID); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}
