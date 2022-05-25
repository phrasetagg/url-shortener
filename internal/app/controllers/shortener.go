package controllers

import (
	"io"
	"net/http"
	"phrasetagg/url-shortener/internal/app/models"
)

// Index обработчик запросов пути /.
func Index(w http.ResponseWriter, r *http.Request) {
	res := ""
	shortener := models.GetInstanceShortener()

	switch r.Method {
	case "GET":
		id := r.URL.Path[1:]
		res = shortener.GetFullUrl(id)
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "POST":
		b, _ := io.ReadAll(r.Body)
		res = shortener.Shorten(string(b))

		_, err := w.Write([]byte(res))
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
