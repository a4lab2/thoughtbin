package main

import (
	"path/filepath"
	"text/template"
	"time"

	"a4lab2.com/thoughtbin/pkg/forms"
	"a4lab2.com/thoughtbin/pkg/models"
)

type templateData struct {
	Thought         *models.Thought
	Thoughts        []*models.Thought
	CurrentYear     int
	Form            *forms.Form
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.gohtml"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.gohtml"))

		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.gohtml"))

		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
