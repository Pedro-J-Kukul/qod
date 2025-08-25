// File: cmd/api/main.go

package main

import (
	"flag"
	"log/slog"
	"os"
)

// application configuration with the help of an environment file
type configuration struct {
	port    int
	env     string
	version string
}

// dependency injec6tion for the application. Uses one instance of each service
type application struct {
	config configuration
	logger *slog.Logger
}

func main() {
	// Load application configuration
	cfg := loadConfig()

	// Setup application logger
	logger := setupLogger(cfg.env)

	app := application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func loadConfig() configuration {
	var cfg configuration

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&cfg.version, "version", "1.0.0", "Application version")
	flag.Parse()

	return cfg
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger = logger.With("environment", env)

	return logger
}
