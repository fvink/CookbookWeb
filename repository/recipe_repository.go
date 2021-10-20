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
			return Recipe{}, &NotFound{"recipes", id}
		default:
			return Recipe{}, &InternalError{err.Error()}
		}
	}
	ingredients, err := r.getRecipeIngredientsByIds([]int64{id})
	if err != nil {
		return Recipe{}, err
	}
	recipe.Ingredients = ingredients[id]
	return recipe, nil
}

func (r RecipeRepository) GetAll() (recipes []Recipe, err error) {
	results, err := r.db.Query(context.Background(), "SELECT * FROM recipes")
	if err != nil {
		return []Recipe{}, &InternalError{err.Error()}
	}
	recipeIngredients, err := r.getAllRecipeIngredients()
	if err != nil {
		return []Recipe{}, err
	}
	return r.parseRecipeRows(results, recipeIngredients), nil
}

func (r RecipeRepository) GetList(ids []int64) (recipes []Recipe, err error) {
	results, err := r.db.Query(context.Background(), "SELECT * FROM recipes WHERE id IN ("+JoinIds(ids)+")")
	if err != nil {
		return []Recipe{}, &InternalError{err.Error()}
	}
	recipeIngredients, err := r.getRecipeIngredientsByIds(ids)
	if err != nil {
		return []Recipe{}, err
	}
	return r.parseRecipeRows(results, recipeIngredients), nil
}

func (r RecipeRepository) parseRecipeRows(rows pgx.Rows, recipeIngredients map[int64][]IngredientShort) (recipes []Recipe) {
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(&recipe.Id, &recipe.Name, &recipe.Steps)
		if err != nil {
			log.Println(err.Error())
		}
		recipe.Ingredients = recipeIngredients[recipe.Id]
		if err != nil {
			log.Println(err.Error())
		}
		recipes = append(recipes, recipe)
	}
	return recipes
}

func (r RecipeRepository) getRecipeIngredientsByIds(recipeIds []int64) (map[int64][]IngredientShort, error) {
	return r.getRecipeIngredients("SELECT recipe_id, ingredient_id, amount, unit FROM recipe_ingredients WHERE recipe_id IN (" + JoinIds(recipeIds) + ") ORDER BY recipe_id, recipe_ingredients.index")
}

func (r RecipeRepository) getAllRecipeIngredients() (map[int64][]IngredientShort, error) {
	return r.getRecipeIngredients("SELECT recipe_id, ingredient_id, amount, unit FROM recipe_ingredients ORDER BY recipe_id, recipe_ingredients.index")
}

func (r RecipeRepository) getRecipeIngredients(query string) (map[int64][]IngredientShort, error) {
	ingredients := make(map[int64][]IngredientShort)
	results, err := r.db.Query(context.Background(), query)
	if err != nil {
		log.Println(err.Error())
		return nil, &InternalError{err.Error()}
	}
	for results.Next() {
		var i IngredientShort
		var recipeId int64
		err := results.Scan(&recipeId, &i.Id, &i.Amount, &i.Unit)
		if err != nil {
			log.Println(err.Error())
			return nil, &InternalError{err.Error()}
		}
		if _, ok := ingredients[recipeId]; !ok {
			ingredients[recipeId] = make([]IngredientShort, 0)
		}
		ingredients[recipeId] = append(ingredients[recipeId], i)
	}
	return ingredients, nil
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
