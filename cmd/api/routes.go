package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/items", app.createItemHandler)
	router.HandlerFunc(http.MethodGet, "/v1/item/:id", app.showItemHandler)
	return router
}
