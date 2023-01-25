package main

import (
	"Liriker/snippetBox/pkg/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Инициализируем карту, которая будет хранить кеш
	cache := map[string]*template.Template{}

	// Получаем срез всех файловых путей "*.page.html"
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Извлекаем конечное название файла и присваиваем его переменной
		name := filepath.Base(page)

		// Обрабатываем итерируемый файл шаблона
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// добавляем все каркасные шаблоны
		// В нашем случае это base.layout.html
		ts, err = ts.ParseGlob(filepath.Join("*.layout.html"))
		if err != nil {
			return nil, err
		}

		// Добавляем вспомагательные шаблоны
		// В нашем случае это footer.partial.html
		ts, err = ts.ParseGlob(filepath.Join("*.partial.html"))
		if err != nil {
			return nil, err
		}

		// Добавляем полученный набор шаблонов в кэш
		cache[name] = ts

	}

	// Возвращаем полученную карту
	return cache, nil
}
