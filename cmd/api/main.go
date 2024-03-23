package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"chatapp.kyriakidis.net/internal/data"
	"github.com/redis/go-redis/v9"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models data.Models
	logger *slog.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server's Port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment (production | development)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "Database dsn")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		config: cfg,
		models: data.NewModels(db),
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}

}

func connectDB(cfg config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	return client, nil
}
