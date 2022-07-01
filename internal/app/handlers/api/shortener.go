package api

import (
	"encoding/json"
	"io"
	"net/http"
	"phrasetagg/url-shortener/internal/app/middlewares"
	"phrasetagg/url-shortener/internal/app/models"
)

func GetUserURLs(shortener models.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")

		rawUserID := r.Context().Value(middlewares.UserID)
		var userID uint32

		switch uidType := rawUserID.(type) {
		case uint32:
			userID = uidType
		}

		userURLs := shortener.GetUserURLs(userID)
		responseBytes, _ := json.Marshal(userURLs)

		if len(userURLs) == 0 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		_, err := w.Write(responseBytes)
		if err != nil {
			return
		}
	}
}

func ShortenURL(shortener models.Shortener) http.HandlerFunc {
	type request struct {
		URL string `json:"url"`
	}

	type response struct {
		Result string `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rawUserID := r.Context().Value(middlewares.UserID)
		var userID uint32

		switch uidType := rawUserID.(type) {
		case uint32:
			userID = uidType
		}

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

		shortURL := shortener.Shorten(userID, request.URL)

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

func ShortenURLBatch(shortener models.Shortener) http.HandlerFunc {
	type dataToShorten struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	type resultData struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var request []dataToShorten
		var result []resultData

		rawUserID := r.Context().Value(middlewares.UserID)
		var userID uint32

		switch uidType := rawUserID.(type) {
		case uint32:
			userID = uidType
		}

		b, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(b, &request)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(request) == 0 {
			http.Error(w, `{"error":"URLs in body are required"}`, http.StatusBadRequest)
			return
		}

		for _, data := range request {
			shortURL := shortener.Shorten(userID, data.OriginalURL)
			result = append(result, resultData{CorrelationID: data.CorrelationID, ShortURL: shortURL})
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		responseBytes, _ := json.Marshal(result)

		_, err = w.Write(responseBytes)
		if err != nil {
			return
		}
	}
}
