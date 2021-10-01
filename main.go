package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/cookbook/repository"
	"github.com/cookbook/service"
	"github.com/gorilla/mux"
)

var serv service.IngredientService
var recipeServ service.RecipeService

func ingredients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ings, err := serv.GetAll()
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(ings)
}

func ingredientByName(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ing, err := serv.Get(id)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(ing)
}

func createIngredient(w http.ResponseWriter, r *http.Request) {
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
	err = serv.Create(i)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func updateIngredient(w http.ResponseWriter, r *http.Request) {
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
	err = serv.Update(i)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func deleteIngredient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	err = serv.Delete(id)
	if err != nil {
		handleError(w, err)
		return
	}
}

func recipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	recipes, err := recipeServ.GetAll()
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(recipes)
}

func recipeById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ing, err := recipeServ.Get(id)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(ing)
}

func createRecipe(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var recipe service.RecipeCreate
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&recipe)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = recipeServ.Create(recipe)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func updateRecipe(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var recipe service.RecipeCreate
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&recipe)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = recipeServ.Update(recipe)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func deleteRecipe(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	err = recipeServ.Delete(id)
	if err != nil {
		handleError(w, err)
		return
	}
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

func main() {
	r := mux.NewRouter()
	repo, _ := repository.NewIngredientRepository()
	recipeRepo, _ := repository.NewRecipeRepository()
	serv = service.NewIngredientService(repo)
	recipeServ = service.NewRecipeService(recipeRepo, serv)
	//r.HandleFunc("/", home)
	ingredientsR := r.PathPrefix("/ingredients").Subrouter()
	ingredientsR.Path("").Methods(http.MethodGet).HandlerFunc(ingredients)
	ingredientsR.Path("").Methods(http.MethodPost).HandlerFunc(createIngredient)
	ingredientsR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(ingredientByName)
	ingredientsR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateIngredient)
	ingredientsR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteIngredient)

	recipeR := r.PathPrefix("/recipes").Subrouter()
	recipeR.Path("").Methods(http.MethodGet).HandlerFunc(recipes)
	recipeR.Path("").Methods(http.MethodPost).HandlerFunc(createRecipe)
	recipeR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(recipeById)
	recipeR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateRecipe)
	recipeR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteRecipe)
	http.ListenAndServe(":3000", r)
}
