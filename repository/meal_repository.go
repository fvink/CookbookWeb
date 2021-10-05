package repository

import (
	"database/sql"
	"log"
)

type MealRepository struct {
	db *sql.DB
}

func NewMealRepository() (*MealRepository, error) {
	r := new(MealRepository)
	var err error
	r.db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/cookbook")
	if err != nil {
		log.Println(err.Error())
	}
	return r, err
}

func (r MealRepository) Get(id int64) (meal Meal, e error) {
	err := r.db.QueryRow("SELECT * FROM meals WHERE id = ?", id).Scan(&meal.Id, &meal.Name)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case sql.ErrNoRows:
			e = &NotFound{"meals", id}
		default:
			e = &InternalError{err.Error()}
		}
	}
	meal.Recipes, e = r.getMealRecipes(id)
	return
}

func (r MealRepository) getMealRecipes(id int64) (recipes []int64, err error) {
	results, err := r.db.Query("SELECT recipe_id FROM meal_recipes WHERE meal_id = ? ORDER BY meal_recipes.index", id)
	if err != nil {
		return []int64{}, &InternalError{err.Error()}
	}
	for results.Next() {
		var recipeId int64
		err = results.Scan(&recipeId)
		if err != nil {
			log.Println(err.Error())
		}
		recipes = append(recipes, recipeId)
	}
	return recipes, nil
}

func (r MealRepository) GetAll() (meals []Meal, e error) {
	results, err := r.db.Query("SELECT * FROM meals")
	if err != nil {
		return []Meal{}, &InternalError{err.Error()}
	}
	for results.Next() {
		var meal Meal
		err = results.Scan(&meal.Id, &meal.Name)
		if err != nil {
			log.Println(err.Error())
		}
		meal.Recipes, err = r.getMealRecipes(meal.Id)
		if err != nil {
			log.Println(err.Error())
		}
		meals = append(meals, meal)
	}
	return meals, nil
}

func (r MealRepository) Create(meal Meal) error {
	result, err := r.db.Exec("INSERT INTO meals (name) VALUES (?)", meal.Name)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return &InternalError{"meal not created"}
	}
	return r.createMealRecipes(meal)
}

func (r MealRepository) Update(meal Meal) error {
	err := r.deleteMealRecipes(meal.Id)
	if err != nil {
		return err
	}
	result, err := r.db.Exec("UPDATE meals SET name = ? WHERE id = ?", meal.Name, meal.Id)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case sql.ErrNoRows:
			return &NotFound{"meals", meal.Id}
		default:
			return &InternalError{err.Error()}
		}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return &InternalError{"meal not updated"}
	}
	return r.createMealRecipes(meal)
}

func (r MealRepository) createMealRecipes(meal Meal) error {
	for index, recipeId := range meal.Recipes {
		result, err := r.db.Exec("INSERT INTO meal_recipes (meal_id, recipe_id, index) VALUES (?, ?, ?)", meal.Id, recipeId, index)
		if err != nil {
			return &InternalError{err.Error()}
		}
		rowCnt, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowCnt != 1 {
			return &InternalError{"meal recipe not created"}
		}
	}
	return nil
}

func (r MealRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM meals WHERE id = ?", id)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return &InternalError{err.Error()}
	}
	if rowCnt != 1 {
		return &NotFound{"meals", id}
	}
	return nil
}

func (r MealRepository) deleteMealRecipes(mealId int64) error {
	_, err := r.db.Exec("DELETE FROM meal_recipes WHERE meal_id = ?", mealId)
	if err != nil {
		return &InternalError{err.Error()}
	}
	return nil
}
