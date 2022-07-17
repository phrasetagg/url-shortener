package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"phrasetagg/url-shortener/internal/app/middlewares"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"
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

func DeleteUserURLs(shortener models.Shortener) http.HandlerFunc {
	var request []string

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

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
			http.Error(w, `{"error":"Invalid body"}`, http.StatusBadRequest)
			return
		}

		if len(request) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, `{"error":"Body must be not empty array"}`, http.StatusBadRequest)
			return
		}

		shortener.DeleteURLs(userID, request)
		w.WriteHeader(http.StatusAccepted)
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
		w.Header().Set("content-type", "application/json")

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
			http.Error(w, `{"error":"URL in body is required"}`, http.StatusBadRequest)
			return
		}

		if request.URL == "" {
			http.Error(w, `{"error":"URL in body is required"}`, http.StatusBadRequest)
			return
		}

		shortURL, err := shortener.Shorten(userID, request.URL)
		response := response{Result: shortURL}

		var iae *storage.ItemAlreadyExistsError

		if errors.As(err, &iae) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		if err != nil && !errors.As(err, &iae) {
			http.Error(w, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
			return
		}

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
		w.Header().Set("content-type", "application/json")

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
			http.Error(w, `{"error":"URLs in body are required"}`, http.StatusBadRequest)
			return
		}

		if len(request) == 0 {
			http.Error(w, `{"error":"URLs in body are required"}`, http.StatusBadRequest)
			return
		}

		var shortenError error

		for _, data := range request {
			var shortURL string
			shortURL, shortenError = shortener.Shorten(userID, data.OriginalURL)

			result = append(result, resultData{CorrelationID: data.CorrelationID, ShortURL: shortURL})
		}

		var iae *storage.ItemAlreadyExistsError
		if errors.As(shortenError, &iae) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		if shortenError != nil && !errors.As(shortenError, &iae) {
			http.Error(w, `{"error":"Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		responseBytes, _ := json.Marshal(result)

		_, err = w.Write(responseBytes)
		if err != nil {
			return
		}
	}
}
