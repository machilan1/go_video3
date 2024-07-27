package main

import (
	"net/http"

	"github.com/gotailwindcss/tailwind/twembed"
	"github.com/gotailwindcss/tailwind/twhandler"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./statics/"))

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	protected := dynamic.Append(app.requireAuthentication)
	updater := protected.Append(app.requireUpdater, app.requireAuthor)
	// TODO : Reserved for admin UI
	// admin := protected.Append(app.requireAdmin)

	mux.Handle("GET /statics/", http.StripPrefix("/statics", fileServer))
	mux.Handle("GET /css/", twhandler.New(http.Dir("statics/ui"), "/css", twembed.New()))

	// Publics
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.viewLogin))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.viewSignUp))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.login))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.signUp))
	//Protected
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))
	mux.Handle("POST /files/upload", protected.ThenFunc(app.upload))
	mux.Handle("GET /{$}", protected.ThenFunc(app.viewHome))
	mux.Handle("POST /{$}", protected.ThenFunc(app.viewHome))
	mux.Handle("GET /courses/{ID}", protected.ThenFunc(app.viewCourse))
	mux.Handle("GET /courses/{courseID}/chapters/{chapterID}", protected.ThenFunc(app.viewChapter))

	mux.Handle("GET /users/{userID}/courses", updater.ThenFunc(app.viewUserCourse))
	// TODO : Create course undeveloped.
	mux.Handle("GET /users/{userID}/courses/create", updater.ThenFunc(app.viewCreateUserCourse))
	mux.Handle("POST /users/{userID}/courses", updater.ThenFunc(app.createCourse))
	mux.Handle("GET /users/{userID}/courses/{courseID}", updater.ThenFunc(app.viewUserCourseDetail))
	mux.Handle("GET /users/{userID}/courses/{courseID}/edit", updater.ThenFunc(app.viewEditUserCourse))
	mux.Handle("PATCH /users/{userID}/courses/{courseID}", updater.ThenFunc(app.editCourse))
	mux.Handle("GET /users/{userID}/courses/{courseID}/chapters/{chapterID}/edit", updater.ThenFunc(app.viewEditCourseChapter))
	mux.Handle("PATCH /users/{userID}/courses/{courseID}/chapters/{chapterID}", updater.ThenFunc(app.editChapter))
	// TODO : delete course undeveloped.
	mux.Handle("DELETE /users/{userID}/courses/{courseID}", updater.ThenFunc(app.deleteCourse))

	// mux.HandleFunc("DELETE /user/{userId}", app.userDeactivate)

	standardMiddlewares := alice.New(app.changeMethod)
	return standardMiddlewares.Then(mux)
}
