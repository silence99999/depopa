package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/items", app.createItemHandler)
	router.HandlerFunc(http.MethodGet, "/v1/items", app.showItemsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/item/:id", app.showItemHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/item/:id", app.deleteItemHandler)
	router.HandlerFunc(http.MethodPut, "/v1/item/:id", app.updateItemHandler)
	return router
}
