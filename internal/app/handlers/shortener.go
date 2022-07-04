package handlers

import (
	"errors"
	"io"
	"net/http"
	"phrasetagg/url-shortener/internal/app/middlewares"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"

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

		rawUserID := r.Context().Value(middlewares.UserID)
		var userID uint32

		switch uidType := rawUserID.(type) {
		case uint32:
			userID = uidType
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		URL := string(b)

		if URL == "" {
			http.Error(w, "URL in body is required", http.StatusBadRequest)
			return
		}

		res, err := shortener.Shorten(userID, URL)

		var iae *storage.ItemAlreadyExistsError

		if errors.As(err, &iae) {
			w.WriteHeader(http.StatusConflict)
		}

		if err != nil {
			http.Error(w, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			http.Error(w, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
			return
		}
	}
}
