package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.Handler(http.MethodGet, "/", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/signup", http.HandlerFunc(app.getSignup))
	router.Handler(http.MethodGet, "/login", http.HandlerFunc(app.getLogin))
	router.Handler(http.MethodGet, "/todo", http.HandlerFunc(app.getTodo))

	//api's for work with db
	router.Handler(http.MethodGet, "/todo/view/:id", http.HandlerFunc(app.viewTodo))
	router.Handler(http.MethodPost, "/todo/create/:id", http.HandlerFunc(app.createTodo))
	router.Handler(http.MethodPost, "/user/signup", http.HandlerFunc(app.postSignup))
	router.Handler(http.MethodPost, "/user/login", http.HandlerFunc(app.postLogin))

	// fmt.Println("Import \"fmt\"")

	return router
}
