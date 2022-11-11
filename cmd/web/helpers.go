package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	//2 is depth of trace tree
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// // Return true if the current request is from an authenticated user, otherwise
// // return false.
// func (app *application) isAuthenticated(r *http.Request) bool {
// 	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
// }

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		// app.serverError(w, err)
		app.errorLog.Fatal(err)
		return
	}
	// Initialize a new buffer.
	buf := new(bytes.Buffer)
	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		// app.serverError(w, err)
		app.errorLog.Fatal(err)
		return
	}
	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)
	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		// Add the flash message to the template data, if one exists.
		// Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}
