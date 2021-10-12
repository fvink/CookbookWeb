package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cookbook/service"
	"github.com/gorilla/mux"
)

type IngredientHandler struct {
	Service service.IngredientService
}

func (handler IngredientHandler) Register(router *mux.Router) {
	subrouter := router.PathPrefix("/ingredients").Subrouter()
	subrouter.Path("").Methods(http.MethodGet).HandlerFunc(handler.ingredients)
	subrouter.Path("").Methods(http.MethodPost).HandlerFunc(handler.createIngredient)
	subrouter.Path("/{id}").Methods(http.MethodGet).HandlerFunc(handler.ingredientByName)
	subrouter.Path("/{id}").Methods(http.MethodPut).HandlerFunc(handler.updateIngredient)
	subrouter.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(handler.deleteIngredient)
}

func (handler IngredientHandler) ingredients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ings, err := handler.Service.GetAll()
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(ings)
}

func (handler IngredientHandler) ingredientByName(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ing, err := handler.Service.Get(id)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(ing)
}

func (handler IngredientHandler) createIngredient(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var i service.Ingredient
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&i)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = handler.Service.Create(i)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func (handler IngredientHandler) updateIngredient(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var i service.Ingredient
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&i)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = handler.Service.Update(i)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func (handler IngredientHandler) deleteIngredient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	err = handler.Service.Delete(id)
	if err != nil {
		handleError(w, err)
		return
	}
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func handleError(w http.ResponseWriter, err error) {
	switch x := err.(type) {
	case *service.NotFound:
		errorResponse(w, x.Error(), http.StatusNotFound)
	case *service.ValidationError:
		errorResponse(w, x.Error(), http.StatusBadRequest)
	default:
		errorResponse(w, "Internal server error, if the error persists contact server admin", http.StatusInternalServerError)
	}
}
