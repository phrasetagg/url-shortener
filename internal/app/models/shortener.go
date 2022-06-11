package models

import (
	"phrasetagg/url-shortener/internal/app/storage"
)

type Shortener struct {
	storage storage.Storager
}

const HostURL = "http://localhost:8080/"

var (
	firstShortURL = "a"
	lastShortURL  string
	maxCharCode   = rune(122) // Буква z
)

func NewShortener(storage storage.Storager) Shortener {
	return Shortener{
		storage: storage,
	}
}

func (s Shortener) GetFullURL(shortURL string) (string, error) {
	fullURL, err := s.storage.GetItem(shortURL)

	return fullURL, err
}

func (s Shortener) Shorten(URL string) string {
	shortURL := ""

	// Если мапа пустая (первый запрос после запуска), то используем в качестве сокращенной ссылки firstShortURL.
	// Его же записываем в последнюю созданную сокращенную ссылку lastShortURL.
	// Добавляем все в мапу.
	if len(s.storage.GetItems()) == 0 {
		shortURL := firstShortURL
		lastShortURL = firstShortURL
		s.storage.AddItem(shortURL, URL)

		return HostURL + shortURL
	}

	// Разбиваем последнюю созданную короткую ссылку на коды.
	shortURLRune := []rune(lastShortURL)
	// Получаем код последнего символа короткой ссылки.
	lastCharCode := shortURLRune[len(shortURLRune)-1]

	// Если этот код равен коду максимально допустимого символа maxCharCode,
	// то конкатинируем в конец короткой ссылки символ firstShortURL.
	if lastCharCode == maxCharCode {
		shortURL = lastShortURL + firstShortURL
		lastShortURL = shortURL
		s.storage.AddItem(shortURL, URL)

		return HostURL + shortURL
	}

	// Если код НЕ равен коду максимально допустимого символа maxCharCode,
	// то добавляем коду последнего символа 1, чтобы символ изменился на последующий.
	shortURLRune[len(shortURLRune)-1] = shortURLRune[len(shortURLRune)-1] + 1
	// Приводим к строке.
	shortURL = string(shortURLRune)
	lastShortURL = shortURL

	s.storage.AddItem(shortURL, URL)

	return HostURL + shortURL
}
