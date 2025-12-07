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
	Description string   `json:"description,omitempty"`
	Colors      []string `json:"colors,omitempty"`
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

func (i ItemModel) GetAll() ([]*Item, error) {
	query := `SELECT id,name,condition,price,size FROM items`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := i.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var items []*Item

	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Condition,
			&item.Price,
			&item.Size,
		)

		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

func (i ItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM items WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := i.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (i ItemModel) Update(item *Item) error {
	query := `UPDATE items
	SET name = $1, condition = $2, description = $3, colors = $4, price = $5, size = $6

	WHERE id = $7
`

	args := []interface{}{
		item.Name,
		item.Condition,
		item.Description,
		pq.Array(item.Colors),
		item.Price,
		item.Size,
		item.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := i.DB.ExecContext(ctx, query, args)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}

	return nil
}
