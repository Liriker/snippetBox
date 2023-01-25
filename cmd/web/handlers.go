package main

import (
	"Liriker/snippetBox/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// home - обработчик главной страницы
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что бы домашняя страница вызывалась только при URL .../
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Получаем данные последних 10-ти заметок
	s, err := app.snippets.Lastest()
	if err != nil {
		app.serverError(w, err)
	}

	// Отображаем шаблон страницы с данными
	app.render(w, r, "home.page.html", &templateData{Snippets: s})
}

// showSnippet - обработчик для отображения содержимого заметки
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// Извлекаем значение параметра id из URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	// Извлеккаем данные из записи по её ID
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Отображаем шаблон с данными
	app.render(w, r, "show.page.html", &templateData{Snippet: s})
}

// createSnippet - обработчик создания новой заметки
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Проверяем, я вляется ли запрос POST
	if r.Method != http.MethodPost {
		// Добавляем заголовок "Allow: POST" в карту http заголовков.
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Временно создаём переменные для теста
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

	// Передаём данные в Insert для создания заметки, получаем id созданной заметки
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправляем пользователя на соответствующую страницу заметки
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
