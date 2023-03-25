package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"a4lab2.com/thoughtbin/pkg/forms"
	"a4lab2.com/thoughtbin/pkg/models"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	data, err := app.thoughts.Latest()
	if err != nil {
		app.notFound(w)
		return
	}

	dat := &templateData{
		Thoughts: data,
	}

	app.render(w, r, "home.page.gohtml", dat)

}

func (app *Application) showThought(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.notFound(w)
		return
	}
	data, err := app.thoughts.Get(id)
	if err != nil {
		app.notFound(w)

	}

	dat := &templateData{
		Thought: data,
	}

	app.render(w, r, "show.page.gohtml", dat)
}

func (app *Application) createThought(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLenght("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.gohtml", &templateData{
			Form: form,
		})
	}

	id, err := app.thoughts.Insert(form.Get("title"), form.Get("content"), form.Get("content"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	//set session and redirect to the created thought
	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/thought/%d", id), http.StatusSeeOther)

}

func (app *Application) createThoughtForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "create.page.gohtml", &templateData{
		Form: forms.New(nil),
	})

}

//auths

func (app *Application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.gohtml", &templateData{
		Form: forms.New(nil),
	})
}
func (app *Application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLenght("name", 255)
	form.MaxLenght("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.gohtml", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.gohtml", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	fmt.Fprintln(w, "Create a new user...")
}
func (app *Application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", &templateData{
		Form: forms.New(nil),
	})
}
func (app *Application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or password is incorrect")
			app.render(w, r, "login.page.gohtml", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	http.Redirect(w, r, "/thought/create", http.StatusSeeOther)

}
func (app *Application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
