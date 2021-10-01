package service

type Ingredient struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	NutritionalValue `json:"nutritional_value"`
}
