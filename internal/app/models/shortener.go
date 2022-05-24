package models

import (
	"crypto/sha1"
	"encoding/hex"
)

type Shortener struct {
	Urls map[string]string
}

var singleton *Shortener

func init() {
	singleton = &Shortener{
		Urls: map[string]string{},
	}
}

func GetInstanceShortener() *Shortener {
	return singleton
}

func (s *Shortener) GetFullUrl(shortUrl string) string {
	return s.Urls[shortUrl]
}

func (s *Shortener) Shorten(url string) string {
	h := sha1.New()
	h.Write([]byte(url))

	encodedUrl := hex.EncodeToString(h.Sum(nil))
	s.Urls[encodedUrl] = url

	return encodedUrl
}
