package main

import (
	"log"
	"net/http"
)

// Define a home handler function which writes a slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it does't, use
	// the http.notFound() function to send a 404 response to the client.
	// Importantly, we then return from handler. If we don't return the handler
	// would keep executing and also write the "Hello from SnippetBox" message
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from SnippetBox"))
}

// Add s snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	if r.Method != "POST" {
		// Use the Header().Set() method to add an 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second parameter is the header value.
		w.Header().Set("Allow", "POST")

		// If it's not, use the w.WriteHeader() method to send a 405 status
		// code and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code it not executed.
		//w.WriteHeader(405)
		//w.Write([]byte("Method Not Allowed"))

		// Use the http.Error() function to send 405 status code and "Method Not
		// Allowed" string as the response body.
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Register the two handler functions and corresponding URL patterns with
	// the servemux, in exactly the same way that we did before.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use the http.ListeAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() return an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
