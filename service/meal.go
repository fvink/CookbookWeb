package service

type MealGet struct {
	Id      int64       `json:"id"`
	Name    string      `json:"name"`
	Recipes []RecipeGet `json:"recipes"`
}

type MealCreate struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Recipes []int64 `json:"recipes"`
}
