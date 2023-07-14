package routes

import "net/http"

func (app *Application) serverError(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func (app *Application) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	err := app.templates.Render(w, "not-found.page.html", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}
