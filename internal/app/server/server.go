package server

import (
	"log"
	"net/http"
	"os"
	"phrasetagg/url-shortener/internal/app/handlers"
	"phrasetagg/url-shortener/internal/app/handlers/api"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartServer() {
	var urlStorage storage.IStorager

	urlsFilePath := os.Getenv("FILE_STORAGE_PATH")

	if urlsFilePath == "" {
		urlStorage = storage.NewInMemoryURLStorage()
	} else {
		urlStorage = storage.NewFileURLStorage(urlsFilePath)
	}

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

	serveAddr := os.Getenv("SERVER_ADDRESS")

	if serveAddr == "" {
		serveAddr = "localhost:8080"
	}

	log.Fatal(http.ListenAndServe(serveAddr, r))
}
