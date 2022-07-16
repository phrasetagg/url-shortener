package handlers

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"
	"strings"
	"testing"
)

func TestSuccessShortenWithInMemoryURLStorage(t *testing.T) {
	type args struct {
		URL  string
		body string
	}

	type want struct {
		code                   string
		redirectLocationHeader string
		redirectCode           string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Status 201",
			args: args{
				URL:  "/",
				body: "https://reddit.com",
			},
			want: want{
				code:                   "201 Created",
				redirectCode:           "307 Temporary Redirect",
				redirectLocationHeader: "https://reddit.com",
			},
		},
		{
			name: "Status 201 N2",
			args: args{
				URL:  "/",
				body: "https://vk.com",
			},
			want: want{
				code:                   "201 Created",
				redirectCode:           "307 Temporary Redirect",
				redirectLocationHeader: "https://vk.com",
			},
		},
	}

	urlStorage := storage.NewInMemoryURLStorage()
	shortener := models.NewShortener(urlStorage, "http://localhost:8080/")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.args.body)

			request := httptest.NewRequest(http.MethodPost, tt.args.URL, buffer)
			w := httptest.NewRecorder()
			h := ShortenURL(shortener)
			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.Status, "Unexpected empty body")
			assert.NotEmpty(t, w.Body)
			shortURL := w.Body.String()

			// Проверяем работу редиректа
			request = httptest.NewRequest(http.MethodGet, shortURL, buffer)
			w = httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shortURL", strings.TrimPrefix(shortURL, "http://localhost:8080/"))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			h = GetFullURL(shortener)
			h.ServeHTTP(w, request)

			res = w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.redirectLocationHeader, res.Header.Get("Location"), "Unexpected Location header value")
			assert.Equal(t, tt.want.redirectCode, res.Status)

		})
	}
}

func TestFailShortenWithInMemoryURLStorage(t *testing.T) {
	type args struct {
		URL           string
		method        string
		shortURL      string
		body          string
		checkRedirect bool
	}

	type want struct {
		code                   string
		body                   string
		redirectLocationHeader string
		redirectCode           string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Status 400 GET Empty short url",
			args: args{
				URL:    "/",
				method: http.MethodGet,
			},
			want: want{
				code: "400 Bad Request",
				body: "shortURL is required\n",
			},
		},
		{
			name: "Status 400 GET Undefined short url",
			args: args{
				URL:      "/{shortURL}",
				method:   http.MethodGet,
				shortURL: "some_undefined_short_url",
			},
			want: want{
				code: "404 Not Found",
				body: "not found",
			},
		},
		{
			name: "Status 400 POST Empty body",
			args: args{
				URL:    "/",
				method: http.MethodPost,
			},
			want: want{
				code: "400 Bad Request",
				body: "URL in body is required\n",
			},
		},
	}

	urlStorage := storage.NewInMemoryURLStorage()
	shortener := models.NewShortener(urlStorage, "http://localhost:8080/")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.args.body)

			request := httptest.NewRequest(tt.args.method, tt.args.URL, buffer)
			w := httptest.NewRecorder()

			if tt.args.shortURL != "" {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("shortURL", tt.args.shortURL)

				request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			}

			h := GetFullURL(shortener)

			if tt.args.method == http.MethodPost {
				h = ShortenURL(shortener)
			}

			h.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.Status)
			assert.NotEmpty(t, w.Body)
			assert.Equal(t, tt.want.body, w.Body.String(), "Unexpected body")
		})
	}
}
