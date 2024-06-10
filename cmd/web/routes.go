package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"snippetbox.example.org/ui"
)

// The routes() method returns a servemux containing our application routes.
// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServerMux.
func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Create a handler function which wraps our notFound() helper. and then
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Allowed responses by setting
	// router.NotFoundAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Take the ui.Files embeded filesystem and convert it to a http.FS type so
	// that it satisfiles the http.FileSystem interface. We then pass that to the
	// http.FileServer() function to create the files server handler.
	fileServer := http.FileServer(http.FS(ui.Files))

	// Our static files are contained in the "static" folder of the ui.Files
	// embeded filesystem. So, form example, our CSS stylesheet is located at
	// "static/css/main.css". This means the we now longer need to strip the
	// prefix from the request URL -- any requests that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// Add a new GET /ping route.
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to ir later.
	// Alice provides a convenient way to chain your HTTP middleware functions
	// and the app handler.
	// Use the nosurf middleware on all our 'dynamic' routes.
	// Add the authenticate() middleware to the chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Update these routes to use the new dynamic middleware chain followed by
	// the appropiate handler function. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switch to registering to route using the router.Handler() method.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	// Add the five new routes, all of which use our 'dynamic' middleware chain.
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	// Because the 'protected' middleware chain appends to the 'dynamic' chain
	// the noSurf middleware will also be used on the routes below too.
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Wrap the existing chain with the logRequest middleware
	// Wrap the existing chain with the recoverPanic middleware.
	// Wrap the existing chain with the chain your HTTP middleware functions
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
