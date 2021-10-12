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
	mealPlans, err := s.convertRepoModel(rMealPlan)
	if err != nil {
		return MealPlanGet{}, handleError(err)
	}
	if len(mealPlans) != 1 {
		fmt.Println("Found multiple meal plans with the same ID")
		return MealPlanGet{}, &InternalError{message: "internal error, if the problem persists contact server admin"}
	}
	return mealPlans[0], nil
}

func (s MealPlanServiceImpl) GetAll() ([]MealPlanGet, error) {
	rMealPlans, err := s.repo.GetAll()
	if err != nil {
		return []MealPlanGet{}, handleError(err)
	}
	return s.convertRepoModel(rMealPlans...)
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

func (s MealPlanServiceImpl) convertRepoModel(repoMealPlans ...repository.MealPlan) ([]MealPlanGet, error) {
	var mealPlans = make([]MealPlanGet, len(repoMealPlans))
	usedMeals, err := s.getAllMeals(repoMealPlans...)
	if err != nil {
		return []MealPlanGet{}, handleError(err)
	}
	for index, rMealPlan := range repoMealPlans {
		mealPlans[index] = MealPlanGet{
			Id:          rMealPlan.Id,
			Name:        rMealPlan.Name,
			DateStarted: rMealPlan.StartDate,
		}
		mealPlans[index].Meals = make([][]MealGet, len(rMealPlan.Meals))
		for day, dayMeals := range rMealPlan.Meals {
			for _, mealId := range dayMeals {
				mealPlans[index].Meals[day] = append(mealPlans[index].Meals[day], (*usedMeals)[mealId])
			}
		}
	}
	return mealPlans, nil
}

func (s MealPlanServiceImpl) getAllMeals(mealPlans ...repository.MealPlan) (*map[int64]MealGet, error) {
	meals := make(map[int64]MealGet)
	var ids []int64
	for _, mealPlan := range mealPlans {
		for _, dayMeals := range mealPlan.Meals {
			for _, mealId := range dayMeals {
				if _, ok := meals[mealId]; !ok {
					meals[mealId] = MealGet{}
					ids = append(ids, mealId)
				}
			}
		}
	}
	rMeals, err := s.mealService.GetList(ids)
	if err != nil {
		return nil, err
	}
	for _, meal := range rMeals {
		meals[meal.Id] = meal
	}
	return &meals, nil
}

func validateMealPlan(meal MealPlanCreate) error {
	return nil
}
