package storage

import (
	"errors"
	"sync"
)

type InMemoryURLStorage struct {
	urls  map[string]URLItem
	mutex *sync.RWMutex
}

type URLItem struct {
	userID uint32
	url    string
}

func NewInMemoryURLStorage() *InMemoryURLStorage {
	return &InMemoryURLStorage{
		urls:  make(map[string]URLItem),
		mutex: new(sync.RWMutex),
	}
}

func (s *InMemoryURLStorage) AddRecord(itemID string, value string, userID uint32) error {
	s.mutex.Lock()
	s.urls[itemID] = URLItem{url: value, userID: userID}
	s.mutex.Unlock()

	return nil
}

func (s *InMemoryURLStorage) GetOriginalURLByShortURI(itemID string) (OriginalURL, error) {
	var originalURL OriginalURL

	s.mutex.RLock()
	urlItem, ok := s.urls[itemID]
	originalURL.OriginalURL = urlItem.url
	s.mutex.RUnlock()

	if !ok {
		return originalURL, errors.New("not found")
	}

	return originalURL, nil
}

func (s InMemoryURLStorage) GetShortURIByOriginalURL(originalURL string) (string, error) {
	return originalURL, errors.New("in_memory_url_storage doesn't support this method")
}

func (s InMemoryURLStorage) GetRecordsByUserID(userID uint32) []UserURLs {
	var userURLs []UserURLs

	for key, value := range s.urls {
		if value.userID == userID {
			userURLs = append(userURLs, UserURLs{ShortURL: key, URL: value.url})
		}
	}

	return userURLs
}

func (s InMemoryURLStorage) DeleteUserRecordsByShortURLs(_ uint32, _ []string) error {
	return errors.New("in_memory_url_storage doesn't support this method")
}
