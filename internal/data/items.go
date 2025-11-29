package data

import (
	"context"
	"database/sql"
	"depopa/internal/validator"
	"errors"
	"github.com/lib/pq"
	"time"
)

type Item struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Condition   string   `json:"condition"`
	Description string   `json:"description"`
	Colors      []string `json:"colors"`
	Price       int      `json:"price"`
	Size        int      `json:"size"`
}

func ValidateItem(v *validator.Validator, item *Item) {
	v.Check(item.Name != "", "name", "must be provided")
	v.Check(len(item.Name) <= 200, "name", "must not be more than 200 bytes long")
	v.Check(item.Description != "", "description", "must be provided")
	v.Check(len(item.Description) <= 500, "description", "must not be more than 500 bytes long")
	v.Check(item.Condition != "", "condition", "must be provided")
	v.Check(len(item.Condition) <= 50, "condition", "must not be more than 50 bytes long")
	v.Check(item.Size != 0, "size", "must be provided")
	v.Check(item.Size < 100, "size", "must be less than 100")
	v.Check(item.Price > 0, "price", "must be greater than 0")
	v.Check(validator.Unique(item.Colors), "colors", "must not contain duplicate values")
}

type ItemModel struct {
	DB *sql.DB
}

func (i ItemModel) Insert(item *Item) error {

	query := `INSERT INTO items (name, condition, description, colors, price, size)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id`

	args := []interface{}{item.Name, item.Condition, item.Description, pq.Array(item.Colors), item.Price, item.Size}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return i.DB.QueryRowContext(ctx, query, args...).Scan(&item.ID)
}

func (i ItemModel) Get(id int64) (*Item, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id,name,condition,description,colors,price,size FROM items WHERE id = $1`

	var item Item

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := i.DB.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Condition,
		&item.Description,
		pq.Array(&item.Colors),
		&item.Price,
		&item.Size,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &item, nil
}
