package forum

import (
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	Status     string
	StatusText string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Status:     r.FormValue("s"),
		StatusText: r.FormValue("st"),
	}

	t, _ := template.ParseFiles("templates/error.html")
	statusCode, _ := strconv.Atoi(data.Status)
	w.WriteHeader(statusCode)
	t.Execute(w, data)
}
