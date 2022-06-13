package server

import (
	"log"
	"net/http"
	"phrasetagg/url-shortener/internal/app/handlers"
	"phrasetagg/url-shortener/internal/app/handlers/api"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type serverCfg struct {
	Addr string `env:"SERVER_ADDRESS,required"`
}

func StartServer() {
	urlStorage := storage.NewURLStorage()
	shortener := models.NewShortener(urlStorage)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{shortURL}", handlers.GetFullURL(shortener))
		r.Post("/", handlers.ShortenURL(shortener))

		// /api routes
		r.Route("/api/", func(r chi.Router) {
			r.Post("/shorten", api.ShortenURL(shortener))
		})
	})

	var serverCfg serverCfg
	err := env.Parse(&serverCfg)
	if err != nil {
		panic("config error")
	}

	log.Fatal(http.ListenAndServe(serverCfg.Addr, r))
}
