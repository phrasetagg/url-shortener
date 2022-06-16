package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GzipRequestDecoder() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if r.ContentLength == 0 {
				next.ServeHTTP(w, r)
				return
			}

			if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gzReader.Close()

			r.Body = gzReader

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
