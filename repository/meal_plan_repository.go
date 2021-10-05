package repository

import (
	"database/sql"
	"log"
)

type MealPlanRepository struct {
	db *sql.DB
}

func NewMealPlanRepository() (*MealPlanRepository, error) {
	r := new(MealPlanRepository)
	var err error
	r.db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/cookbook")
	if err != nil {
		log.Println(err.Error())
	}
	return r, err
}

func (r MealPlanRepository) Get(id int64) (mealPlan MealPlan, e error) {
	err := r.db.QueryRow("SELECT * FROM meal_plans WHERE id = ?", id).Scan(&mealPlan.Id, &mealPlan.Name, &mealPlan.StartDate)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case sql.ErrNoRows:
			e = &NotFound{"meal_plans", id}
		default:
			e = &InternalError{err.Error()}
		}
	}
	mealPlan.Meals, e = r.getMealPlanMeals(id)
	return
}

func (r MealPlanRepository) getMealPlanMeals(id int64) (meals []int64, err error) {
	results, err := r.db.Query("SELECT meal_id FROM meal_plan_meals WHERE meal_plan_id = ? ORDER BY meal_plan_meals.index", id)
	if err != nil {
		return []int64{}, &InternalError{err.Error()}
	}
	for results.Next() {
		var mealId int64
		err = results.Scan(&mealId)
		if err != nil {
			log.Println(err.Error())
		}
		meals = append(meals, mealId)
	}
	return meals, nil
}

func (r MealPlanRepository) GetAll() (mealPlans []MealPlan, e error) {
	results, err := r.db.Query("SELECT * FROM meal_plans")
	if err != nil {
		return []MealPlan{}, &InternalError{err.Error()}
	}
	for results.Next() {
		var mealPlan MealPlan
		err = results.Scan(&mealPlan.Id, &mealPlan.Name, &mealPlan.StartDate)
		if err != nil {
			log.Println(err.Error())
		}
		mealPlan.Meals, err = r.getMealPlanMeals(mealPlan.Id)
		if err != nil {
			log.Println(err.Error())
		}
		mealPlans = append(mealPlans, mealPlan)
	}
	return mealPlans, nil
}

func (r MealPlanRepository) Create(mealPlan MealPlan) error {
	result, err := r.db.Exec("INSERT INTO meal_plans (name, start_date) VALUES (?, ?)", mealPlan.Name, mealPlan.StartDate)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return &InternalError{"meal not created"}
	}
	return r.createMealPlanMeals(mealPlan)
}

func (r MealPlanRepository) Update(mealPlan MealPlan) error {
	err := r.deleteMealPlanMeals(mealPlan.Id)
	if err != nil {
		return err
	}
	result, err := r.db.Exec("UPDATE meal_plans SET name = ?, start_date = ? WHERE id = ?", mealPlan.Name, mealPlan.StartDate, mealPlan.Id)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case sql.ErrNoRows:
			return &NotFound{"meal_plans", mealPlan.Id}
		default:
			return &InternalError{err.Error()}
		}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return &InternalError{"meal plan not updated"}
	}
	return r.createMealPlanMeals(mealPlan)
}

func (r MealPlanRepository) createMealPlanMeals(mealPlan MealPlan) error {
	for index, mealId := range mealPlan.Meals {
		result, err := r.db.Exec("INSERT INTO meal_plan_meals (meal_plan_id, meal_id, index) VALUES (?, ?, ?)", mealPlan.Id, mealId, index)
		if err != nil {
			return &InternalError{err.Error()}
		}
		rowCnt, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowCnt != 1 {
			return &InternalError{"meal plan meal not created"}
		}
	}
	return nil
}

func (r MealPlanRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM meal_plans WHERE id = ?", id)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return &InternalError{err.Error()}
	}
	if rowCnt != 1 {
		return &NotFound{"meal_plans", id}
	}
	return nil
}

func (r MealPlanRepository) deleteMealPlanMeals(mealPlanId int64) error {
	_, err := r.db.Exec("DELETE FROM meal_plan_meals WHERE meal_plan_id = ?", mealPlanId)
	if err != nil {
		return &InternalError{err.Error()}
	}
	return nil
}
