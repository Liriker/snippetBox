package main

import (
	"log"
	"net/http"
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
	w.Write([]byte("Отображение заметки"))
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
		// Если это не так, то мы отправляем код ошибки 405
		// и тело ответа "Метод запрещён"
		// Затем вызываем return, что бы последующий код не выполнялся
		w.WriteHeader(405)
		w.Write([]byte("GET Метод запрещён!"))
		return
	}
	w.Write([]byte("Форма создания заметки"))
}

func main() {
	// Регистрируем обработчики и соответствующие url- шаблоны в маршрутизаторе
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
