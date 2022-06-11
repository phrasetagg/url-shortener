package api

import (
	"encoding/json"
	"io"
	"net/http"
	"phrasetagg/url-shortener/internal/app/models"
)

func ShortenURL(shortener models.Shortener) http.HandlerFunc {
	type request struct {
		URL string `json:"url"`
	}

	type response struct {
		Result string `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)

		var request request

		err := json.Unmarshal(b, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if request.URL == "" {
			http.Error(w, `{"error":"URL in body is required"}`, http.StatusBadRequest)
			return
		}

		shortURL := shortener.Shorten(request.URL)

		response := response{Result: shortURL}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		responseBytes, _ := json.Marshal(response)

		_, err = w.Write(responseBytes)
		if err != nil {
			return
		}
	}
}
