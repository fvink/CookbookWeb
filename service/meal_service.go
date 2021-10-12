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

func (s MealServiceImpl) Get(id int64) (MealGet, error) {
	rMeal, err := s.repo.Get(id)
	if err != nil {
		return MealGet{}, handleError(err)
	}
	meals, err := s.convertRepoModel(rMeal)
	if err != nil {
		return MealGet{}, handleError(err)
	}
	if len(meals) != 1 {
		fmt.Println("Found multiple meals with the same ID")
		return MealGet{}, &InternalError{message: "internal error, if the problem persists contact server admin"}
	}
	return meals[0], nil
}

func (s MealServiceImpl) GetList(ids []int64) ([]MealGet, error) {
	rMeals, err := s.repo.GetList(ids)
	if err != nil {
		fmt.Println(err.Error())
		return []MealGet{}, handleError(err)
	}
	return s.convertRepoModel(rMeals...)
}

func (s MealServiceImpl) GetAll() ([]MealGet, error) {
	rMeals, err := s.repo.GetAll()
	if err != nil {
		fmt.Println(err.Error())
		return []MealGet{}, handleError(err)
	}
	return s.convertRepoModel(rMeals...)
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

func (s MealServiceImpl) convertRepoModel(repoMeals ...repository.Meal) ([]MealGet, error) {
	var meals = make([]MealGet, len(repoMeals))
	usedRecipes, err := s.getAllRecipes(repoMeals...)
	if err != nil {
		return []MealGet{}, handleError(err)
	}
	for index, rMeal := range repoMeals {
		meals[index] = MealGet{
			Id:   rMeal.Id,
			Name: rMeal.Name,
		}
		for _, recipeId := range rMeal.Recipes {
			meals[index].Recipes = append(meals[index].Recipes, (*usedRecipes)[recipeId])
		}
	}
	return meals, nil
}

func (s MealServiceImpl) getAllRecipes(meals ...repository.Meal) (*map[int64]RecipeGet, error) {
	recipes := make(map[int64]RecipeGet)
	var ids []int64
	for _, meal := range meals {
		for _, recipeId := range meal.Recipes {
			if _, ok := recipes[recipeId]; !ok {
				recipes[recipeId] = RecipeGet{}
				ids = append(ids, recipeId)
			}
		}
	}
	rRecipes, err := s.rcpService.GetList(ids)
	if err != nil {
		return nil, err
	}
	for _, recipe := range rRecipes {
		recipes[recipe.Id] = recipe
	}
	return &recipes, nil
}

func validateMeal(meal MealCreate) error {
	return nil
}
