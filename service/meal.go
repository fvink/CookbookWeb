package service

type MealGet struct {
	Id      int64
	Name    string
	Recipes []RecipeGet
}

type MealCreate struct {
	Id      int64
	Name    string
	Recipes []int64
}
