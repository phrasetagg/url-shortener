package models

import (
	"phrasetagg/url-shortener/internal/app/storage"
)

type Shortener struct {
	storage       storage.IStorager
	baseURL       string
	firstShortURL string
	lastShortURL  string
	maxCharCode   rune
}

func NewShortener(storage storage.IStorager, baseURL string) Shortener {
	return Shortener{
		storage:       storage,
		baseURL:       baseURL,
		firstShortURL: "a",
		maxCharCode:   rune(122), // Буква z
	}
}

func (s Shortener) GetFullURL(shortURL string) (string, error) {
	fullURL, err := s.storage.GetItem(shortURL)

	return fullURL, err
}

func (s Shortener) Shorten(userID uint32, URL string) string {
	shortURL := ""

	// Если короткая ссылка генерируется первый раз и при этом в хранилище нет ссылок,
	// то используем в качестве сокращенной ссылки firstShortURL.
	// Его же записываем в последнюю созданную сокращенную ссылку lastShortURL.
	// Добавляем все в хранилище.
	if s.lastShortURL == "" && s.storage.GetLastElementID() == "" {
		shortURL := s.firstShortURL
		s.lastShortURL = s.firstShortURL
		s.storage.AddItem(shortURL, URL, userID)

		return s.baseURL + shortURL
	}

	if s.storage.GetLastElementID() != "" {
		s.lastShortURL = s.storage.GetLastElementID()
	}

	// Разбиваем последнюю созданную короткую ссылку на коды.
	shortURLRune := []rune(s.lastShortURL)
	// Получаем код последнего символа короткой ссылки.
	lastCharCode := shortURLRune[len(shortURLRune)-1]

	// Если этот код равен коду максимально допустимого символа maxCharCode,
	// то конкатинируем в конец короткой ссылки символ firstShortURL.
	if lastCharCode == s.maxCharCode {
		shortURL = s.lastShortURL + s.firstShortURL
		s.lastShortURL = shortURL
		s.storage.AddItem(shortURL, URL, userID)

		return s.baseURL + shortURL
	}

	// Если код НЕ равен коду максимально допустимого символа maxCharCode,
	// то добавляем коду последнего символа 1, чтобы символ изменился на последующий.
	shortURLRune[len(shortURLRune)-1] = shortURLRune[len(shortURLRune)-1] + 1
	// Приводим к строке.
	shortURL = string(shortURLRune)
	s.lastShortURL = shortURL

	s.storage.AddItem(shortURL, URL, userID)

	return s.baseURL + shortURL
}

func (s Shortener) GetUserURLs(userID uint32) []storage.UserURLs {
	var preparedUserURLs []storage.UserURLs

	userURLs := s.storage.GetItemsByUserID(userID)

	for _, value := range userURLs {
		preparedUserURLs = append(preparedUserURLs, storage.UserURLs{ShortURL: s.baseURL + value.ShortURL, URL: value.URL})
	}

	return preparedUserURLs
}
