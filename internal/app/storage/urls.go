package storage

import (
	"errors"
	"sync"
)

type URLStorage struct {
	urls  map[string]string
	mutex *sync.RWMutex
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		urls:  make(map[string]string),
		mutex: new(sync.RWMutex),
	}
}

func (u URLStorage) GetItem(itemID string) (string, error) {
	u.mutex.RLock()
	item, ok := u.urls[itemID]
	u.mutex.RUnlock()

	if !ok {
		return "", errors.New("not found")
	}

	return item, nil
}

func (u *URLStorage) AddItem(itemID string, value string) {
	u.mutex.Lock()
	u.urls[itemID] = value
	u.mutex.Unlock()
}

func (u URLStorage) GetItems() map[string]string {
	return u.urls
}
