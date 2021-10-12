package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cookbook/config"
	"github.com/cookbook/repository"
	"github.com/cookbook/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

var serv service.IngredientService
var recipeServ service.RecipeService
var mealServ service.MealService
var mealPlanServ service.MealPlanService

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

func meals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	meals, err := mealServ.GetAll()
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(meals)
}

func mealById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	meal, err := mealServ.Get(id)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(meal)
}

func createMeal(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var meal service.MealCreate
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&meal)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = mealServ.Create(meal)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func updateMeal(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var meal service.MealCreate
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&meal)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = mealServ.Update(meal)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func deleteMeal(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	err = mealServ.Delete(id)
	if err != nil {
		handleError(w, err)
		return
	}
}

func mealPlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mealPlans, err := mealPlanServ.GetAll()
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(mealPlans)
}

func mealPlanById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	mealPlan, err := mealPlanServ.Get(id)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(mealPlan)
}

func createMealPlan(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var mealPlan service.MealPlanCreate
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&mealPlan)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = mealPlanServ.Create(mealPlan)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func updateMealPlan(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var mealPlan service.MealPlanCreate
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&mealPlan)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = mealPlanServ.Update(mealPlan)
	if err != nil {
		handleError(w, err)
		return
	}
	errorResponse(w, "Success", http.StatusOK)
}

func deleteMealPlan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}
	err = mealPlanServ.Delete(id)
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
	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	router := mux.NewRouter()
	databaseUrl := fmt.Sprintf("%s://%s:%s@%s:%s/%s", conf.Database.Protocol,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Server,
		conf.Database.Port,
		conf.Database.Name)
	dbConn, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbConn.Close()
	repo := repository.NewIngredientRepository(dbConn)
	recipeRepo := repository.NewRecipeRepository(dbConn)
	mealRepo := repository.NewMealRepository(dbConn)
	mealPlanRepo := repository.NewMealPlanRepository(dbConn)
	serv = service.NewIngredientService(repo)
	recipeServ = service.NewRecipeService(recipeRepo, serv)
	mealServ = service.NewMealService(mealRepo, recipeServ)
	mealPlanServ = service.NewMealPlanService(mealPlanRepo, mealServ)
	//r.HandleFunc("/", home)
	ingredientsR := router.PathPrefix("/ingredients").Subrouter()
	ingredientsR.Path("").Methods(http.MethodGet).HandlerFunc(ingredients)
	ingredientsR.Path("").Methods(http.MethodPost).HandlerFunc(createIngredient)
	ingredientsR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(ingredientByName)
	ingredientsR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateIngredient)
	ingredientsR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteIngredient)

	recipeR := router.PathPrefix("/recipes").Subrouter()
	recipeR.Path("").Methods(http.MethodGet).HandlerFunc(recipes)
	recipeR.Path("").Methods(http.MethodPost).HandlerFunc(createRecipe)
	recipeR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(recipeById)
	recipeR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateRecipe)
	recipeR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteRecipe)

	mealR := router.PathPrefix("/meals").Subrouter()
	mealR.Path("").Methods(http.MethodGet).HandlerFunc(meals)
	mealR.Path("").Methods(http.MethodPost).HandlerFunc(createMeal)
	mealR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(mealById)
	mealR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateMeal)
	mealR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteMeal)

	mealPlanR := router.PathPrefix("/meal-plans").Subrouter()
	mealPlanR.Path("").Methods(http.MethodGet).HandlerFunc(mealPlans)
	mealPlanR.Path("").Methods(http.MethodPost).HandlerFunc(createMealPlan)
	mealPlanR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(mealPlanById)
	mealPlanR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateMealPlan)
	mealPlanR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteMealPlan)
	http.ListenAndServe(conf.Server.Host+":"+conf.Server.Port, router)
}
