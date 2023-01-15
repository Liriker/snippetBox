package main

import (
	"log"
	"net/http"
)

// home - функция обработчик, сейчас имитирует ответ от пользователя
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Привет из SnippetBox"))
}

func main() {
	// http.NewServeMux - функция инициализации нового рутера
	// функция home регистрируется как обработчик URL шаблона "/"
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// http.ListenAndServe() - функция для запуска нового веб-сервера
	// Мы передаем два параметра: TCP-адрес сети для прослушивания
	//(в данном случае это "localhost:4000") и созданный рутер.
	// log.Fatal() используется  для логирования ошибок.
	//Любая ошибка, возвращаемая ListenAndServe() всегда non-nil.
	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
