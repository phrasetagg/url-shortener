package storage

import "fmt"

type IURLStorager interface {
	AddRecord(shortURI string, originalURL string, userID uint32) error
	GetOriginalURLByShortURI(shortURI string) (OriginalURL, error)
	GetShortURIByOriginalURL(originalURL string) (string, error)
	GetRecordsByUserID(userID uint32) []UserURLs
	DeleteUserRecordsByShortURLs(userID uint32, shortURLs []string) error
}

type OriginalURL struct {
	OriginalURL string
	IsDeleted   bool
}

type UserURLs struct {
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
}

type ItemAlreadyExistsError struct {
	value string
}

func (iae *ItemAlreadyExistsError) Error() string {
	return fmt.Sprintf("record for %s already exists", iae.value)
}
