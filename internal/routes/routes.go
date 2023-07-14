package routes

import (
	"database/sql"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/hunterwilkins2/go-template/internal/config"
	"github.com/hunterwilkins2/go-template/internal/data"
	"github.com/hunterwilkins2/go-template/internal/templates"
	"go.uber.org/zap"
)

type Application struct {
	config    *config.Config
	logger    *zap.Logger
	queries   *data.Queries
	templates *templates.HTMLTemplate
}

func New(config *config.Config, logger *zap.Logger, db *sql.DB) *Application {
	templates, err := templates.New("./ui/html/", config.HotReload)
	if err != nil {
		logger.Fatal("could not parse templates", zap.String("error", err.Error()))
	}

	return &Application{
		config:    config,
		logger:    logger,
		queries:   data.New(db),
		templates: templates,
	}
}

func (app *Application) Routes() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		app.notFound(w)
	})

	mux.HandleFunc("/", app.home, http.MethodGet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/...", http.StripPrefix("/static", fileServer))

	if app.config.HotReload {
		mux.HandleFunc("/hot-reload", app.hotReload, http.MethodGet)
		mux.HandleFunc("/hot-reload/ready", app.testAlive, http.MethodGet)
	}

	return mux
}
