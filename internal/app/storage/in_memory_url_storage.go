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

func (s *InMemoryURLStorage) GetOriginalURLByShortURI(itemID string) (string, error) {
	s.mutex.RLock()
	urlItem, ok := s.urls[itemID]
	s.mutex.RUnlock()

	if !ok {
		return "", errors.New("not found")
	}

	return urlItem.url, nil
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
