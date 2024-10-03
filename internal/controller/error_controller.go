package controller

import (
	"html/template"
	"log/slog"
	"net/http"
)

func ErrorController(w http.ResponseWriter, statusCode int, Info string) {
	tmpl := template.Must(template.ParseFiles(GetTmpPath("error")))

	w.WriteHeader(statusCode)
	if err := tmpl.Execute(w, Info); err != nil {
		slog.Error(err.Error())
	}
}
