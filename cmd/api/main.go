// File: cmd/api/main.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Pedro-J-Kukul/qod/internal/data"

	_ "github.com/lib/pq"
)

// AppVersion
const AppVersion = "2.5.2"

// server configuration structure
type serverConfig struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	cors struct {
		trustedOrigins []string
	}
}

// application dependencies
type appDependencies struct {
	config serverConfig
	logger *slog.Logger
	model  data.QuoteModel
}

// Function: main
// Description: Entry point for the application
func main() {
	cfg := loadConfig()
	logger := setupLogger(cfg.env)

	db, err := openDB(cfg)
	if err != nil {
		logger.Error("Error opening database: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	app := appDependencies{
		config: cfg,
		logger: logger,
		model:  data.QuoteModel{DB: db},
	}

	err = app.serve()
	if err != nil {
		logger.Error("Error starting server: " + err.Error())
		os.Exit(1)
	}
}

// Function: serverConfig
// Description: Loads the server configuration from environment variables
func loadConfig() serverConfig {
	var cfg serverConfig
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "postgreSQL DSN")

	// custom command line for cors
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(s string) error {
		cfg.cors.trustedOrigins = strings.Fields(s)
		return nil
	})
	flag.Parse()
	return cfg
}

// Function: setupLogger
// Description: Sets up a logger for the application
func setupLogger(env string) *slog.Logger {
	// Create a new logger instance
	var logger *slog.Logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger = logger.With("environment", env)
	return logger
}

// Function: openDB
// Description: Opens a database connection
func openDB(cfg serverConfig) (*sql.DB, error) {
	// Create a connection pool with "sql.OPEN"
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a timeout to ensure DB operations don't hang
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to ensure a connection is established
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Return the database connection
	return db, nil
}
