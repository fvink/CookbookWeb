package repository

type Recipe struct {
	Id          int64
	Name        string
	Steps       string
	Ingredients []IngredientShort
}

type IngredientShort struct {
	Id     int64
	Amount float32
	Unit   string
}
