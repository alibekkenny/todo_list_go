package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo_list/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type todoCreateForm struct {
	Title               string    `form:"title"`
	Description         string    `form:"description"`
	Expires             time.Time `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Write([]byte("Not found!"))
	}
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) getTodo(w http.ResponseWriter, r *http.Request) {
	userId := 1
	todos, err := app.todos.GetByUserId(userId)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	data := app.newTemplateData(r)
	data.Form = todoCreateForm{}
	data.Todos = todos
	app.render(w, http.StatusOK, "todo.html", data)
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
	var form todoCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.errorLog.Fatal(err)
		return
	}
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "todo.html", data)
		return
	}
	userId := 1
	err = app.todos.Insert(userId, form.Title, form.Description, form.Expires)
	if err != nil {
		app.errorLog.Fatal(err)
		return
	}

	http.Redirect(w, r, "/todos/", http.StatusSeeOther)
	// w.Write([]byte("create todo"))
}

func (app *application) getSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "register.html", data)
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
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post login"))
}
