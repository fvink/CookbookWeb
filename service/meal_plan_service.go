package service

import "github.com/cookbook/repository"

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
	for _, mealId := range rMealPlan.Meals {
		rMeal, _ := s.mealService.Get(mealId)
		mealPlan.Meals = append(mealPlan.Meals, rMeal)
	}
	return
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
	rMealPlan.Meals = append(rMealPlan.Meals, mealPlan.Meals...)
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
	rMealPlan.Meals = append(rMealPlan.Meals, mealPlan.Meals...)
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
		for _, mealId := range rMealPlan.Meals {
			rMeal, _ := s.mealService.Get(mealId)
			if err != nil {

			}
			mealPlans[index].Meals = append(mealPlans[index].Meals, rMeal)
		}
	}

	return mealPlans, nil
}

func validateMealPlan(meal MealPlanCreate) error {
	return nil
}
