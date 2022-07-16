package models

import (
	"errors"
	"math/rand"
	"phrasetagg/url-shortener/internal/app/storage"
)

type Shortener struct {
	storage storage.IURLStorager
	baseURL string
}

func NewShortener(storage storage.IURLStorager, baseURL string) Shortener {
	return Shortener{
		storage: storage,
		baseURL: baseURL,
	}
}

func (s Shortener) GetFullURL(shortURL string) (string, error) {
	fullURL, err := s.storage.GetOriginalURLByShortURI(shortURL)

	return fullURL, err
}

func (s Shortener) Shorten(userID uint32, URL string) (string, error) {
	var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	shortURIrunes := make([]rune, 10)
	for i := range shortURIrunes {
		shortURIrunes[i] = letters[rand.Intn(len(letters))]
	}
	shortURI := string(shortURIrunes)

	err := s.storage.AddRecord(shortURI, URL, userID)

	var iae *storage.ItemAlreadyExistsError

	if errors.As(err, &iae) {
		var getShortURIErr error
		shortURI, getShortURIErr = s.storage.GetShortURIByOriginalURL(URL)
		if getShortURIErr != nil {
			return "", getShortURIErr
		}
	}

	if err != nil && !errors.As(err, &iae) {
		return "", err
	}

	return s.baseURL + shortURI, err
}

func (s Shortener) GetUserURLs(userID uint32) []storage.UserURLs {
	var preparedUserURLs []storage.UserURLs

	userURLs := s.storage.GetRecordsByUserID(userID)

	for _, value := range userURLs {
		preparedUserURLs = append(preparedUserURLs, storage.UserURLs{ShortURL: s.baseURL + value.ShortURL, URL: value.URL})
	}

	return preparedUserURLs
}
