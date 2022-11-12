package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo_list/internal/models"
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

type userRegisterForm struct {
	Nickname            string `form:"nickname"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) getSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userRegisterForm{}
	app.render(w, http.StatusOK, "register.html", data)
}

func (app *application) postSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("nickname")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	fmt.Println("*****", r.Body, "***", name, email, password)
	// err = app.decodePostForm(r, &form)
	// if err != nil {
	// 	app.errorLog.Fatal(err)
	// 	return
	// }

	form := userRegisterForm{
		Nickname: name,
		Email:    email,
		Password: password,
	}

	form.CheckField(validator.NotBlank(form.Nickname), "nickname", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "register.html", data)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	userExists, err := app.users.CheckUserExists(form.Email)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Something went wrong!")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "register.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	if userExists {
		form.AddFieldError("email", "Email is unavaliable!")
		data := app.newTemplateData(r)
		data.Form = form
		// http.Redirect(w, r, "/signup/", http.StatusSeeOther)
		// return
		app.render(w, http.StatusUnprocessableEntity, "register.html", data)
	}

	createdUser, err := app.users.Insert(form.Nickname, form.Email, form.Password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", createdUser.Id)

	http.Redirect(w, r, "/todo/", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) getLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	// Decode the form data into the userLoginForm struct.
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	// The SessionManager.RenewToken() method that we’re using in the code above
	// will change the ID of the current user’s session but retain any data
	// associated with the session. It’s good practice to do this before login
	// to mitigate the risk of a session fixation attack.
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/todo/", http.StatusSeeOther)
}
