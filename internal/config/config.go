package config

// Config Server config
type Config struct {
	Port      int    // Server port
	DSN       string // Database DSN
	Version   string // Service Version
	Env       string // Service Environment (development|production|testing)
	HotReload bool   // Refreshes browser page when change is made
}
