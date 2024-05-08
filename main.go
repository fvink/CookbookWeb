package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cookbook/handler"
	"github.com/cookbook/repository"
	"github.com/cookbook/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("PORT")
	
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

	serv := service.NewIngredientService(repo)
	recipeServ := service.NewRecipeService(recipeRepo, serv)
	mealServ := service.NewMealService(mealRepo, recipeServ)
	mealPlanServ := service.NewMealPlanService(mealPlanRepo, mealServ)

	router := handler.NewRestRouter()
	ingredientHandler := handler.IngredientHandler{Service: serv}
	recipeHandler := handler.RecipeHandler{Service: recipeServ}
	mealHandler := handler.MealHandler{Service: mealServ}
	mealPlanHandler := handler.MealPlanHandler{Service: mealPlanServ}

	router.Register("ingredients", ingredientHandler)
	router.Register("recipes", recipeHandler)
	router.Register("meals", mealHandler)
	router.Register("meal-plans", mealPlanHandler)

	http.ListenAndServe(serverHost + ":" + serverPort, router)
}
