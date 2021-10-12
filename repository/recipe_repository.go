package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RecipeRepository struct {
	db *pgxpool.Pool
}

func NewRecipeRepository(dbConn *pgxpool.Pool) *RecipeRepository {
	r := new(RecipeRepository)
	r.db = dbConn
	return r
}

func (r RecipeRepository) Get(id int64) (recipe Recipe, e error) {
	err := r.db.QueryRow(context.Background(), "SELECT * FROM recipes WHERE id = $1", id).Scan(&recipe.Id, &recipe.Name, &recipe.Steps)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case pgx.ErrNoRows:
			e = &NotFound{"recipes", id}
		default:
			e = &InternalError{err.Error()}
		}
	}
	recipe.Ingredients, e = r.getRecipeIngredients(id)
	return
}

func (r RecipeRepository) getRecipeIngredients(recipeId int64) (ingredients []IngredientShort, e error) {

	results, err := r.db.Query(context.Background(), "SELECT ingredient_id, amount, unit FROM recipe_ingredients WHERE recipe_id = $1 ORDER BY recipe_ingredients.index", recipeId)
	if err != nil {
		log.Println(err.Error())
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
	results, err := r.db.Query(context.Background(), "SELECT * FROM recipes")
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

func (r RecipeRepository) GetList(ids []int64) (recipes []Recipe, err error) {
	results, err := r.db.Query(context.Background(), "SELECT * FROM recipes WHERE id IN ("+JoinIds(ids)+")")
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
	err := r.db.QueryRow(context.Background(), "INSERT INTO recipes (name, steps) VALUES ($1, $2) RETURNING id", recipe.Name, recipe.Steps).Scan(&recipe.Id)
	if err != nil {
		log.Println(err.Error())
		return &InternalError{err.Error()}
	}
	return r.createRecipeIngredients(recipe)
}

func (r RecipeRepository) Update(recipe Recipe) error {
	err := r.deleteRecipeIngredients(recipe.Id)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(context.Background(), "UPDATE recipes SET name = $1, steps = $2 WHERE id = $3", recipe.Name, recipe.Steps, recipe.Id)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case pgx.ErrNoRows:
			return &NotFound{"recipes", recipe.Id}
		default:
			return &InternalError{err.Error()}
		}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &InternalError{"recipe not updated"}
	}
	return r.createRecipeIngredients(recipe)
}

func (r RecipeRepository) createRecipeIngredients(recipe Recipe) error {
	for index, ing := range recipe.Ingredients {
		result, err := r.db.Exec(context.Background(), "INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount, unit, index) VALUES ($1, $2, $3, $4, $5)", recipe.Id, ing.Id, ing.Amount, ing.Unit, index)
		if err != nil {
			log.Println(err.Error())
			return &InternalError{err.Error()}
		}
		rowCnt := result.RowsAffected()
		if rowCnt != 1 {
			log.Println(err.Error())
			return &InternalError{"recipe ingredient not created"}
		}
	}
	return nil
}

func (r RecipeRepository) Delete(id int64) error {
	result, err := r.db.Exec(context.Background(), "DELETE FROM recipes WHERE id = $1", id)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &NotFound{"recipes", id}
	}
	return nil
}

func (r RecipeRepository) deleteRecipeIngredients(recipeId int64) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM recipe_ingredients WHERE recipe_id = $1", recipeId)
	if err != nil {
		return &InternalError{err.Error()}
	}
	return nil
}
