package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.Handler(http.MethodGet, "/", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/todo/view/:id", http.HandlerFunc(app.viewTask))
	router.Handler(http.MethodPost, "/todo/create", http.HandlerFunc(app.createTask))
	router.Handler(http.MethodGet, "/user/signup", http.HandlerFunc(app.getSignup))
	router.Handler(http.MethodPost, "/user/signup", http.HandlerFunc(app.postSignup))
	router.Handler(http.MethodGet, "/user/login", http.HandlerFunc(app.getLogin))
	router.Handler(http.MethodPost, "/user/login", http.HandlerFunc(app.postLogin))

	// fmt.Println("Import \"fmt\"")

	return router
}
