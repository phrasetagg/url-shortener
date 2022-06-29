package handlers

import (
	"net/http"
	"phrasetagg/url-shortener/internal/app/models"
)

func Ping(db *models.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer db.Close()

		conn, err := db.GetConn(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = conn.Ping(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
