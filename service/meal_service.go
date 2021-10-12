package service

import (
	"fmt"

	"github.com/cookbook/repository"
)

type MealService interface {
	Get(int64) (MealGet, error)
	GetList([]int64) ([]MealGet, error)
	GetAll() ([]MealGet, error)
	Create(MealCreate) error
	Update(MealCreate) error
	Delete(int64) error
}

type MealServiceImpl struct {
	repo       *repository.MealRepository
	rcpService RecipeService
}

func NewMealService(r *repository.MealRepository, rs RecipeService) MealService {
	return MealServiceImpl{
		repo:       r,
		rcpService: rs,
	}
}

func (s MealServiceImpl) Get(id int64) (meal MealGet, err error) {
	rMeal, err := s.repo.Get(id)
	if err != nil {
		return MealGet{}, handleError(err)
	}

	meal = MealGet{
		Id:   rMeal.Id,
		Name: rMeal.Name,
	}
	for _, recipeId := range rMeal.Recipes {
		rRecipe, _ := s.rcpService.Get(recipeId)
		meal.Recipes = append(meal.Recipes, rRecipe)
	}
	return
}

func (s MealServiceImpl) Create(meal MealCreate) (err error) {
	err = validateMeal(meal)
	if err != nil {
		return err
	}
	rMeal := repository.Meal{
		Name: meal.Name,
	}
	rMeal.Recipes = append(rMeal.Recipes, meal.Recipes...)
	err = s.repo.Create(rMeal)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s MealServiceImpl) Update(meal MealCreate) (err error) {
	err = validateMeal(meal)
	if err != nil {
		return err
	}
	rMeal := repository.Meal{
		Id:   meal.Id,
		Name: meal.Name,
	}
	rMeal.Recipes = append(rMeal.Recipes, meal.Recipes...)
	err = s.repo.Update(rMeal)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s MealServiceImpl) Delete(id int64) (err error) {
	err = s.repo.Delete(id)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s MealServiceImpl) GetAll() ([]MealGet, error) {
	rMeals, err := s.repo.GetAll()
	if err != nil {
		return []MealGet{}, handleError(err)
	}
	var meals = make([]MealGet, len(rMeals))

	for index, rMeal := range rMeals {
		meals[index] = MealGet{
			Id:   rMeal.Id,
			Name: rMeal.Name,
		}
		for _, recipeId := range rMeal.Recipes {
			rRecipe, _ := s.rcpService.Get(recipeId)
			if err != nil {

			}
			meals[index].Recipes = append(meals[index].Recipes, rRecipe)
		}
	}

	return meals, nil
}

func (s MealServiceImpl) GetList(ids []int64) ([]MealGet, error) {
	rMeals, err := s.repo.GetList(ids)
	if err != nil {
		fmt.Println(err.Error())
		return []MealGet{}, handleError(err)
	}
	var meals = make([]MealGet, len(rMeals))

	for index, rMeal := range rMeals {
		meals[index] = MealGet{
			Id:   rMeal.Id,
			Name: rMeal.Name,
		}
		for _, recipeId := range rMeal.Recipes {
			rRecipe, _ := s.rcpService.Get(recipeId)
			if err != nil {

			}
			meals[index].Recipes = append(meals[index].Recipes, rRecipe)
		}
	}

	return meals, nil
}

func validateMeal(meal MealCreate) error {
	return nil
}
