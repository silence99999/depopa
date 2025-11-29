package main

import (
	"depopa/internal/data"
	"depopa/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string   `json:"name"`
		Condition   string   `json:"condition"`
		Description string   `json:"description"`
		Colors      []string `json:"colors"`
		Price       int      `json:"price"`
		Size        int      `json:"size"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	item := &data.Item{
		Name:        input.Name,
		Condition:   input.Condition,
		Description: input.Description,
		Colors:      input.Colors,
		Price:       input.Price,
		Size:        input.Size,
	}

	v := validator.New()

	if data.ValidateItem(v, item); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Items.Insert(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/items/%d", item.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"item": item}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	item, err := app.models.Items.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
