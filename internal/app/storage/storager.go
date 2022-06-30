package storage

type IStorager interface {
	GetItem(itemID string) (string, error)
	AddItem(itemID string, value string, userID uint32)
	GetLastElementID() string
	GetItemsByUserID(userID uint32) []UserURLs
}
