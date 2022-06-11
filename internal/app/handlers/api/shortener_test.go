package api

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"
	"strings"
	"testing"
)

func TestShortenURL(t *testing.T) {
	type args struct {
		URL      string
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
				URL:  "/api/shorten",
				body: `{"url":"https://reddit.com"}`,
			},
			want: want{
				code: "201 Created",
				body: `{"result":"http://localhost:8080/a"}`,
			},
		},
		{
			name: "Status 400 POST invalid json body",
			args: args{
				URL:  "/api/shorten",
				body: `{"}`,
			},
			want: want{
				code: "400 Bad Request",
			},
		},
		{
			name: "Status 400 POST !isset url",
			args: args{
				URL:  "/api/shorten",
				body: `{}`,
			},
			want: want{
				code: "400 Bad Request",
				body: `{"error":"URL in body is required"}`,
			},
		},
		{
			name: "Status 400 POST Empty url",
			args: args{
				URL:  "/api/shorten",
				body: `{"url":""}`,
			},
			want: want{
				code: "400 Bad Request",
				body: `{"error":"URL in body is required"}`,
			},
		},
	}

	urlStorage := storage.NewURLStorage()
	shortener := models.NewShortener(urlStorage)

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

			fmt.Println(w.Body.String())

			assert.Equal(t, tt.want.code, res.Status)
			assert.Equal(t, tt.want.body, strings.Trim(w.Body.String(), "\n"), "Unexpected body")

			if tt.want.locationHeader != "" {
				assert.Equal(t, tt.want.locationHeader, res.Header.Get("Location"), "Unexpected Location header value")
			}
		})
	}
}
