package forum

import (
	"html/template"
	"net/http"
)

// function that display the error page
func ErrorHandler(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	tmpl, err := template.ParseFiles("./templates/error.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}

	var title, message string
	switch status {
	case http.StatusNotFound:
		title = "Page Not Found"
		message = "The page you are looking for does not exist."
	case http.StatusInternalServerError:
		title = "Internal Server Error"
		message = "Something went wrong on our server."
	case http.StatusMethodNotAllowed:
		title = "Method Not Allowed"
		message = "The HTTP method used is not allowed for this endpoint."
	case http.StatusBadRequest:
		title = "Bad Request"
		message = "somthing wrong from your part"
	default:
		title = "Error"
		message = "An unexpected error occurred."
	}
	data := struct{
		Status int
		Title, Message string
	}{
		status, title, message,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}
