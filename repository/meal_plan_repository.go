package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MealPlanRepository struct {
	db *pgxpool.Pool
}

func NewMealPlanRepository(dbConn *pgxpool.Pool) *MealPlanRepository {
	r := new(MealPlanRepository)
	r.db = dbConn
	return r
}

func (r MealPlanRepository) Get(id int64) (mealPlan MealPlan, e error) {
	err := r.db.QueryRow(context.Background(), "SELECT * FROM meal_plans WHERE id = $1", id).Scan(&mealPlan.Id, &mealPlan.Name, &mealPlan.StartDate)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case pgx.ErrNoRows:
			e = &NotFound{"meal_plans", id}
		default:
			e = &InternalError{err.Error()}
		}
	}
	meals, e := r.getMealPlanMeals([]int64{id})
	var ok bool
	if mealPlan.Meals, ok = meals[id]; !ok {
		mealPlan.Meals = [][]int64{}
	}
	return
}

func (r MealPlanRepository) GetAll() (mealPlans []MealPlan, e error) {
	results, err := r.db.Query(context.Background(), "SELECT * FROM meal_plans")
	if err != nil {
		return []MealPlan{}, &InternalError{err.Error()}
	}
	meals, err := r.getAllMealPlanMeals()
	if err != nil {
		return []MealPlan{}, &InternalError{err.Error()}
	}
	for results.Next() {
		var mealPlan MealPlan
		err = results.Scan(&mealPlan.Id, &mealPlan.Name, &mealPlan.StartDate)
		if err != nil {
			log.Println(err.Error())
		}
		var ok bool
		if mealPlan.Meals, ok = meals[mealPlan.Id]; !ok {
			mealPlan.Meals = [][]int64{}
		}
		mealPlans = append(mealPlans, mealPlan)
	}
	return mealPlans, nil
}

func (r MealPlanRepository) getMealPlanMeals(ids []int64) (meals map[int64][][]int64, err error) {
	meals = make(map[int64][][]int64)
	results, err := r.db.Query(context.Background(), "SELECT meal_plan_id, meal_id, day FROM meal_plan_meals WHERE meal_plan_id IN ("+JoinIds(ids)+") ORDER BY meal_plan_meals.day DESC, meal_plan_meals.index")
	if err != nil {
		return nil, &InternalError{err.Error()}
	}
	for results.Next() {
		var mealPlanId, mealId, day int64
		err = results.Scan(&mealPlanId, &mealId, &day)
		if err != nil {
			log.Println(err.Error())
		}
		if _, ok := meals[mealPlanId]; !ok {
			meals[mealPlanId] = make([][]int64, day+1)
		}
		meals[mealPlanId][day] = append(meals[mealPlanId][day], mealId)
	}
	return meals, nil
}

func (r MealPlanRepository) getAllMealPlanMeals() (meals map[int64][][]int64, err error) {
	meals = make(map[int64][][]int64)
	results, err := r.db.Query(context.Background(), "SELECT meal_plan_id, meal_id, day FROM meal_plan_meals ORDER BY meal_plan_meals.day DESC, meal_plan_meals.index")
	if err != nil {
		return nil, &InternalError{err.Error()}
	}
	for results.Next() {
		var mealPlanId, mealId, day int64
		err = results.Scan(&mealPlanId, &mealId, &day)
		if err != nil {
			log.Println(err.Error())
		}
		if _, ok := meals[mealPlanId]; !ok {
			meals[mealPlanId] = make([][]int64, day+1)
		}
		meals[mealPlanId][day] = append(meals[mealPlanId][day], mealId)
	}
	return meals, nil
}

func (r MealPlanRepository) Create(mealPlan MealPlan) error {
	err := r.db.QueryRow(context.Background(), "INSERT INTO meal_plans (name, start_date, days) VALUES ($1, $2, $3) RETURNING id", mealPlan.Name, mealPlan.StartDate, len(mealPlan.Meals)).Scan(&mealPlan.Id)
	if err != nil {
		log.Println(err.Error())
		return &InternalError{err.Error()}
	}
	return r.createMealPlanMeals(mealPlan)
}

func (r MealPlanRepository) Update(mealPlan MealPlan) error {
	err := r.deleteMealPlanMeals(mealPlan.Id)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(context.Background(), "UPDATE meal_plans SET name = $1, start_date = $2, days = $3 WHERE id = $4", mealPlan.Name, mealPlan.StartDate, len(mealPlan.Meals), mealPlan.Id)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case pgx.ErrNoRows:
			return &NotFound{"meal_plans", mealPlan.Id}
		default:
			return &InternalError{err.Error()}
		}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &InternalError{"meal plan not updated"}
	}
	return r.createMealPlanMeals(mealPlan)
}

func (r MealPlanRepository) createMealPlanMeals(mealPlan MealPlan) error {
	for day, dayMeals := range mealPlan.Meals {
		for index, mealId := range dayMeals {
			result, err := r.db.Exec(context.Background(), "INSERT INTO meal_plan_meals (meal_plan_id, meal_id, day, index) VALUES ($1, $2, $3, $4)", mealPlan.Id, mealId, day, index)
			if err != nil {
				log.Println(err.Error())
				return &InternalError{err.Error()}
			}
			rowCnt := result.RowsAffected()
			if rowCnt != 1 {
				log.Println(err.Error())
				return &InternalError{"meal plan meal not created"}
			}
		}
	}
	return nil
}

func (r MealPlanRepository) Delete(id int64) error {
	result, err := r.db.Exec(context.Background(), "DELETE FROM meal_plans WHERE id = $1", id)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &NotFound{"meal_plans", id}
	}
	return nil
}

func (r MealPlanRepository) deleteMealPlanMeals(mealPlanId int64) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM meal_plan_meals WHERE meal_plan_id = $1", mealPlanId)
	if err != nil {
		return &InternalError{err.Error()}
	}
	return nil
}
