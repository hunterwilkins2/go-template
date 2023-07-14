package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hunterwilkins2/go-template/internal/config"
	"github.com/hunterwilkins2/go-template/internal/routes"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	db     *sql.DB
	server *http.Server
}

func New(cfg *config.Config) *Server {
	logger := getLogger(cfg)
	db, err := openDB(cfg.DSN)
	if err != nil {
		logger.Fatal("could not connect to the database", zap.String("error", err.Error()))
	}

	app := routes.New(cfg, logger, db)
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        app.Routes(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		logger: logger,
		db:     db,
		server: &server,
	}
}

func (s *Server) Start() {
	defer s.logger.Sync()
	defer s.db.Close()

	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		s.logger.Info("shutting down server", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		shutdownErr <- s.server.Shutdown(ctx)
	}()

	s.logger.Info("Starting server", zap.String("addr", s.server.Addr))
	err := s.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal("uncaught error occurred", zap.String("error", err.Error()))
	}

	if err = <-shutdownErr; err != nil {
		s.logger.Fatal("error shutting down server", zap.String("error", err.Error()))
	}

	s.logger.Info("stopped server")
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getLogger(cfg *config.Config) *zap.Logger {
	switch cfg.Env {
	case "production":
		return zap.Must(zap.NewProduction())
	case "testing":
		return zap.Must(zap.NewDevelopment())
	default:
		return zap.Must(zap.NewDevelopment())
	}
}
