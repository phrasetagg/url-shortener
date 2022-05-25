package main

import (
	"net/http"
	"phrasetagg/url-shortener/internal/app/controllers"
)

func main() {
	http.HandleFunc("/", controllers.Index)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
