package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"phrasetagg/url-shortener/internal/app/db"
	"strings"
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

func (s *DBURLStorage) AddRecord(itemID string, value string, userID uint32) error {
	conn, err := s.db.GetConn(context.Background())
	if err != nil {
		panic(err)
	}

	defer s.db.Close()

	_, err = conn.Exec(context.Background(), "INSERT INTO urls (short_uri, original_url, user_id, created_at) VALUES ($1,$2,$3,$4)", itemID, value, userID, time.Now())

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return &ItemAlreadyExistsError{value: value}
	}

	if err != nil {
		panic(err)
	}

	return nil
}

func (s *DBURLStorage) GetOriginalURLByShortURI(itemID string) (OriginalURL, error) {
	var originalURL OriginalURL

	conn, err := s.db.GetConn(context.Background())
	if err != nil {
		return originalURL, err
	}

	defer s.db.Close()

	err = conn.QueryRow(context.Background(), "SELECT original_url, is_deleted FROM urls WHERE short_uri = $1 LIMIT 1", itemID).Scan(&originalURL.OriginalURL, &originalURL.IsDeleted)
	if err != nil {
		panic(err)
	}

	if originalURL.OriginalURL == "" {
		return originalURL, errors.New("not found")
	}

	return originalURL, nil
}

func (s DBURLStorage) GetShortURIByOriginalURL(originalURL string) (string, error) {
	var shortURI string

	conn, err := s.db.GetConn(context.Background())
	if err != nil {
		return "", err
	}

	defer s.db.Close()

	err = conn.QueryRow(context.Background(), "SELECT short_uri FROM urls WHERE original_url = $1 LIMIT 1", originalURL).Scan(&shortURI)
	if err != nil {
		panic(err)
	}

	if shortURI == "" {
		return "", errors.New("not found")
	}

	return shortURI, nil
}

func (s DBURLStorage) GetRecordsByUserID(userID uint32) []UserURLs {
	userURLs := make([]UserURLs, 0)

	conn, err := s.db.GetConn(context.Background())
	if err != nil {
		return userURLs
	}

	defer s.db.Close()

	rows, err := conn.Query(context.Background(), "SELECT short_uri, original_url FROM urls WHERE user_id = $1", userID)
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

func (s DBURLStorage) DeleteUserRecordsByShortURLs(userID uint32, shortURLs []string) error {
	conn, err := s.db.GetConn(context.Background())
	if err != nil {
		return err
	}

	defer s.db.Close()

	_, err = conn.Exec(context.Background(), "UPDATE urls SET is_deleted = true WHERE short_uri = ANY($1) AND user_id = $2", shortURLs, userID)
	if err != nil {
		return err
	}

	return nil
}
