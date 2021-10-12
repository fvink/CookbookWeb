package service

import (
	"fmt"

	"github.com/cookbook/repository"
)

type MealPlanService interface {
	Get(int64) (MealPlanGet, error)
	GetAll() ([]MealPlanGet, error)
	Create(MealPlanCreate) error
	Update(MealPlanCreate) error
	Delete(int64) error
}

type MealPlanServiceImpl struct {
	repo        *repository.MealPlanRepository
	mealService MealService
}

func NewMealPlanService(r *repository.MealPlanRepository, ms MealService) MealPlanService {
	return MealPlanServiceImpl{
		repo:        r,
		mealService: ms,
	}
}

func (s MealPlanServiceImpl) Get(id int64) (mealPlan MealPlanGet, err error) {
	rMealPlan, err := s.repo.Get(id)
	if err != nil {
		return MealPlanGet{}, handleError(err)
	}

	mealPlan = MealPlanGet{
		Id:          rMealPlan.Id,
		Name:        rMealPlan.Name,
		DateStarted: rMealPlan.StartDate,
	}
	meals := make(map[int64]MealGet)
	for _, dayMeals := range rMealPlan.Meals {
		for _, mealId := range dayMeals {
			meals[mealId] = MealGet{}
		}
	}
	err = s.getMeals(&meals)
	if err != nil {
		return MealPlanGet{}, handleError(err)
	}
	mealPlan.Meals = make([][]MealGet, len(rMealPlan.Meals))
	for day, dayMeals := range rMealPlan.Meals {
		for _, mealId := range dayMeals {
			mealPlan.Meals[day] = append(mealPlan.Meals[day], meals[mealId])
		}
	}
	return
}

func (s MealPlanServiceImpl) getMeals(meals *map[int64]MealGet) error {
	var ids []int64
	for id, _ := range *meals {
		ids = append(ids, id)
	}
	rMeals, err := s.mealService.GetList(ids)
	if err != nil {
		return err
	}
	for _, meal := range rMeals {
		(*meals)[meal.Id] = meal
	}
	return nil
}

func (s MealPlanServiceImpl) Create(mealPlan MealPlanCreate) (err error) {
	err = validateMealPlan(mealPlan)
	if err != nil {
		return err
	}
	rMealPlan := repository.MealPlan{
		Name:      mealPlan.Name,
		StartDate: mealPlan.DateStarted,
	}
	rMealPlan.Meals = mealPlan.Meals
	err = s.repo.Create(rMealPlan)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s MealPlanServiceImpl) Update(mealPlan MealPlanCreate) (err error) {
	err = validateMealPlan(mealPlan)
	if err != nil {
		return err
	}
	rMealPlan := repository.MealPlan{
		Id:        mealPlan.Id,
		Name:      mealPlan.Name,
		StartDate: mealPlan.DateStarted,
	}
	rMealPlan.Meals = mealPlan.Meals
	err = s.repo.Update(rMealPlan)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s MealPlanServiceImpl) Delete(id int64) (err error) {
	err = s.repo.Delete(id)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s MealPlanServiceImpl) GetAll() ([]MealPlanGet, error) {
	rMealPlans, err := s.repo.GetAll()
	if err != nil {
		return []MealPlanGet{}, handleError(err)
	}
	var mealPlans = make([]MealPlanGet, len(rMealPlans))

	for index, rMealPlan := range rMealPlans {
		mealPlans[index] = MealPlanGet{
			Id:          rMealPlan.Id,
			Name:        rMealPlan.Name,
			DateStarted: rMealPlan.StartDate,
		}
		mealPlans[index].Meals = make([][]MealGet, len(rMealPlan.Meals))
		for day, dayMeals := range rMealPlan.Meals {
			for _, mealId := range dayMeals {
				fmt.Println(mealId)
				rMeal, _ := s.mealService.Get(mealId)
				if err != nil {

				}
				mealPlans[index].Meals[day] = append(mealPlans[index].Meals[day], rMeal)
			}
		}
	}

	return mealPlans, nil
}

func validateMealPlan(meal MealPlanCreate) error {
	return nil
}
