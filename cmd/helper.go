package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) parseQueryInt(r *http.Request, key string) (int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return 0, nil
	} else {
		res, err := strconv.Atoi(param)
		if err != nil {
			return 0, err
		}
		return res, nil
	}
}
func (app *application) parseQuery(r *http.Request, key string) string {
	param := r.URL.Query().Get(key)
	return param
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {

	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data any) {

	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template  %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, fmt.Errorf("failed to execute templates"))
		return
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "currentUserID")
}
