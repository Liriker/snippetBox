package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// home - обработчик главной страницы
func home(w http.ResponseWriter, r *http.Request) {
	// Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/". Если нет, вызывается
	// функция http.NotFound() для возвращения клиенту ошибки 404.
	// Важно, чтобы мы завершили работу обработчика через return. Если мы забудем про "return", то обработчик
	// продолжит работу и выведет сообщение "Привет из SnippetBox" как ни в чем не бывало.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Привет из SnippetBox"))
}

// showSnippet - обработчик для отображения содержимого заметки
func showSnippet(w http.ResponseWriter, r *http.Request) {

	// Извлекаем значение параметра id из URL и пытаемся конвертировать строку в int
	// используя функцию strconv.Atoi.
	// Если конвертирование не удалось или значение меньше 1,
	// то возвращаем 404 - Страница не найдена.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Fprintf вставляет значение id в строку
	// и записывает его в w - http.ResponseWriter.
	fmt.Fprintf(w, "Отображение заметки с ID %d...", id)
}

// createSnippet - обработчик создания новой заметки
func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Проверяем, я вляется ли запрос POST
	// http.MethodPost является строкой и содержит POST
	if r.Method != http.MethodPost {
		// Header().Set() добавляет заголовок "Allow: POST"
		// в карту http заголовков.
		// Первый параметр -название заголовка, второй - значение
		w.Header().Set("Allow", http.MethodPost)

		// http.Error() отправляет код состояния с телом ошибки
		// Под капотом тут так же есть w.Write и w.WriteHeader
		http.Error(w, "Метод запрещён!", 405)
		return
	}
	w.Write([]byte("Форма создания заметки"))
}
