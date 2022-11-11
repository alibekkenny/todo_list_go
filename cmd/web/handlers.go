package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Write([]byte("Not found!"))
	}
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "index_2.tmpl", data)
}

func (app *application) viewTodo(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || userId < 1 {
		app.errorLog.Fatal(err)
		return
	}

	todos, err := app.todos.GetByUserId(userId)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	fmt.Fprintf(w, fmt.Sprintf("User id: %d\n", userId))
	for _, todo := range todos {
		fmt.Fprintf(w, "%+v\n", todo)
	}
}

func (app *application) createTodo(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || userId < 1 {
		app.errorLog.Fatal(err)
		return
	}

	createdtodo, err := app.todos.Insert(userId, "Exam preparation", "Prepare for an exam on the course DBMS", "study", time.Now().AddDate(0, 0, 10))
	fmt.Fprintf(w, "%+v\n", createdtodo)
	// w.Write([]byte("create todo"))
}

func (app *application) getSignup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get registration template"))
}
func (app *application) postSignup(w http.ResponseWriter, r *http.Request) {
	createdUser, err := app.users.Insert("aliba", "example_email", "123456")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	// w.Write([]byte(fmt.Sprintf("Id:%d, nickname:%s, email:%s", createdUser.Id, createdUser.Nickname, createdUser.Email)))
	fmt.Fprintf(w, "%+v\n", createdUser)
	// fmt.Println(userId)
}

func (app *application) getLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get login template"))
}

func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post login"))
}
