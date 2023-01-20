package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError записывает сообщение об ошибке в errorLog и
// отправляет пользователю ответ 500
func (app *application) serverError(w http.ResponseWriter, err error) {
	// debag.Stack используется для получения трассировки стека текущей горутины и добавления её в логгер
	trace := fmt.Sprintf("%s\n%s,", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError отправляет определённый код состояния пользователю
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound - оболочка для clientError для отправки 404 ошибки
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
