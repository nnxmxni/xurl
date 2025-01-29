package main

import (
	"github.com/xurl/internal/db"
	"github.com/xurl/internal/env"
	"github.com/xurl/internal/store"
	"log"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		dbCfg: dbConfig{
			dsn:          env.GetString("DB_ADDR", "postgres://doadmin:pa55w0rd@localhost/chaos?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	database, err := db.New(cfg.dbCfg.dsn, cfg.dbCfg.maxOpenConns, cfg.dbCfg.maxIdleConns, cfg.dbCfg.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	log.Println("Database connection pool established")

	app := &application{
		config: cfg,
		store:  &store.URLStore{Db: database},
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
