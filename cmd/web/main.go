package main

import (
	"flag"

	"github.com/hunterwilkins2/go-template/internal/config"
	"github.com/hunterwilkins2/go-template/internal/server"
)

var version = "v1.0.0"

func main() {
	cfg := config.Config{}

	flag.IntVar(&cfg.Port, "port", 4000, "Server port")
	flag.StringVar(&cfg.DSN, "dsn", "", "MySQL Database DSN")
	flag.StringVar(&cfg.Version, "version", version, "Package version")
	flag.StringVar(&cfg.Env, "env", "development", "Server environment (development|production|testing)")
	flag.BoolVar(&cfg.HotReload, "hot-reload", false, "Hot reloads the browers when a change is made")
	flag.Parse()

	server := server.New(&cfg)
	server.Start()
}
