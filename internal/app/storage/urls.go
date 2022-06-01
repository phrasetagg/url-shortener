package storage

import (
	"errors"
	"sync"
)

var (
	urls  = map[string]string{}
	mutex = &sync.RWMutex{}
)

type URLStorage struct{}

func NewURLStorage() *URLStorage {
	return &URLStorage{}
}

func (u URLStorage) GetItem(itemID string) (string, error) {
	mutex.RLock()
	item, ok := urls[itemID]
	mutex.RUnlock()

	if !ok {
		return "", errors.New("not found")
	}

	return item, nil
}

func (u *URLStorage) AddItem(itemID string, value string) {
	mutex.Lock()
	urls[itemID] = value
	mutex.Unlock()
}

func (u URLStorage) GetItems() map[string]string {
	return urls
}
