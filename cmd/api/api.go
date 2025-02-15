package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xurl/internal/store"
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
	store  store.Store
}

type config struct {
	addr  string
	dbCfg dbConfig
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/create", app.createShortenedUrlHandler)
	r.Get("/api/{slug}", app.getOriginalUrlHandler)

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Minute,
	}

	log.Println("Starting server on", app.config.addr)

	err := srv.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}
