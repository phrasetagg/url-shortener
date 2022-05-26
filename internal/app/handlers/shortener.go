package handlers

import (
	"io"
	"net/http"
	"phrasetagg/url-shortener/internal/app/models"
)

func Index(shortener models.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Path[1:]

			if id == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fullURL, err := shortener.GetFullURL(id)

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				_, err = w.Write([]byte(err.Error()))
				return
			}

			w.Header().Set("Location", fullURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		case http.MethodPost:
			b, _ := io.ReadAll(r.Body)

			URL := string(b)

			if URL == "" {
				w.WriteHeader(http.StatusBadRequest)
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
}
