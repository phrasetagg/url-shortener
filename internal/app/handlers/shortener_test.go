package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"phrasetagg/url-shortener/internal/app/models"
	"phrasetagg/url-shortener/internal/app/storage"
	"testing"
)

func TestIndex(t *testing.T) {
	type args struct {
		URL    string
		method string
		body   string
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
				body: "http://localhost:8080/e7f4f2110990a57302e2639e33e465092613f209",
			},
		},
		{
			name: "Status 307",
			args: args{
				URL:    "/e7f4f2110990a57302e2639e33e465092613f209",
				method: http.MethodGet,
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
				body: "http://localhost:8080/5f12e5a8cc3d801aea41913df4fc427919aa0799",
			},
		},
		{
			name: "Status 307 N2",
			args: args{
				URL:    "/5f12e5a8cc3d801aea41913df4fc427919aa0799",
				method: http.MethodGet,
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
			},
		},
		{
			name: "Status 400 GET Undefined short url",
			args: args{
				URL:    "/some_undefined_short_url",
				method: http.MethodGet,
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
			},
		},
		{
			name: "Status 307 Check that first link is still working",
			args: args{
				URL:    "/e7f4f2110990a57302e2639e33e465092613f209",
				method: http.MethodGet,
			},
			want: want{
				code:           "307 Temporary Redirect",
				locationHeader: "https://reddit.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := models.NewShortener(storage.GetURLsInstance())

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.args.body)

			request := httptest.NewRequest(tt.args.method, tt.args.URL, buffer)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Index(shortener))
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
