package routes

import "net/http"

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	err := app.templates.Render(w, "home.page.html", nil)
	if err != nil {
		app.serverError(w, err)
	}
}
