package main

import (
	"net/http"
	"strconv"
)

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireUpdater(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.sessionManager.GetBool(r.Context(), "isAdmin") {
			next.ServeHTTP(w, r)
			return
		}

		if !app.sessionManager.GetBool(r.Context(), "canUpdate") {
			app.clientError(w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) changeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		method := r.PostForm.Get("_method")

		if method != "" {
			r.Method = method
		}
		next.ServeHTTP(w, r)
	})
}

// TODO : Use this in admin UIs
func (app *application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.sessionManager.GetBool(r.Context(), "isAdmin") {
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (app *application) requireAuthor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.sessionManager.GetBool(r.Context(), "isAdmin") {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := strconv.Atoi(r.PathValue("userID"))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		sessionUserID := app.sessionManager.GetInt(r.Context(), "currentUserID")

		if userID != sessionUserID {
			app.clientError(w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
