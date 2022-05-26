package server

import (
	"log"
	"net/http"
	"phrasetagg/url-shortener/internal/app/handlers"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"
)

func StartServer() {
	shortener := models.NewShortener(storage.GetURLsInstance())

	http.HandleFunc("/", handlers.Index(shortener))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
