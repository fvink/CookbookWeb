package repository

import (
	"database/sql"
	"log"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository() (*RecipeRepository, error) {
	r := new(RecipeRepository)
	var err error
	r.db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/cookbook")
	if err != nil {
		log.Println(err.Error())
	}
	return r, err
}

func (r RecipeRepository) Get(id int64) (recipe Recipe, e error) {
	err := r.db.QueryRow("SELECT * FROM recipes WHERE id = ?", id).Scan(&recipe.Id, &recipe.Name, &recipe.Steps)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case sql.ErrNoRows:
			e = &NotFound{"recipes", id}
		default:
			e = &InternalError{err.Error()}
		}
	}
	recipe.Ingredients, e = r.getRecipeIngredients(id)
	return
}

func (r RecipeRepository) getRecipeIngredients(recipeId int64) (ingredients []IngredientShort, e error) {
	results, err := r.db.Query("SELECT ingredient_id, amount, unit FROM recipe_ingredients WHERE recipe_id = ? ORDER BY recipe_ingredients.index", recipeId)
	if err != nil {
		return []IngredientShort{}, &InternalError{err.Error()}
	}
	for results.Next() {
		var i IngredientShort
		err = results.Scan(&i.Id, &i.Amount, &i.Unit)
		if err != nil {
			log.Println(err.Error())
		}
		ingredients = append(ingredients, i)
	}
	return ingredients, nil
}

func (r RecipeRepository) GetAll() (recipes []Recipe, err error) {
	results, err := r.db.Query("SELECT * FROM recipes")
	if err != nil {
		return []Recipe{}, &InternalError{err.Error()}
	}
	for results.Next() {
		var recipe Recipe
		err = results.Scan(&recipe.Id, &recipe.Name, &recipe.Steps)
		if err != nil {
			log.Println(err.Error())
		}
		recipe.Ingredients, err = r.getRecipeIngredients(recipe.Id)
		if err != nil {
			log.Println(err.Error())
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r RecipeRepository) Create(recipe Recipe) error {
	result, err := r.db.Exec("INSERT INTO recipes (name, steps) VALUES (?, ?)", recipe.Name, recipe.Steps)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return &InternalError{"recipe not created"}
	}
	return r.createRecipeIngredients(recipe)
}

func (r RecipeRepository) Update(recipe Recipe) error {
	err := r.deleteRecipeIngredients(recipe.Id)
	if err != nil {
		return err
	}
	result, err := r.db.Exec("UPDATE recipes SET name = ?, steps = ? WHERE id = ?", recipe.Name, recipe.Steps, recipe.Id)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case sql.ErrNoRows:
			return &NotFound{"recipes", recipe.Id}
		default:
			return &InternalError{err.Error()}
		}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return &InternalError{"recipe not updated"}
	}
	return r.createRecipeIngredients(recipe)
}

func (r RecipeRepository) createRecipeIngredients(recipe Recipe) error {
	for index, ing := range recipe.Ingredients {
		result, err := r.db.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount, unit, index) VALUES (?, ?, ?, ?, ?)", recipe.Id, ing.Id, ing.Amount, ing.Unit, index)
		if err != nil {
			return &InternalError{err.Error()}
		}
		rowCnt, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowCnt != 1 {
			return &InternalError{"recipe ingredient not created"}
		}
	}
	return nil
}

func (r RecipeRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM recipes WHERE id = ?", id)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return &InternalError{err.Error()}
	}
	if rowCnt != 1 {
		return &NotFound{"recipes", id}
	}
	return nil
}

func (r RecipeRepository) deleteRecipeIngredients(recipeId int64) error {
	_, err := r.db.Exec("DELETE FROM recipe_ingredients WHERE recipe_id = ?", recipeId)
	if err != nil {
		return &InternalError{err.Error()}
	}
	return nil
}
