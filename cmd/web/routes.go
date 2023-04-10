package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {

	//alice to manage middlewares
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	//use pat for routing
	mux := pat.New()
	// mux.Get("/thought/:id", http.HandlerFunc(app.showThought))
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/thought/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createThoughtForm))
	mux.Post("/thought/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createThought))
	mux.Get("/thought/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showThought))

	//Auths
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	//test route
	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
