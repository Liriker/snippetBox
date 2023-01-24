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
	// Описываем SQL запрос
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Exec выполняет запрос и возвращает sql.Result,
	// содержащий основные данные о том, что произошло после выполнения запроса
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// LastinsertId возвращает последний ID созданной записи в таблицу
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Возвращаемый ID имеет тип int64, конвертируем его и возвращаем результат
	return int(id), nil
}

// Get получает заметку по ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Lastest показывает последние 10 заметок
func (m *SnippetModel) Lastest() ([]*models.Snippet, error) {
	return nil, nil
}
