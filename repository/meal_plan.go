package repository

import "time"

type MealPlan struct {
	Id        int64
	Name      string
	StartDate time.Time
	Meals     [][]int64
}
