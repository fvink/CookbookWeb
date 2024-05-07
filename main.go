package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cookbook/config"
	"github.com/cookbook/handler"
	"github.com/cookbook/repository"
	"github.com/cookbook/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// conf, err := config.LoadConfig("config.yaml")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }
	// databaseUrl := fmt.Sprintf("%s://%s:%s@%s:%s/%s", conf.Database.Protocol,
	// 	conf.Database.Username,
	// 	conf.Database.Password,
	// 	conf.Database.Server,
	// 	conf.Database.Port,
	// 	conf.Database.Name)

	dbConn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
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

	http.ListenAndServe(conf.Server.Host+":"+conf.Server.Port, router)
}
