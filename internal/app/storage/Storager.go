package storage

type Storager interface {
	GetItem(itemId string) (string, error)
	AddItem(itemId string, value string)
}
