package storage

import (
	"errors"
	"sync"
)

type InMemoryURLStorage struct {
	urls  map[string]string
	mutex *sync.RWMutex
}

func NewInMemoryURLStorage() *InMemoryURLStorage {
	return &InMemoryURLStorage{
		urls:  make(map[string]string),
		mutex: new(sync.RWMutex),
	}
}

func (s *InMemoryURLStorage) GetItem(itemID string) (string, error) {
	s.mutex.RLock()
	item, ok := s.urls[itemID]
	s.mutex.RUnlock()

	if !ok {
		return "", errors.New("not found")
	}

	return item, nil
}

func (s *InMemoryURLStorage) AddItem(itemID string, value string) {
	s.mutex.Lock()
	s.urls[itemID] = value
	s.mutex.Unlock()
}

func (s InMemoryURLStorage) GetLastElementID() string {
	var shortURL string

	for key := range s.urls {
		shortURL = key
	}

	return shortURL
}
