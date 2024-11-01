package forum

import (
	"net/http"
)

func Auth(funcNext http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Token")
		_, ok := TokenMap[c.Value]
		if err != nil || !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			funcNext.ServeHTTP(w, r)
		}
	}
}

func AuthLG(funcNext http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Token")
		if err == nil {
			_, ok := TokenMap[c.Value]
			if !ok {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			funcNext.ServeHTTP(w, r)
		}
	}
}