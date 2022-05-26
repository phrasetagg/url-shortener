package storage

type Storager interface {
	GetItem(itemID string) (string, error)
	AddItem(itemID string, value string)
}
