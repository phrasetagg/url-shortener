package models

import (
	"crypto/sha1"
	"encoding/hex"
)

type Shortener struct {
	URLs map[string]string
}

var singleton *Shortener

func init() {
	singleton = &Shortener{
		URLs: map[string]string{},
	}
}

func GetInstanceShortener() *Shortener {
	return singleton
}

func (s *Shortener) GetFullURL(shortURL string) string {
	return s.URLs[shortURL]
}

func (s *Shortener) Shorten(URL string) string {
	h := sha1.New()
	h.Write([]byte(URL))

	encodedURL := hex.EncodeToString(h.Sum(nil))
	s.URLs[encodedURL] = URL

	return "https://localhost:8080/" + encodedURL
}
