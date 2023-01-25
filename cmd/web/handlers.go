package main

import (
	"Liriker/snippetBox/pkg/models"
	"errors"
	"fmt"
	"html/template"
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

	s, err := app.snippets.Lastest()
	if err != nil {
		app.serverError(w, err)
	}

	data := &templateData{Snippets: s}

	//Инициализируем срез, содержащий пути к файлам.
	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partical.html",
	}

	// Парсим шаблоны из среза
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	// Записываем шаблоны и данные в тело ответа
	err = ts.Execute(w, data)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
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

	data := &templateData{Snippet: s}

	// Инициализируем срез строк, у котором будут храниться пути к html шаблонам
	files := []string{
		"./ui/html/show.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partical.html",
	}

	// Парсим шаблоны
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Выполняем шаблоны, в качестве динамических данных указываем заметки
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// createSnippet - обработчик создания новой заметки
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Проверяем, я вляется ли запрос POST
	// http.MethodPost является строкой и содержит POST
	if r.Method != http.MethodPost {
		// Header().Set() добавляет заголовок "Allow: POST"
		// в карту http заголовков.
		// Первый параметр -название заголовка, второй - значение
		w.Header().Set("Allow", http.MethodPost)
		// http.Error() отправляет код состояния с телом ошибки
		// Под капотом тут так же есть w.Write и w.WriteHeader
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
