package server

import (
	"log"
	"net/http"
	"phrasetagg/url-shortener/internal/app/config"
	"phrasetagg/url-shortener/internal/app/handlers"
	"phrasetagg/url-shortener/internal/app/handlers/api"
	"phrasetagg/url-shortener/internal/app/middlewares"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var cfg = config.PrepareCfg()

func StartServer() {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middlewares.GzipRequestDecoder())
	r.Use(middlewares.GzipResponseEncode())
	r.Use(middlewares.GenerateAuthToken())

	urlStorage := createURLStorage()
	shortener := models.NewShortener(urlStorage, cfg.BaseURL)

	r.Route("/", func(r chi.Router) {
		r.Get("/{shortURL}", handlers.GetFullURL(shortener))
		r.Post("/", handlers.ShortenURL(shortener))

		// /api routes
		r.Route("/api/", func(r chi.Router) {
			r.Get("/user/urls", api.GetUserURLs(shortener))
			r.Post("/shorten", api.ShortenURL(shortener))
		})
	})

	log.Fatal(http.ListenAndServe(cfg.ServerAddr, r))
}

func createURLStorage() storage.IStorager {
	var urlStorage storage.IStorager

	if cfg.FileStoragePath == "" {
		urlStorage = storage.NewInMemoryURLStorage()
	} else {
		urlStorage = storage.NewFileURLStorage(cfg.FileStoragePath)
	}

	return urlStorage
}
