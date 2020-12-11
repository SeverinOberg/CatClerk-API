package api

import (
	"cat-clerk-api/util"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (api *API) createShoppingList(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	request := StorageRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusUnprocessableEntity, w)
		return
	}

	if err := api.DB.CreateShoppingList(username, request.Title, request.Owner); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) getShoppingLists(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetShoppingLists(username)
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

func (api *API) getShoppingListsCount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	payload, err := api.DB.GetShoppingListsCount(username)
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

func (api *API) updateShoppingListTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	shoppingListIDString := vars["shopping_list_id"]
	shoppingListID, _ := strconv.Atoi(shoppingListIDString)

	if err := api.DB.UpdateShoppingListTitle(title, shoppingListID); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) deleteShoppingList(w http.ResponseWriter, r *http.Request) {
	shoppingListIDString := mux.Vars(r)["shopping_list_id"]
	shoppingListID, _ := strconv.Atoi(shoppingListIDString)

	if err := api.DB.DeleteShoppingList(shoppingListID); err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			util.WriteJSON(nil, http.StatusNotFound, w)
			return
		}
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) shareShoppingList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shoppingListIDString := vars["shopping_list_id"]
	shoppingListID, _ := strconv.Atoi(shoppingListIDString)

	usernameRequest := vars["username_request"]

	if err := api.DB.ShareShoppingList(usernameRequest, shoppingListID); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}

func (api *API) getShoppingListOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shoppingListIDString := vars["shopping_list_id"]
	shoppingListID, _ := strconv.Atoi(shoppingListIDString)

	owner := vars["owner"]

	payload, err := api.DB.GetShoppingListOwner(owner, shoppingListID)
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

func (api *API) removeShareShoppingList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shoppingListIDString := vars["shopping_list_id"]
	shoppingListID, _ := strconv.Atoi(shoppingListIDString)

	usernameRequest := vars["username_request"]

	if err := api.DB.RemoveShareShoppingList(usernameRequest, shoppingListID); err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	util.WriteJSON(nil, http.StatusNoContent, w)
}
