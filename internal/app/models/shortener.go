package models

import (
	"crypto/sha1"
	"encoding/hex"
	"phrasetagg/url-shortener/internal/app/storage"
)

type Shortener struct {
	storage storage.Storager
}

func NewShortener(storage storage.Storager) Shortener {
	return Shortener{
		storage: storage,
	}
}

func (s Shortener) GetFullURL(shortURL string) (string, error) {
	fullUrl, err := s.storage.GetItem(shortURL)

	return fullUrl, err
}

func (s Shortener) Shorten(URL string) (string, error) {
	h := sha1.New()
	h.Write([]byte(URL))

	encodedURL := hex.EncodeToString(h.Sum(nil))

	s.storage.AddItem(encodedURL, URL)

	return "http://localhost:8080/" + encodedURL, nil
}
