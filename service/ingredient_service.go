package service

import (
	"fmt"
	"strings"

	"github.com/cookbook/repository"
)

type IngredientService interface {
	Get(int64) (Ingredient, error)
	GetAll() ([]Ingredient, error)
	GetList(ids []int64) ([]Ingredient, error)
	Create(Ingredient) error
	Update(Ingredient) error
	Delete(int64) error
}

type NotFound struct {
	message string
}

func (e *NotFound) Error() string {
	return e.message
}

type InternalError struct {
	message string
}

func (e *InternalError) Error() string {
	return e.message
}

type ValidationError struct {
	messages []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.messages, "\n")
}

type ServiceImpl struct {
	repo *repository.IngredientRepository
}

func NewIngredientService(r *repository.IngredientRepository) IngredientService {
	var s ServiceImpl
	s.repo = r
	return s
}

func (s ServiceImpl) Get(id int64) (Ingredient, error) {
	ri, err := s.repo.Get(id)
	if err != nil {
		return Ingredient{}, handleError(err)
	}
	ingredients := s.convertRepoModel(ri)
	if len(ingredients) != 1 {
		fmt.Println("Found multiple ingredients with the same ID")
		return Ingredient{}, &InternalError{message: "internal error, if the problem persists contact server admin"}
	}
	return ingredients[0], nil
}

func (s ServiceImpl) GetList(ids []int64) ([]Ingredient, error) {
	ri, err := s.repo.GetList(ids)
	if err != nil {
		fmt.Println(err.Error())
		return []Ingredient{}, handleError(err)
	}
	return s.convertRepoModel(ri...), nil
}

func (s ServiceImpl) GetAll() ([]Ingredient, error) {
	ri, err := s.repo.GetAll()
	if err != nil {
		fmt.Println(err.Error())
		return []Ingredient{}, handleError(err)
	}
	return s.convertRepoModel(ri...), nil
}

func (s ServiceImpl) Create(i Ingredient) (err error) {
	err = validateIngredient(i)
	if err != nil {
		return err
	}

	ri := repository.Ingredient{
		Name:     i.Name,
		Calories: i.Calories,
		Protein:  i.Protein,
		Carbs:    i.Carbs,
		Fat:      i.Fat,
		Amount:   i.Amount,
		Unit:     i.Unit,
	}

	err = s.repo.Create(ri)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s ServiceImpl) Update(i Ingredient) (err error) {
	err = validateIngredient(i)
	if err != nil {
		return err
	}

	ri := repository.Ingredient{
		Id:       i.Id,
		Name:     i.Name,
		Calories: i.Calories,
		Protein:  i.Protein,
		Carbs:    i.Carbs,
		Fat:      i.Fat,
		Amount:   i.Amount,
		Unit:     i.Unit,
	}

	err = s.repo.Update(ri)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s ServiceImpl) Delete(id int64) (err error) {

	err = s.repo.Delete(id)
	if err != nil {
		err = handleError(err)
	}
	return
}

func (s ServiceImpl) convertRepoModel(repoIngredients ...repository.Ingredient) []Ingredient {
	var ings = make([]Ingredient, len(repoIngredients))

	for index, i := range repoIngredients {
		ings[index] = Ingredient{
			Id:   i.Id,
			Name: i.Name,
			NutritionalValue: NutritionalValue{
				Quantity: Quantity{
					Amount: i.Amount,
					Unit:   i.Unit,
				},
				Calories: i.Calories,
				Fat:      i.Fat,
				Carbs:    i.Carbs,
				Protein:  i.Protein,
			},
		}
	}
	return ings
}

func validateIngredient(i Ingredient) error {
	var messages []string
	if i.Name == "" {
		messages = append(messages, "Ingredient name must be provided")
	}
	if !isUnitValid(i.Unit) {
		messages = append(messages, fmt.Sprintf("Invalid measurement unit: %s", i.Unit))
	}
	if i.Amount <= 0.0 {
		messages = append(messages, "Ingredient amount must be greater then 0")
	}
	if i.Calories < 0.0 {
		messages = append(messages, "Calories amount must be a positive value")
	}
	if i.Protein < 0.0 {
		messages = append(messages, "Protein amount must be a positive value")
	}
	if i.Carbs < 0.0 {
		messages = append(messages, "Carb amount must be a positive value")
	}
	if i.Fat < 0.0 {
		messages = append(messages, "Fat amount must be a positive value")
	}
	if len(messages) > 0 {
		return &ValidationError{messages: messages}
	}
	return nil
}

func handleError(e error) error {
	switch x := e.(type) {
	case *repository.NotFound:
		return &NotFound{x.Error()}
	default:
		return &InternalError{x.Error()}
	}
}
