package main

import (
	"log"
	"net/http"
)

func main() {
	// Регистрируем обработчики и соответствующие url- шаблоны в маршрутизаторе
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Инициалзируем FileServer. Он будет обрабатывать
	// HTTP-запросык статическим файлам в папке "./ui/static"
	// Путь является относительным корневой папки проекта
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Handle регистрирует обработчик
	// Он обрабатывает запросы к статическим файлам в папке" "./ui/static"
	// StripPrefix убирает префикс "/static"прежде чем запрос достигнет http.fileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
