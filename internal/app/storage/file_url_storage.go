package storage

import (
	"bufio"
	"errors"
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

func (s *FileURLStorage) GetItem(itemID string) (string, error) {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
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
			return res[1], nil
		}

		if err == io.EOF {
			break
		}
	}

	return "", errors.New("not found")
}

func (s *FileURLStorage) AddItem(itemID string, value string) {
	file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(itemID + " " + value + "\n")
	if err != nil {
		return
	}

	err = writer.Flush()
	if err != nil {
		return
	}
}

func (s FileURLStorage) GetLastElementID() string {
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return ""
	}

	reader := bufio.NewReader(file)

	var row string

	for {
		var item []byte

		item, err = reader.ReadBytes('\n')

		if err == nil {
			row = string(item)
		}

		if err == io.EOF {
			break
		}
	}

	return strings.Split(row, " ")[0]
}
