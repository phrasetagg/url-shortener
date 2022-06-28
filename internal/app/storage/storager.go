package storage

type IStorager interface {
	GetItem(itemID string) (string, error)
	AddItem(itemID string, value string, userID uint64)
	GetLastElementID() string
	GetItemsByUserID(userID uint64) []UserURLs
}
