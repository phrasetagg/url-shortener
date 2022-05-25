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

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res = shortener.GetFullURL(id)

		w.Header().Set("Location", res)
		//w.WriteHeader(http.StatusTemporaryRedirect)
	case "POST":
		b, _ := io.ReadAll(r.Body)

		URL := string(b)

		if URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res = shortener.Shorten(URL)

		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(res))
		if err != nil {
			return
		}
	}
}
