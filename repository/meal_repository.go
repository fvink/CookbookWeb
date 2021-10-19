package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MealRepository struct {
	db *pgxpool.Pool
}

func NewMealRepository(dbConn *pgxpool.Pool) *MealRepository {
	r := new(MealRepository)
	r.db = dbConn
	return r
}

func (r MealRepository) Get(id int64) (meal Meal, e error) {
	err := r.db.QueryRow(context.Background(), "SELECT * FROM meals WHERE id = $1", id).Scan(&meal.Id, &meal.Name)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case pgx.ErrNoRows:
			return Meal{}, &NotFound{"meals", id}
		default:
			return Meal{}, &InternalError{err.Error()}
		}
	}
	recipes, err := r.getMealRecipes([]int64{id})
	if err != nil {
		return Meal{}, err
	}
	meal.Recipes = recipes[id]
	return meal, nil
}

func (r MealRepository) getMealRecipes(ids []int64) (recipes map[int64][]int64, err error) {
	recipes = make(map[int64][]int64)
	results, err := r.db.Query(context.Background(), "SELECT meal_id, recipe_id FROM meal_recipes WHERE meal_id IN ("+JoinIds(ids)+") ORDER BY recipe_id, meal_recipes.index")
	if err != nil {
		return nil, &InternalError{err.Error()}
	}
	for results.Next() {
		var mealId, recipeId int64
		err = results.Scan(&mealId, &recipeId)
		if err != nil {
			log.Println(err.Error())
		}
		if _, ok := recipes[mealId]; !ok {
			recipes[mealId] = make([]int64, 0)
		}
		recipes[mealId] = append(recipes[mealId], recipeId)
	}
	return recipes, nil
}

func (r MealRepository) getAllMealRecipes() (recipes map[int64][]int64, err error) {
	recipes = make(map[int64][]int64)
	results, err := r.db.Query(context.Background(), "SELECT meal_id, recipe_id FROM meal_recipes ORDER BY recipe_id, meal_recipes.index")
	if err != nil {
		return nil, &InternalError{err.Error()}
	}
	for results.Next() {
		var mealId, recipeId int64
		err = results.Scan(&mealId, &recipeId)
		if err != nil {
			log.Println(err.Error())
		}
		if _, ok := recipes[mealId]; !ok {
			recipes[mealId] = make([]int64, 0)
		}
		recipes[mealId] = append(recipes[mealId], recipeId)
	}
	return recipes, nil
}

func (r MealRepository) GetAll() (meals []Meal, e error) {
	results, err := r.db.Query(context.Background(), "SELECT * FROM meals")
	if err != nil {
		return []Meal{}, &InternalError{err.Error()}
	}
	mealRecipes, err := r.getAllMealRecipes()
	if err != nil {
		return []Meal{}, err
	}
	for results.Next() {
		var meal Meal
		err = results.Scan(&meal.Id, &meal.Name)
		if err != nil {
			log.Println(err.Error())
		}
		meal.Recipes = mealRecipes[meal.Id]
		if err != nil {
			log.Println(err.Error())
		}
		meals = append(meals, meal)
	}
	return meals, nil
}

func (r MealRepository) GetList(ids []int64) (meals []Meal, e error) {
	results, err := r.db.Query(context.Background(), "SELECT * FROM meals WHERE id IN ("+JoinIds(ids)+")")
	if err != nil {
		return []Meal{}, &InternalError{err.Error()}
	}
	mealRecipes, err := r.getMealRecipes(ids)
	if err != nil {
		return []Meal{}, err
	}
	for results.Next() {
		var meal Meal
		err = results.Scan(&meal.Id, &meal.Name)
		if err != nil {
			log.Println(err.Error())
		}
		meal.Recipes = mealRecipes[meal.Id]
		if err != nil {
			log.Println(err.Error())
		}
		meals = append(meals, meal)
	}
	return meals, nil
}

func (r MealRepository) Create(meal Meal) error {
	err := r.db.QueryRow(context.Background(), "INSERT INTO meals (name) VALUES ($1) RETURNING id", meal.Name).Scan(&meal.Id)
	if err != nil {
		log.Println(err.Error())
		return &InternalError{err.Error()}
	}
	return r.createMealRecipes(meal)
}

func (r MealRepository) Update(meal Meal) error {
	err := r.deleteMealRecipes(meal.Id)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(context.Background(), "UPDATE meals SET name = $1 WHERE id = $2", meal.Name, meal.Id)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case pgx.ErrNoRows:
			return &NotFound{"meals", meal.Id}
		default:
			return &InternalError{err.Error()}
		}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &InternalError{"meal not updated"}
	}
	return r.createMealRecipes(meal)
}

func (r MealRepository) createMealRecipes(meal Meal) error {
	for index, recipeId := range meal.Recipes {
		result, err := r.db.Exec(context.Background(), "INSERT INTO meal_recipes (meal_id, recipe_id, index) VALUES ($1, $2, $3)", meal.Id, recipeId, index)
		if err != nil {
			log.Println(err.Error())
			return &InternalError{err.Error()}
		}
		rowCnt := result.RowsAffected()
		if rowCnt != 1 {
			log.Println(err.Error())
			return &InternalError{"meal recipe not created"}
		}
	}
	return nil
}

func (r MealRepository) Delete(id int64) error {
	result, err := r.db.Exec(context.Background(), "DELETE FROM meals WHERE id = $1", id)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &NotFound{"meals", id}
	}
	return nil
}

func (r MealRepository) deleteMealRecipes(mealId int64) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM meal_recipes WHERE meal_id = $1", mealId)
	if err != nil {
		return &InternalError{err.Error()}
	}
	return nil
}
