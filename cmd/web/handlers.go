package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func (app *application) viewTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("view todo"))
}

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create todo"))
}

func (app *application) getSignup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get registration template"))
}
func (app *application) postSignup(w http.ResponseWriter, r *http.Request) {

	userId, err := app.users.Insert("danik", "myEmail", "myPassword")
	if err != nil {
		app.errorLog.Fatal(err)
	}

	fmt.Println(userId)

}

func (app *application) getLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get login template"))
}

func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post login"))
}
