package forum

import "net/http"

func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))

	if r.URL.Path == "/static/" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	static.ServeHTTP(w, r)
}
