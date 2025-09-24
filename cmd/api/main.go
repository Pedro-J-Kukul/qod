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
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

// application dependencies
type appDependencies struct {
	config     serverConfig
	logger     *slog.Logger
	quoteModel data.QuoteModel
	userModel  data.UserModel
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
		config:     cfg,
		logger:     logger,
		quoteModel: data.QuoteModel{DB: db},
		userModel:  data.UserModel{DB: db},
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
	var cfg serverConfig                                                                            // Initialize a new serverConfig struct
	flag.IntVar(&cfg.port, "port", 4000, "API server port")                                         //  for port
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")   // for environment
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "postgreSQL DSN")                                     // for database DSN
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second") // for rate limiter rps
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 5, "Rate limiter maximum burst")               // for rate limiter burst
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")              // for enabling rate limiter

	// for trusted CORS origins
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(s string) error {
		cfg.cors.trustedOrigins = strings.Fields(s)
		return nil
	})

	flag.Parse() // Parse the command-line flags
	return cfg   // Return the populated serverConfig struct
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
