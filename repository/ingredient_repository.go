package repository

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type IngredientRepository struct {
	db *pgxpool.Pool
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

func NewIngredientRepository(dbConn *pgxpool.Pool) *IngredientRepository {
	r := new(IngredientRepository)
	r.db = dbConn
	return r
}

func (r IngredientRepository) Get(id int64) (i Ingredient, e error) {
	err := r.db.QueryRow(context.Background(), "SELECT * FROM ingredients WHERE id = $1", id).Scan(&i.Id, &i.Name, &i.Calories, &i.Protein, &i.Carbs, &i.Fat, &i.Amount, &i.Unit)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case pgx.ErrNoRows:
			e = &NotFound{"ingredients", id}
		default:
			e = &InternalError{err.Error()}
		}
	}
	return
}

func (r IngredientRepository) Create(i Ingredient) error {
	result, err := r.db.Exec(context.Background(), "INSERT INTO ingredients (name, calories, protein, carbs, fat, amount, unit) VALUES ($1, $2, $3, $4, $5, $6, $7)", i.Name, i.Calories, i.Protein, i.Carbs, i.Fat, i.Amount, i.Unit)
	if err != nil {
		log.Println(err.Error())
		return &InternalError{err.Error()}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &InternalError{"ingredient not created"}
	}
	return nil
}

func (r IngredientRepository) Update(i Ingredient) error {
	result, err := r.db.Exec(context.Background(), "UPDATE ingredients SET name = $1, calories = $2, protein = $3, carbs = $4, fat = $5, amount = $6, unit = $7 WHERE id = $8", i.Name, i.Calories, i.Protein, i.Carbs, i.Fat, i.Amount, i.Unit, i.Id)
	if err != nil {
		log.Println(err.Error())
		switch err {
		case pgx.ErrNoRows:
			return &NotFound{"ingredients", i.Id}
		default:
			return &InternalError{err.Error()}
		}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &InternalError{"ingredient not updated"}
	}
	return nil
}

func (r IngredientRepository) Delete(id int64) error {
	result, err := r.db.Exec(context.Background(), "DELETE FROM ingredients WHERE id = $1", id)
	if err != nil {
		return &InternalError{err.Error()}
	}
	rowCnt := result.RowsAffected()
	if rowCnt != 1 {
		return &NotFound{"ingredients", id}
	}
	return nil
}

func (r IngredientRepository) GetAll() (ingredients []Ingredient, err error) {
	results, err := r.db.Query(context.Background(), "SELECT * FROM ingredients")
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

func (r IngredientRepository) GetList(ids []int64) (ingredients []Ingredient, err error) {
	results, err := r.db.Query(context.Background(), "SELECT * FROM ingredients WHERE id IN ("+JoinIds(ids)+")")
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

func JoinIds(ids []int64) string {
	if len(ids) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for i, id := range ids {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(strconv.FormatInt(id, 10))
	}
	return buf.String()
}
