package models

import (
	"errors"
	"time"
)

// Ошибка при отсутствии искомой записи
var ErrNoRecord = errors.New("models: подходящей записи не найдено")

// Модель заметки
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
