package service

import (
	"strings"

	"github.com/cookbook/repository"
)

type IngredientService interface {
	Get(int64) (Ingredient, error)
	GetAll() ([]Ingredient, error)
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

func (s ServiceImpl) Get(id int64) (ing Ingredient, e error) {
	ri, err := s.repo.Get(id)
	if err != nil {
		return Ingredient{}, handleError(err)
	}

	ing.Id = ri.Id
	ing.Name = ri.Name
	ing.Calories = ri.Calories
	ing.Amount = ri.Amount
	ing.Unit = ri.Unit
	ing.Fat = ri.Fat
	ing.Carbs = ri.Carbs
	ing.Protein = ri.Protein

	return ing, e
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

func (s ServiceImpl) GetAll() ([]Ingredient, error) {
	ri, err := s.repo.GetAll()
	if err != nil {
		return []Ingredient{}, handleError(err)
	}
	var ings = make([]Ingredient, len(ri))

	for index, i := range ri {
		ings[index].Id = i.Id
		ings[index].Name = i.Name
		ings[index].Calories = i.Calories
		ings[index].Amount = i.Amount
		ings[index].Unit = i.Unit
		ings[index].Fat = i.Fat
		ings[index].Carbs = i.Carbs
		ings[index].Protein = i.Protein
	}

	return ings, nil
}

func validateIngredient(i Ingredient) error {
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
