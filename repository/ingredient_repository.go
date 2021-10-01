package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type IngredientRepository struct {
	db *sql.DB
}

type NotFound struct {
	Table string
	Id    int64
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%s with id %d doesn't exist", e.Table, e.Id)
}

type InternalError struct {
	message string
}

func (e *InternalError) Error() string {
	return e.message
}

func NewIngredientRepository() (*IngredientRepository, error) {
	r := new(IngredientRepository)
	var err error
	r.db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/cookbook")
	if err != nil {
		log.Println(err.Error())
	}
	return r, err
}

func (r IngredientRepository) Get(id int64) (i Ingredient, e error) {
	err := r.db.QueryRow("SELECT * FROM ingredients WHERE id = ?", id).Scan(&i.Id, &i.Name, &i.Calories, &i.Protein, &i.Carbs, &i.Fat, &i.Amount, &i.Unit)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case sql.ErrNoRows:
			e = &NotFound{"ingredients", id}
		default:
			e = &InternalError{err.Error()}
		}
	}
	return
}

func (r IngredientRepository) Create(i Ingredient) error {
	result, err := r.db.Exec("INSERT INTO ingredients (name, calories, protein, carbs, fat, amount, unit) VALUES (?, ?, ?, ?, ?, ?, ?)", i.Name, i.Calories, i.Protein, i.Carbs, i.Fat, i.Amount, i.Unit)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return &InternalError{"ingredient not created"}
	}
	return nil
}

func (r IngredientRepository) Update(i Ingredient) error {
	result, err := r.db.Exec("UPDATE ingredients SET name = ?, calories = ?, protein = ?, carbs = ?, fat = ?, amount = ?, unit = ? WHERE id = ?", i.Name, i.Calories, i.Protein, i.Carbs, i.Fat, i.Amount, i.Unit, i.Id)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case sql.ErrNoRows:
			return &NotFound{"ingredients", i.Id}
		default:
			return &InternalError{err.Error()}
		}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return &InternalError{"ingredient not updated"}
	}
	return nil
}

func (r IngredientRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM ingredients WHERE id = ?", id)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return &InternalError{err.Error()}
	}
	if rowCnt != 1 {
		return &NotFound{"ingredients", id}
	}
	return nil
}

func (r IngredientRepository) GetAll() (ingredients []Ingredient, err error) {
	results, err := r.db.Query("SELECT * FROM ingredients")
	if err != nil {
		return []Ingredient{}, &InternalError{err.Error()}
	}
	for results.Next() {
		var i Ingredient
		err = results.Scan(&i.Id, &i.Name, &i.Calories, &i.Protein, &i.Carbs, &i.Fat, &i.Amount, &i.Unit)
		if err != nil {
			log.Println(err.Error())
		}
		ingredients = append(ingredients, i)
	}
	return ingredients, nil
}

func (r IngredientRepository) Close() {
	r.db.Close()
}
