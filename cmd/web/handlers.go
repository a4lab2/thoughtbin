package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	data, err := app.thoughts.Latest()
	if err != nil {
		app.notFound(w)
		// return nil, err
	}

	// data, err := app.thoughts.Get(id)

	dat := &templateData{
		Thoughts: data,
	}

	// fmt.Fprintf(w, "Show thought %+v\n", data)

	files := []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
	}
	err = ts.Execute(w, dat)
	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}
}

func (app *Application) showThought(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}
	data, err := app.thoughts.Get(id)

	dat := &templateData{
		Thought: data,
	}
	files := []string{
		"./ui/html/show.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}
	if err != nil {
		app.notFound(w)
		// return nil, err
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
	}
	err = ts.Execute(w, dat)
	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}
}

func (app *Application) createThought(w http.ResponseWriter, r *http.Request) {

	// if r.Method != http.MethodPost {

	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// f := models.Thought{
	// 	Title:   "x snail",
	// 	Content: "x snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n Kobayashi Issa",
	// 	Expires: 7,
	// }

	title1 := "1 snail"
	content1 := "1 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n Kobayashi Issa"
	expires1 := "7"

	id, err := app.thoughts.Insert(title1, content1, expires1)
	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}
	fmt.Fprintf(w, "Created thought ID: %d ", id)
}
