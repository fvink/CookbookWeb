package service

import "time"

type MealPlanGet struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	DateStarted time.Time   `json:"date_started"`
	Meals       [][]MealGet `json:"meals"`
}

type MealPlanCreate struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	DateStarted time.Time `json:"date_started"`
	Meals       [][]int64 `json:"meals"`
}
