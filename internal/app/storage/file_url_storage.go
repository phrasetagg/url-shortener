package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type FileURLStorage struct {
	filePath string
	mutex    *sync.RWMutex
}

func NewFileURLStorage(filePath string) *FileURLStorage {
	return &FileURLStorage{
		filePath: filePath,
		mutex:    new(sync.RWMutex),
	}
}

func (s *FileURLStorage) AddRecord(itemID string, value string, userID uint32) error {
	file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(itemID + " " + value + " " + fmt.Sprint(userID) + "\n")
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (s *FileURLStorage) GetOriginalURLByShortURI(itemID string) (OriginalURL, error) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0777)

	var originalURL OriginalURL

	if err != nil {
		return originalURL, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	reader := bufio.NewReader(file)

	for {
		item, err := reader.ReadBytes('\n')

		row := string(item)
		row = strings.Trim(row, "\n")
		res := strings.Split(row, " ")

		if res[0] == itemID {
			originalURL.OriginalURL = res[1]
			return originalURL, nil
		}

		if err == io.EOF {
			break
		}
	}

	return originalURL, errors.New("not found")
}

func (s FileURLStorage) GetShortURIByOriginalURL(originalURL string) (string, error) {
	return originalURL, errors.New("file_url_storage doesn't support this method")
}

func (s FileURLStorage) GetRecordsByUserID(userID uint32) []UserURLs {

	var userURLs []UserURLs

	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return userURLs
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	reader := bufio.NewReader(file)

	for {
		item, err := reader.ReadBytes('\n')

		row := string(item)
		row = strings.Trim(row, "\n")
		res := strings.Split(row, " ")

		if err == nil && res[2] == fmt.Sprint(userID) {
			userURLs = append(userURLs, UserURLs{ShortURL: res[0], URL: res[1]})
		}

		if err == io.EOF {
			break
		}
	}

	return userURLs
}

func (s FileURLStorage) DeleteUserRecordsByShortURLs(_ uint32, _ []string) error {
	return errors.New("file_url_storage doesn't support this method")
}
