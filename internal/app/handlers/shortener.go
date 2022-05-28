package handlers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"phrasetagg/url-shortener/internal/app/models"
)

func GetFullURL(shortener models.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "shortURL")

		if id == "" {
			http.Error(w, "shortURL is required", http.StatusBadRequest)
			return
		}

		fullURL, err := shortener.GetFullURL(id)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}

		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func ShortenLink(shortener models.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)

		URL := string(b)

		if URL == "" {
			http.Error(w, "URL in body is required", http.StatusBadRequest)
			return
		}

		res, err := shortener.Shorten(URL)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			return
		}
	}
}
