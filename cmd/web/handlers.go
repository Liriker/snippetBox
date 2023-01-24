package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// home - обработчик главной страницы
// сигнатура application опередяет home как метод, что позволяет использовать зависимости
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Инициализируем срез, содержащий пути к файлам.
	// ВАЖНО home.page.html должен быть первым в срезе
	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partical.html",
	}

	// template.ParseFiles читает файл шаблона
	// В случае возврата ошибки записываем в лог детальное сообщение ошибки с помощью http.Error()
	// и отправляем пользователю ответ в виде 500 ошибки
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Используем логгер из структуры application вместо стандартного
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	// Execute записывает содержимое шаблона в тело HTTP ответа
	// Последний параметр предоставлет возможность отправки динамических данных в шаблон
	err = ts.Execute(w, nil)
	if err != nil {
		// Так же используем логер из структры
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

// showSnippet - обработчик для отображения содержимого заметки
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// Извлекаем значение параметра id из URL и пытаемся конвертировать строку в int
	// используя функцию strconv.Atoi.
	// Если конвертирование не удалось или значение меньше 1,
	// то возвращаем 404 - Страница не найдена.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Fprintf вставляет значение id в строку
	// и записывает его в w - http.ResponseWriter.
	fmt.Fprintf(w, "Отображение заметки с ID %d...", id)
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
