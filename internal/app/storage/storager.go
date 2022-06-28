package storage

type IStorager interface {
	GetItem(itemID string) (string, error)
	AddItem(itemID string, value string)
	GetLastElementID() string
}
