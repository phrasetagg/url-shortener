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
	userID uint64
	url    string
}

type UserURLs struct {
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
}

func NewInMemoryURLStorage() *InMemoryURLStorage {
	return &InMemoryURLStorage{
		urls:  make(map[string]URLItem),
		mutex: new(sync.RWMutex),
	}
}

func (s *InMemoryURLStorage) GetItem(itemID string) (string, error) {
	s.mutex.RLock()
	urlItem, ok := s.urls[itemID]
	s.mutex.RUnlock()

	if !ok {
		return "", errors.New("not found")
	}

	return urlItem.url, nil
}

func (s *InMemoryURLStorage) AddItem(itemID string, value string, userID uint64) {
	s.mutex.Lock()
	s.urls[itemID] = URLItem{url: value, userID: userID}
	s.mutex.Unlock()
}

func (s InMemoryURLStorage) GetLastElementID() string {
	var shortURL string

	s.mutex.RLock()
	for key := range s.urls {
		shortURL = key
	}
	s.mutex.RUnlock()

	return shortURL
}

func (s InMemoryURLStorage) GetItemsByUserID(userID uint64) []UserURLs {
	var userURLs []UserURLs

	for key, value := range s.urls {
		if value.userID == userID {
			userURLs = append(userURLs, UserURLs{ShortURL: key, URL: value.url})
		}
	}

	return userURLs
}
