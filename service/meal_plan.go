package service

import "time"

type MealPlanGet struct {
	Id          int64
	Name        string
	DateStarted time.Time
	Meals       []MealGet
}

type MealPlanCreate struct {
	Id          int64
	Name        string
	DateStarted time.Time
	Meals       []int64
}
