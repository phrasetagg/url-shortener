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

func (u URLs) GetItem(itemId string) (string, error) {
	item, ok := u.URLs[itemId]

	if !ok {
		return "", errors.New("not found")
	}

	return item, nil
}

func (u *URLs) AddItem(itemId string, value string) {
	u.URLs[itemId] = value
}
