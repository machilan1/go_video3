package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/machilan1/go_video/internal/store"
)

type templateData struct {
	Timestamp       time.Time
	Form            any
	Flash           string
	UserID          int
	CourseID        int
	IsAuthenticated bool
	CanUpdate       bool
	IsAdmin         bool
	CourseCards     []courseCardVM
	UserCards       []userCardVM
	Course          CourseVM
	Chapter         ChapterVM
	Options         OptionsVM
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		Timestamp:       time.Now(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CanUpdate:       app.sessionManager.GetBool(r.Context(), "canUpdate"),
		IsAdmin:         app.sessionManager.GetBool(r.Context(), "isAdmin"),
		UserID:          app.sessionManager.GetInt(r.Context(), "currentUserID"),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./statics/ui/pages/*.html")
	if err != nil {
		return nil, err
	}

	partials, err := filepath.Glob("./statics/ui/components/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"./statics/ui/base.html",
			page}

		slice := append(files, partials...)

		ts, err := template.ParseFiles(slice...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

type courseCardVM struct {
	store.Course
}

// TODO make sure you need this or not
type userCardVM struct {
	email     string
	name      string
	updatable bool
	isAdmin   bool
}

type CourseVM struct {
	store.Course
}

type ChapterVM struct {
	store.Chapter
}

type OptionsVM struct {
	Tags []store.CourseTag
}
