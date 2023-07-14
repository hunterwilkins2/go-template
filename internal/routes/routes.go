package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/flow"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/hunterwilkins2/go-template/internal/config"
	"github.com/hunterwilkins2/go-template/internal/data"
	"github.com/hunterwilkins2/go-template/internal/templates"
	"go.uber.org/zap"
)

type Application struct {
	config         *config.Config
	logger         *zap.Logger
	sessionManager *scs.SessionManager
	queries        *data.Queries
	templates      *templates.HTMLTemplate
}

func New(config *config.Config, logger *zap.Logger, db *sql.DB) *Application {
	templates, err := templates.New("./ui/html/", config.HotReload)
	if err != nil {
		logger.Fatal("could not parse templates", zap.String("error", err.Error()))
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = true

	return &Application{
		config:         config,
		logger:         logger,
		sessionManager: sessionManager,
		queries:        data.New(db),
		templates:      templates,
	}
}

func (app *Application) Routes() http.Handler {
	mux := flow.New()

	mux.Use(app.recoverPanic, secureHeaders, app.logRequest)

	mux.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w, r)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/...", http.StripPrefix("/static", fileServer))

	if app.config.HotReload {
		mux.HandleFunc("/hot-reload", app.hotReload, http.MethodGet)
		mux.HandleFunc("/hot-reload/ready", app.testAlive, http.MethodGet)
	}

	mux.Group(func(m *flow.Mux) {
		mux.Use(app.sessionManager.LoadAndSave, noSurf)

		mux.HandleFunc("/", app.home, http.MethodGet)
	})

	return mux
}
