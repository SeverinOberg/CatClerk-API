package api

import (
	"cat-clerk-api/util"
	"net/http"
)

func (api *API) ping(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(nil, http.StatusNoContent, w)
}

// FoodsResponse structure
type FoodsResponse struct {
	Name string `json:"name"`
}

func (api *API) getFoods(w http.ResponseWriter, r *http.Request) {
	payload, err := api.DB.GetFoods()
	if err != nil {
		util.WriteJSON(util.Error(err.Error()), http.StatusInternalServerError, w)
		return
	}

	names := []FoodsResponse{}

	for _, p := range payload {
		name := FoodsResponse{}
		name.Name = p.Name
		names = append(names, name)
	}

	util.WriteJSON(names, http.StatusOK, w)
}
