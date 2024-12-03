package forum

import "net/http"

func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))

	if r.URL.Path == "/static/" {
		ErrorHandler(w, 403)
		return
	}

	static.ServeHTTP(w, r)
}
