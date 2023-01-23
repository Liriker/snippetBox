package mysql

import (
	"Liriker/snippetBox/pkg/models"
	"database/sql"
)

// SnippetModel - тип, который оборачивает пул подключения sql.DB
type SnippetModel struct {
	DB *sql.DB
}

// Insert создаёт заметку
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// Get получает заметку по ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Lastest показывает последние 10 заметок
func (m *SnippetModel) Lastest() ([]*models.Snippet, error) {
	return nil, nil
}
