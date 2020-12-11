package api

import (
	"cat-clerk-api/auth"
	"cat-clerk-api/database"
	"net/http"

	"github.com/gorilla/mux"
)

const path = "/api/v1/"

// API structure
type API struct {
	Router *mux.Router
	DB     *database.Handler
	Auth   *auth.Auth
}

// Init initializes the API package dependencies.
func Init(router *mux.Router, db *database.Handler, auth *auth.Auth) *API {
	return &API{
		Router: router,
		DB:     db,
		Auth:   auth,
	}
}

// Handlers initializes all API handlers.
func (api *API) Handlers() *mux.Router {
	api.Router.Methods(http.MethodGet).
		Path(path + "ping").
		Handler(http.HandlerFunc(api.ping))

	api.Router.Methods(http.MethodPost).
		Path(path + "login").
		Handler(http.HandlerFunc(api.login))

	api.Router.Methods(http.MethodPost).
		Path(path + "sign-up").
		Handler(http.HandlerFunc(api.createAccount))

	api.Router.Methods(http.MethodPost).
		Path(path + "forgotten-password/{email}").
		Handler(http.HandlerFunc(api.sendForgottenPasswordMail))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/foods").
		Handler(http.HandlerFunc(api.getFoods))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}").
		Handler(http.HandlerFunc(api.getAccount))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/email").
		Handler(http.HandlerFunc(api.getAccountEmail))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/email/{email}").
		Handler(http.HandlerFunc(api.updateAccountEmail))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/username/{new_username}").
		Handler(http.HandlerFunc(api.updateAccountUsername))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/password").
		Handler(http.HandlerFunc(api.updateAccountPassword))

	api.Router.Methods(http.MethodPut).
		Path(path + "accounts/{username}").
		Handler(http.HandlerFunc(api.updateAccount))

	api.Router.Methods(http.MethodDelete).
		Path(path + "accounts/{username}").
		Handler(http.HandlerFunc(api.deleteAccount))

	api.Router.Methods(http.MethodPost).
		Path(path + "accounts/{username}/share_requests").
		Handler(http.HandlerFunc(api.createShareRequest))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/share_requests").
		Handler(http.HandlerFunc(api.getShareRequests))

	api.Router.Methods(http.MethodDelete).
		Path(path + "accounts/{username}/share_requests/{share_id}").
		Handler(http.HandlerFunc(api.deleteShareRequest))

	api.Router.Methods(http.MethodPost).
		Path(path + "accounts/{username}/storages").
		Handler(http.HandlerFunc(api.createStorage))

	api.Router.Methods(http.MethodPost).
		Path(path + "accounts/{username}/storages/{storage_id}/share/{username_request}").
		Handler(http.HandlerFunc(api.shareStorage))

	api.Router.Methods(http.MethodDelete).
		Path(path + "accounts/{username}/storages/{storage_id}/share/{username_request}").
		Handler(http.HandlerFunc(api.removeShareStorageFolder))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/storages/{storage_id}/owner/{owner}").
		Handler(http.HandlerFunc(api.getStorageOwner))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/storages").
		Handler(http.HandlerFunc(api.getStorages))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/storages/count").
		Handler(http.HandlerFunc(api.getStoragesCount))

	api.Router.Methods(http.MethodPut).
		Path(path + "accounts/{username}/storages/{storage_id}").
		Handler(http.HandlerFunc(api.updateStorage))

	api.Router.Methods(http.MethodDelete).
		Path(path + "accounts/{username}/storages/{storage_id}").
		Handler(http.HandlerFunc(api.deleteStorage))

	api.Router.Methods(http.MethodPost).
		Path(path + "accounts/{username}/storages/{storage_id}/items").
		Handler(http.HandlerFunc(api.createStorageItem))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/storages/{storage_id}/items").
		Handler(http.HandlerFunc(api.getStorageItems))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/storages/{storage_id}/items/{item_id}").
		Handler(http.HandlerFunc(api.getStorageItem))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/storages/items/count").
		Handler(http.HandlerFunc(api.getStorageItemsCount))

	api.Router.Methods(http.MethodPut).
		Path(path + "accounts/{username}/storages/{storage_id}/items/{item_id}").
		Handler(http.HandlerFunc(api.updateStorageItem))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/storages/{storage_id}/items/{item_id}/quantity/decrement").
		Handler(http.HandlerFunc(api.decrementStorageItemQuantity))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/storages/{storage_id}/items/{item_id}/quantity/increment").
		Handler(http.HandlerFunc(api.incrementStorageItemQuantity))

	api.Router.Methods(http.MethodDelete).
		Path(path + "accounts/{username}/storages/{storage_id}/items/{item_id}").
		Handler(http.HandlerFunc(api.deleteStorageItem))

	api.Router.Methods(http.MethodPost).
		Path(path + "accounts/{username}/shopping-lists").
		Handler(http.HandlerFunc(api.createShoppingList))

	api.Router.Methods(http.MethodPost).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/share/{username_request}").
		Handler(http.HandlerFunc(api.shareShoppingList))

	api.Router.Methods(http.MethodDelete).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/share/{username_request}").
		Handler(http.HandlerFunc(api.removeShareShoppingList))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/owner/{owner}").
		Handler(http.HandlerFunc(api.getShoppingListOwner))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/shopping-lists").
		Handler(http.HandlerFunc(api.getShoppingLists))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/shopping-lists/count").
		Handler(http.HandlerFunc(api.getShoppingListsCount))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/title/{title}").
		Handler(http.HandlerFunc(api.updateShoppingListTitle))

	api.Router.Methods(http.MethodDelete).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}").
		Handler(http.HandlerFunc(api.deleteShoppingList))

	api.Router.Methods(http.MethodPost).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/items").
		Handler(http.HandlerFunc(api.createShoppingListItem))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/items").
		Handler(http.HandlerFunc(api.getShoppingListItems))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/items/{shopping_list_item_id}").
		Handler(http.HandlerFunc(api.getShoppingListItem))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/shopping-lists/items/count").
		Handler(http.HandlerFunc(api.getShoppingListItemsCount))

	api.Router.Methods(http.MethodPut).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/items/{shopping_list_item_id}").
		Handler(http.HandlerFunc(api.updateShoppingListItem))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/items/{shopping_list_item_id}/title/{title}").
		Handler(http.HandlerFunc(api.updateShoppingListItemTitle))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/items/{item_id}/quantity/decrement").
		Handler(http.HandlerFunc(api.decrementShoppingListItemQuantity))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/items/{item_id}/quantity/increment").
		Handler(http.HandlerFunc(api.incrementShoppingListItemQuantity))

	api.Router.Methods(http.MethodDelete).
		Path(path + "accounts/{username}/shopping-lists/{shopping_list_id}/items/{shopping_list_item_id}").
		Handler(http.HandlerFunc(api.deleteShoppingListItem))

	api.Router.Methods(http.MethodGet).
		Path(path + "accounts/{username}/settings/notifications").
		Handler(http.HandlerFunc(api.getNotificatiosSetting))

	api.Router.Methods(http.MethodPatch).
		Path(path + "accounts/{username}/settings/notifications").
		Handler(http.HandlerFunc(api.toggleNotificationSetting))

	return api.Router
}
