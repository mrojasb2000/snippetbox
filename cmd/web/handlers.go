package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Change the signature of the home handler so it is defined as a method againt
// *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	// Use the template.ParseFiles() function to read the template file intro a
	// template set. If there's an error, we log the detailed error message and use
	// the http.error() function to send a generic 500 Internal Server Error
	// response to the user.

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file
	// paths as a variadic parameter?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Because the home handler function is now a method againt application
		// it can access its fields, including the error logger. We'll write the log
		// message to this instead of the standard logger.
		app.errorLog.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// We then use the Execute() method on the template set to write the
	// template content as the response body.
	// represents any dynamic data that we to pass in, which for now we'll
	// leave as nil.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Change the signature of the snippetView handler so it is defined as a method
// againt *application
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < -1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Change the signature of the snippetCreate handler so it is defined as a method
// againt *application
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
