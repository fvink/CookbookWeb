package service

import "github.com/cookbook/repository"

type RecipeService interface {
	Get(int64) (RecipeGet, error)
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

	recipe = RecipeGet{
		Id:    rRecipe.Id,
		Name:  rRecipe.Name,
		Steps: rRecipe.Steps,
	}
	for _, ing := range rRecipe.Ingredients {
		rIng, _ := s.ingService.Get(ing.Id)
		rIng, err = scaleIngredient(rIng, Quantity{Amount: ing.Amount, Unit: ing.Unit})
		if err != nil {

		}
		recipe.Ingredients = append(recipe.Ingredients, rIng)
		recipe.Calories += rIng.Calories
		recipe.Protein += rIng.Protein
		recipe.Carbs += rIng.Carbs
		recipe.Fat += rIng.Fat
	}
	return
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

func (s RecipeServiceImpl) GetAll() ([]RecipeGet, error) {
	rRecipes, err := s.repo.GetAll()
	if err != nil {
		return []RecipeGet{}, handleError(err)
	}
	var recipes = make([]RecipeGet, len(rRecipes))

	for index, rRecipe := range rRecipes {
		recipes[index] = RecipeGet{
			Id:    rRecipe.Id,
			Name:  rRecipe.Name,
			Steps: rRecipe.Steps,
		}
		for _, ing := range rRecipe.Ingredients {
			rIng, _ := s.ingService.Get(ing.Id)
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
	return nil
}
