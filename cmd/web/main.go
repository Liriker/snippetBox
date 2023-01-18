package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {

	// Создаём новый флаг командной строки.
	// Добавляем небольшую справку,объясняющую, что содержит данный флаг
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")

	// flag.Parse() извлекает флаг из командной строки и присваивает его содержимое
	flag.Parse()

	// Регистрируем обработчики и соответствующие url- шаблоны в маршрутизаторе
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Инициалзируем FileServer. Он будет обрабатывать
	// HTTP-запросы к статическим файлам в папке "./ui/static"
	// Путь является относительным корневой папки проекта
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})

	// Handle регистрирует обработчик
	// Он обрабатывает запросы к статическим файлам в папке" "./ui/static"
	// StripPrefix убирает префикс "/static"прежде чем запрос достигнет http.fileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// flag.Srtring
	log.Printf("Запуск веб-сервера на http://%s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

// Структура для проверки пути в http.FileSystem
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open - метод, проверяющий наличие файла index.html в папке path
// если файла index.html нет, то мы возвращаем 404 ошибку
// метод Open удовлетвоярет интерфейс FileSystem,
// что позволяет использовать тип neuteredFileSystem в http.FileServer
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	// Открываем указанный путь
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// Получаем информацию о вызываемом пути
	s, err := f.Stat()
	// Проверяем, является ли это папкой
	if s.IsDir() {
		// Если это папка, то мы проверяем существует ли "index.html" внутри данной папки
		index := filepath.Join(path, "index.html")
		//Если файла нет,то метод возвращает ошибку os.ErrNotExist
		if _, err := nfs.fs.Open(index); err != nil {
			// Закрываем файл index.html, что бы избежать утечки файлового дискриптора
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			// Возвращаем ошибку, которая будет преобразована http.FileServer в 404
			return nil, err
		}
	}
	// В остальных случаях просто возвращаем файл
	return f, nil
}
