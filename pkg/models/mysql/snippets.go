package mysql

import (
	"Liriker/snippetBox/pkg/models"
	"database/sql"
	"errors"
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
	// sql запрос
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// QueryRow выполняет запрос
	// Возвращается указатель на объект sql.Row, содержащий данные записи
	row := m.DB.QueryRow(stmt, id)

	// Инициализируем указатль на новую структуру Snippet
	s := &models.Snippet{}

	// Копируем значения из каждого поля в соответствующее поле в структуре Snippet
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// проверяем на наличие данного ID
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
			//	Если ошибка другая, то просто возвращаем её
		} else {
			return nil, err
		}
	}

	// Если всё в порядке, то возвращаем объект Snippet
	return s, nil

}

// Lastest показывает последние 10 заметок
func (m *SnippetModel) Lastest() ([]*models.Snippet, error) {

	// Пишем SQL запрос
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// Query() выполняет запрос, возвращает sql.Rows,
	// содержащий результат нашего запроса
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Выполняем закрытие набора результатов при любом исходе метода Lastest()
	defer rows.Close()

	// Инициализируем пустой срез для хранения результатов,
	// которые мы выведем
	var snippets []*models.Snippet

	// Перебераем строки результата запроса
	for rows.Next() {
		// Инициализируем указатель на новую структуру,
		// которую мы будем добавлять в итоговый массив
		s := &models.Snippet{}

		// Копируем значения полей в структуру
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Добавляем полученную заполненную структуру в итоговый массив
		snippets = append(snippets, s)
	}

	// Узнаём, возникла ли какая-то ошибка в ходе работы
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если всё в норме, то возвращаем заполненый массив
	return snippets, nil

}
