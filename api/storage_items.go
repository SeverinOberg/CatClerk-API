package api

import (
	"cat-clerk-api/util"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ItemRequest structure
type ItemRequest struct {
	Title               string `json:"title"`
	Image               string `json:"image"`
	Quantity            int    `json:"quantity"`
	QuantityType        string `json:"quantityType"`
	QuantityThreshold   int    `json:"quantityThreshold"`
	ExpirationThreshold int    `json:"expirationThreshold"`
	ExpirationDate      string `json:"expirationDate"`
}

func (api *API) createStorageItem(w http.ResponseWriter, r *http.Request) {
	storageIDString := mux.Vars(r)["storage_id"]
	storageID, _ := strconv.Atoi(storageIDString)

	request := ItemRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	if err := api.DB.CreateStorageItem(
		storageID,
		request.Title,
		request.Quantity,
		request.QuantityType,
		request.QuantityThreshold,
		request.ExpirationThreshold,
		request.ExpirationDate,
	); err != nil {
		if strings.Contains(err.Error(), "Error 1452") {
			util.WriteJSON(util.Error("must be an already existing item"), http.StatusNotFound, w)
			return
		}
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) getStorageItems(w http.ResponseWriter, r *http.Request) {
	storageIDString := mux.Vars(r)["storage_id"]
	storageID, _ := strconv.Atoi(storageIDString)

	payload, err := api.DB.GetStorageItems(storageID)
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

func (api *API) getStorageItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	storageIDString := vars["storage_id"]
	storageID, _ := strconv.Atoi(storageIDString)
	itemIDstring := vars["item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	payload, err := api.DB.GetStorageItem(storageID, itemID)
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

func (api *API) getStorageItemsCount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetStorageItemsCount(username)
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

func (api *API) updateStorageItem(w http.ResponseWriter, r *http.Request) {
	itemIDstring := mux.Vars(r)["item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	request := ItemRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	if err := api.DB.UpdateStorageItem(
		request.Title,
		request.Image,
		request.Quantity,
		request.QuantityType,
		request.QuantityThreshold,
		request.ExpirationThreshold,
		request.ExpirationDate,
		itemID,
	); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) decrementStorageItemQuantity(w http.ResponseWriter, r *http.Request) {
	itemIDstring := mux.Vars(r)["item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	if err := api.DB.DecrementStorageItemQuantity(
		itemID,
	); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) incrementStorageItemQuantity(w http.ResponseWriter, r *http.Request) {
	itemIDstring := mux.Vars(r)["item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	if err := api.DB.IncrementStorageItemQuantity(
		itemID,
	); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) deleteStorageItem(w http.ResponseWriter, r *http.Request) {
	itemIDstring := mux.Vars(r)["item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	if err := api.DB.DeleteStorageItem(itemID); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}
