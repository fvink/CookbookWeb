package service

type IngredientShort struct {
	Id     int64   `json:"id"`
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}

type RecipeCreate struct {
	Id          int64             `json:"id"`
	Name        string            `json:"name"`
	Steps       string            `json:"steps"`
	Ingredients []IngredientShort `json:"ingredients"`
}

type RecipeGet struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Calories    float32      `json:"calories"`
	Protein     float32      `json:"protein"`
	Carbs       float32      `json:"carbs"`
	Fat         float32      `json:"fat"`
	Steps       string       `json:"steps"`
	Ingredients []Ingredient `json:"ingredients"`
}
