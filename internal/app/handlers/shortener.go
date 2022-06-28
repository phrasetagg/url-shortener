package handlers

import (
	"io"
	"net/http"
	"phrasetagg/url-shortener/internal/app/models"

	"github.com/go-chi/chi/v5"
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

func ShortenURL(shortener models.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rawUserID := r.Context().Value("userID")
		var userID uint64

		switch uidType := rawUserID.(type) {
		case uint64:
			userID = uidType
		}

		b, _ := io.ReadAll(r.Body)

		URL := string(b)

		if URL == "" {
			http.Error(w, "URL in body is required", http.StatusBadRequest)
			return
		}

		res := shortener.Shorten(userID, URL)

		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(res))
		if err != nil {
			return
		}
	}
}
