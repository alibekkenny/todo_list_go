package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	//Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(http.HandlerFunc(app.home)))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(http.HandlerFunc(app.getSignup)))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(http.HandlerFunc(app.getLogin)))
	router.Handler(http.MethodGet, "/todo", dynamic.ThenFunc(http.HandlerFunc(app.getTodo)))

	//api's for work with db
	router.Handler(http.MethodGet, "/todo/view/:id", dynamic.ThenFunc(http.HandlerFunc(app.viewTodo)))
	router.Handler(http.MethodPost, "/todo/create/:id", dynamic.ThenFunc(http.HandlerFunc(app.createTodo)))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(http.HandlerFunc(app.postSignup)))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(http.HandlerFunc(app.postLogin)))

	// router.Handler(http.MethodGet, "/", (http.HandlerFunc(app.home)))
	// router.Handler(http.MethodGet, "/signup", (http.HandlerFunc(app.getSignup)))
	// router.Handler(http.MethodGet, "/login", (http.HandlerFunc(app.getLogin)))
	// router.Handler(http.MethodGet, "/todo", (http.HandlerFunc(app.getTodo)))

	// //api's for work with db
	// router.Handler(http.MethodGet, "/todo/view/:id", (http.HandlerFunc(app.viewTodo)))
	// router.Handler(http.MethodPost, "/todo/create/:id", (http.HandlerFunc(app.createTodo)))
	// router.Handler(http.MethodPost, "/user/signup", (http.HandlerFunc(app.postSignup)))
	// router.Handler(http.MethodPost, "/user/login", (http.HandlerFunc(app.postLogin)))

	// fmt.Println("Import \"fmt\"")

	wrapper := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return wrapper.Then(router)
}
