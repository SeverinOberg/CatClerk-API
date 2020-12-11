package api

import (
	"cat-clerk-api/util"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ShoppingListItemRequest structure
type ShoppingListItemRequest struct {
	Title        string `json:"title"`
	Quantity     int    `json:"quantity"`
	QuantityType string `json:"quantityType"`
}

func (api *API) createShoppingListItem(w http.ResponseWriter, r *http.Request) {
	shoppingListIDString := mux.Vars(r)["shopping_list_id"]
	shoppingListID, _ := strconv.Atoi(shoppingListIDString)

	request := ShoppingListItemRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	if err := api.DB.CreateShoppingListItem(
		shoppingListID,
		request.Title,
		request.Quantity,
		request.QuantityType,
	); err != nil {
		if strings.Contains(err.Error(), "Error 1452") {
			util.WriteJSON(util.Error("must be an already existing shopping list"), http.StatusNotFound, w)
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

func (api *API) getShoppingListItems(w http.ResponseWriter, r *http.Request) {
	shoppingListIDString := mux.Vars(r)["shopping_list_id"]
	shoppingListID, _ := strconv.Atoi(shoppingListIDString)

	payload, err := api.DB.GetShoppingListItems(shoppingListID)
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

func (api *API) getShoppingListItem(w http.ResponseWriter, r *http.Request) {
	itemIDString := mux.Vars(r)["shopping_list_item_id"]
	itemID, _ := strconv.Atoi(itemIDString)

	payload, err := api.DB.GetShoppingListItem(itemID)
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

func (api *API) getShoppingListItemsCount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetShoppingListItemsCount(username)
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

func (api *API) updateShoppingListItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	itemIDstring := vars["shopping_list_item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	request := ShoppingListItemRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	if err := api.DB.UpdateShoppingListItem(
		request.Title,
		request.Quantity,
		request.QuantityType,
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

func (api *API) updateShoppingListItemTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	itemIDstring := vars["shopping_list_item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	if err := api.DB.UpdateShoppingListItemTitle(
		title,
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

func (api *API) decrementShoppingListItemQuantity(w http.ResponseWriter, r *http.Request) {
	itemIDstring := mux.Vars(r)["item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	if err := api.DB.DecrementShoppingListItemQuantity(
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

func (api *API) incrementShoppingListItemQuantity(w http.ResponseWriter, r *http.Request) {
	itemIDstring := mux.Vars(r)["item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	if err := api.DB.IncrementShoppingListItemQuantity(
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

func (api *API) deleteShoppingListItem(w http.ResponseWriter, r *http.Request) {
	itemIDstring := mux.Vars(r)["shopping_list_item_id"]
	itemID, _ := strconv.Atoi(itemIDstring)

	if err := api.DB.DeleteShoppingListItem(itemID); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}
