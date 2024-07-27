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
	admin := protected.Append(app.requireAdmin)

	mux.Handle("GET /statics/", http.StripPrefix("/statics", fileServer))
	mux.Handle("GET /css/", twhandler.New(http.Dir("statics/ui"), "/css", twembed.New()))

	// Publics
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.viewLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.login))

	//Protected
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))
	mux.Handle("POST /files/upload", protected.ThenFunc(app.upload))
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.viewHome))
	// TODO : Develope this later
	mux.Handle("GET /courses/search", dynamic.ThenFunc(app.searchCourse))
	mux.Handle("GET /courses/{ID}", protected.ThenFunc(app.viewCourse))
	mux.Handle("GET /courses/{courseID}/chapters/{chapterID}", protected.ThenFunc(app.viewChapter))

	mux.Handle("GET /users/{userID}/courses", updater.ThenFunc(app.viewUserCourse))
	mux.Handle("GET /users/{userID}/courses/create", updater.ThenFunc(app.viewCreateUserCourse))
	mux.Handle("POST /users/{userID}/courses", updater.ThenFunc(app.createCourse))
	mux.Handle("GET /users/{userID}/courses/{courseID}", updater.ThenFunc(app.viewUserCourseDetail))
	mux.Handle("PATCH /users/{userID}/courses/{courseID}", updater.ThenFunc(app.editCourse))
	mux.Handle("DELETE /users/{userID}/courses/{courseID}", updater.ThenFunc(app.deleteCourse))
	mux.Handle("GET /users/{userID}/courses/{courseID}/edit", updater.ThenFunc(app.viewEditUserCourse))
	mux.Handle("GET /users/{userID}/courses/{courseID}/chapters/create", updater.ThenFunc(app.viewCreateChapter))
	mux.Handle("POST /users/{userID}/courses/{courseID}/chapters", updater.ThenFunc(app.createChapter))
	mux.Handle("PATCH /users/{userID}/courses/{courseID}/chapters/{chapterID}", updater.ThenFunc(app.editChapter))
	mux.Handle("DELETE /users/{userID}/courses/{courseID}/chapters/{chapterID}", updater.ThenFunc(app.deleteChapter))
	mux.Handle("GET /users/{userID}/courses/{courseID}/chapters/{chapterID}/edit", updater.ThenFunc(app.viewEditCourseChapter))

	mux.Handle("GET /admin/users", admin.ThenFunc(app.viewAdminUsers))
	// mux.Handle("GET /admin/courses", admin.ThenFunc(app.viewAdminCourses))
	// mux.Handle("GET /user/signup", admin.ThenFunc(app.viewSignUp))
	// mux.Handle("POST /user/signup", admin.ThenFunc(app.signUp))
	// mux.Handle("PATCH /users/{userID}/change-role", admin.ThenFunc(app.changeUserRole))
	// mux.Handle("PATCH /users/{userID}/change-password", admin.ThenFunc(app.changeUserPassword))
	// mux.Handle("GET /tags", admin.ThenFunc(app.getTags))
	// mux.Handle("POST /tags", admin.ThenFunc(app.createTag))
	// mux.Handle("PATCH /tags/{tagID}", admin.ThenFunc(app.updateTag))

	// mux.HandleFunc("DELETE /user/{userId}", app.userDeactivate)

	standardMiddlewares := alice.New(app.changeMethod)
	return standardMiddlewares.Then(mux)
}
