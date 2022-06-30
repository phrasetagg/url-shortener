package storage

import (
	"context"
	"errors"
	"phrasetagg/url-shortener/internal/app/db"
	"time"
)

type DBURLStorage struct {
	db *db.DB
}

func NewDBURLStorage(db *db.DB) *DBURLStorage {
	return &DBURLStorage{
		db: db,
	}
}

func (d *DBURLStorage) GetItem(itemID string) (string, error) {
	var originalUrl string

	conn, err := d.db.GetConn(context.Background())
	if err != nil {
		return "", err
	}

	defer d.db.Close()

	err = conn.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE short_url = $1 LIMIT 1", itemID).Scan(&originalUrl)
	if err != nil {
		panic(err)
	}

	if originalUrl == "" {
		return "", errors.New("not found")
	}

	return originalUrl, nil
}

func (d *DBURLStorage) AddItem(itemID string, value string, userID uint32) {
	conn, err := d.db.GetConn(context.Background())
	if err != nil {
		panic(err)
	}

	defer d.db.Close()

	_, err = conn.Exec(context.Background(), "INSERT INTO urls (short_url, original_url, user_id, created_at) VALUES ($1,$2,$3,$4)", itemID, value, userID, time.Now())
	if err != nil {
		panic(err)
	}
}

func (d DBURLStorage) GetLastElementID() string {
	var lastElementID string

	conn, err := d.db.GetConn(context.Background())
	if err != nil {
		return ""
	}

	defer d.db.Close()

	err = conn.QueryRow(context.Background(), "SELECT short_url FROM urls ORDER BY created_at DESC LIMIT 1").Scan(&lastElementID)

	if err != nil && err.Error() == "no rows in result set" {
		return ""
	}

	if err != nil {
		panic(err)
	}

	return lastElementID
}

func (d DBURLStorage) GetItemsByUserID(userID uint32) []UserURLs {
	userURLs := make([]UserURLs, 0)

	conn, err := d.db.GetConn(context.Background())
	if err != nil {
		return userURLs
	}

	defer d.db.Close()

	rows, err := conn.Query(context.Background(), "SELECT short_url, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var urlData UserURLs

		err := rows.Scan(&urlData.ShortURL, &urlData.URL)
		if err != nil {
			return nil
		}

		userURLs = append(userURLs, urlData)
	}

	return userURLs
}
