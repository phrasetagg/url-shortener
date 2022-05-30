package storage

import (
	"errors"
)

var instance URLs

func init() {
	instance = URLs{URLs: map[string]string{}}
}

func GetURLsInstance() *URLs {
	return &instance
}

type URLs struct {
	URLs map[string]string
}

func (u URLs) GetItem(itemID string) (string, error) {
	item, ok := u.URLs[itemID]

	if !ok {
		return "", errors.New("not found")
	}

	return item, nil
}

func (u *URLs) AddItem(itemID string, value string) {
	u.URLs[itemID] = value
}

func (u URLs) GetItems() map[string]string {
	return u.URLs
}
