package main

import (
	"database/sql"
	"depopa/internal/data"
	"depopa/internal/jsonlog"
	"flag"
	"github.com/joho/godotenv"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	models data.Models
	config config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	godotenv.Load()

	flag.IntVar(&cfg.port, "port", 9696, "The port to run this app on")
	flag.StringVar(&cfg.env, "env", "development", "development|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DEPOPA_DSN_DB"), "PostgeSQL DSN")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		models: data.NewModels(db),
		config: cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
