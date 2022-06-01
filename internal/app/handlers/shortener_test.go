package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	type args struct {
		URL      string
		method   string
		shortURL string
		body     string
	}

	type want struct {
		code           string
		locationHeader string
		body           string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Status 201",
			args: args{
				URL:    "/",
				method: http.MethodPost,
				body:   "https://reddit.com",
			},
			want: want{
				code: "201 Created",
				body: "http://localhost:8080/a",
			},
		},
		{
			name: "Status 307",
			args: args{
				URL:      "/{shortURL}",
				method:   http.MethodGet,
				shortURL: "a",
			},
			want: want{
				code:           "307 Temporary Redirect",
				locationHeader: "https://reddit.com",
			},
		},
		{
			name: "Status 201 N2",
			args: args{
				URL:    "/",
				method: http.MethodPost,
				body:   "https://vk.com",
			},
			want: want{
				code: "201 Created",
				body: "http://localhost:8080/b",
			},
		},
		{
			name: "Status 307 N2",
			args: args{
				URL:      "/{shortURL}",
				method:   http.MethodGet,
				shortURL: "b",
			},
			want: want{
				code:           "307 Temporary Redirect",
				locationHeader: "https://vk.com",
			},
		},
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
		{
			name: "Status 307 Check that first link is still working",
			args: args{
				URL:      "/{shortURL}",
				method:   http.MethodGet,
				shortURL: "a",
			},
			want: want{
				code:           "307 Temporary Redirect",
				locationHeader: "https://reddit.com",
			},
		},
	}

	urlStorage := storage.NewURLStorage()
	shortener := models.NewShortener(urlStorage)

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
			assert.Equal(t, tt.want.body, w.Body.String(), "Unexpected body")

			if tt.want.locationHeader != "" {
				assert.Equal(t, tt.want.locationHeader, res.Header.Get("Location"), "Unexpected Location header value")
			}
		})
	}
}
