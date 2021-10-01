package repository

type Ingredient struct {
	Id       int64
	Name     string
	Calories float32
	Protein  float32
	Carbs    float32
	Fat      float32
	Amount   float32
	Unit     string
}
