package service

import (
	"fmt"
	"log"

	"github.com/cookbook/repository"
)

type RecipeService interface {
	Get(int64) (RecipeGet, error)
	GetList([]int64) ([]RecipeGet, error)
	GetAll() ([]RecipeGet, error)
	Create(RecipeCreate) error
	Update(RecipeCreate) error
	Delete(int64) error
}

type RecipeServiceImpl struct {
	repo       *repository.RecipeRepository
	ingService IngredientService
}

func NewRecipeService(r *repository.RecipeRepository, is IngredientService) RecipeService {
	return RecipeServiceImpl{
		repo:       r,
		ingService: is,
	}
}

func (s RecipeServiceImpl) Get(id int64) (recipe RecipeGet, err error) {
	rRecipe, err := s.repo.Get(id)
	if err != nil {
		return RecipeGet{}, handleError(err)
	}
	recipes, err := s.convertRepoModel(rRecipe)
	if err != nil {
		return RecipeGet{}, handleError(err)
	}
	if len(recipes) != 1 {
		fmt.Println("Found multiple recipes with the same ID")
		return RecipeGet{}, &InternalError{message: "internal error, if the problem persists contact server admin"}
	}
	return recipes[0], nil
}

func (s RecipeServiceImpl) GetList(ids []int64) ([]RecipeGet, error) {
	rRecipes, err := s.repo.GetList(ids)
	if err != nil {
		return []RecipeGet{}, handleError(err)
	}
	return s.convertRepoModel(rRecipes...)
}

func (s RecipeServiceImpl) GetAll() ([]RecipeGet, error) {
	rRecipes, err := s.repo.GetAll()
	if err != nil {
		return []RecipeGet{}, handleError(err)
	}
	return s.convertRepoModel(rRecipes...)
}

func (s RecipeServiceImpl) Create(recipe RecipeCreate) (err error) {
	err = validateRecipe(recipe)
	if err != nil {
		return err
	}
	rRecipe := repository.Recipe{
		Name:  recipe.Name,
		Steps: recipe.Steps,
	}
	for _, ing := range recipe.Ingredients {
		rRecipe.Ingredients = append(rRecipe.Ingredients, repository.IngredientShort{
			Id:     ing.Id,
			Amount: ing.Amount,
			Unit:   ing.Unit,
		})
	}
	err = s.repo.Create(rRecipe)
	if err != nil {
		err = handleError(err)
		log.Println(err.Error())
	}
	return
}

func (s RecipeServiceImpl) Update(recipe RecipeCreate) (err error) {
	err = validateRecipe(recipe)
	if err != nil {
		return err
	}
	rRecipe := repository.Recipe{
		Id:    recipe.Id,
		Name:  recipe.Name,
		Steps: recipe.Steps,
	}
	for _, ing := range recipe.Ingredients {
		rRecipe.Ingredients = append(rRecipe.Ingredients, repository.IngredientShort{
			Id:     ing.Id,
			Amount: ing.Amount,
			Unit:   ing.Unit,
		})
	}
	err = s.repo.Update(rRecipe)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s RecipeServiceImpl) Delete(id int64) (err error) {
	err = s.repo.Delete(id)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s RecipeServiceImpl) convertRepoModel(repoRecipes ...repository.Recipe) ([]RecipeGet, error) {
	var recipes = make([]RecipeGet, len(repoRecipes))
	usedIngredients, err := s.getAllIngredients(repoRecipes...)
	if err != nil {
		return []RecipeGet{}, handleError(err)
	}
	for index, rRecipe := range repoRecipes {
		recipes[index] = RecipeGet{
			Id:    rRecipe.Id,
			Name:  rRecipe.Name,
			Steps: rRecipe.Steps,
		}
		for _, ing := range rRecipe.Ingredients {
			rIng := (*usedIngredients)[ing.Id]
			rIng, err = scaleIngredient(rIng, Quantity{Amount: ing.Amount, Unit: ing.Unit})
			if err != nil {

			}
			recipes[index].Ingredients = append(recipes[index].Ingredients, rIng)
			recipes[index].Calories += rIng.Calories
			recipes[index].Protein += rIng.Protein
			recipes[index].Carbs += rIng.Carbs
			recipes[index].Fat += rIng.Fat
		}
	}
	return recipes, nil
}

func (s RecipeServiceImpl) getAllIngredients(recipes ...repository.Recipe) (*map[int64]Ingredient, error) {
	ingredients := make(map[int64]Ingredient)
	var ids []int64
	for _, recipe := range recipes {
		for _, ing := range recipe.Ingredients {
			if _, ok := ingredients[ing.Id]; !ok {
				ingredients[ing.Id] = Ingredient{}
				ids = append(ids, ing.Id)
			}
		}
	}
	rIngredients, err := s.ingService.GetList(ids)
	if err != nil {
		return nil, err
	}
	for _, ing := range rIngredients {
		ingredients[ing.Id] = ing
	}
	return &ingredients, nil
}

func scaleIngredient(i Ingredient, finalQuantity Quantity) (Ingredient, error) {
	unitScale, err := ConvertUnit(i.Unit, finalQuantity.Unit)
	if err == nil {
		nutritionScale := finalQuantity.Amount / (i.Amount * unitScale)
		i.Quantity = finalQuantity
		i.Calories *= nutritionScale
		i.Protein *= nutritionScale
		i.Carbs *= nutritionScale
		i.Fat *= nutritionScale
	}
	return i, err
}

func validateRecipe(recipe RecipeCreate) error {
	var messages []string
	if recipe.Name == "" {
		messages = append(messages, "Recipe name must be provided")
	}
	for _, ingredient := range recipe.Ingredients {
		if !isUnitValid(ingredient.Unit) {
			messages = append(messages, fmt.Sprintf("Invalid measurement unit %s for %d", ingredient.Unit, ingredient.Id))
		}
		if ingredient.Amount <= 0.0 {
			messages = append(messages, fmt.Sprintf("Ingredient amount must be greater then 0 for %d", ingredient.Id))
		}
	}
	if len(messages) > 0 {
		return &ValidationError{messages: messages}
	}
	return nil
}
