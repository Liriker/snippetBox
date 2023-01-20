package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// application хранит зависимости всего приложения
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Создаём флаг для командной строки, который указывает сетевой адрес.
	// По умолчанию адресс :4000
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	// Парсим адрес
	flag.Parse()

	// Создаём логи для информайции и ошибок
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Инициализируем структуру с зависимостями приложения
	// Указываем в созданные логи
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Инициализируем стуктуру сервера, что бы сервер использовал
	// указанные адрес, логи, и машрутизаторы
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Запускаем сервер, описывая соответствующие логи
	infoLog.Printf("Запуск веб-сервера на http://%s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
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
