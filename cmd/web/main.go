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

	// Создаём новый флаг командной строки.
	// Добавляем небольшую справку,объясняющую, что содержит данный флаг
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")

	// flag.Parse() извлекает флаг из командной строки и присваивает его содержимое
	flag.Parse()

	// log.New создаёт новый логгер, в данном случае для записи информационных сообщений
	// три параметра: место назначения записи лого, префикс
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// Создаём логер для записи сообщений об ошибках
	// Как место для записи используем Stderr
	// log.Lshortfile включает в лог название файла и строку с ошибкой в нём
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Регистрируем обработчики и соответствующие url- шаблоны в маршрутизаторе
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Инициалзируем FileServer. Он будет обрабатывать
	// HTTP-запросы к статическим файлам в папке "./ui/static"
	// Путь является относительным корневой папки проекта
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})

	// Handle регистрирует обработчик
	// Он обрабатывает запросы к статическим файлам в папке" "./ui/static"
	// StripPrefix убирает префикс "/static"прежде чем запрос достигнет http.fileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// flag.Srtring
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
