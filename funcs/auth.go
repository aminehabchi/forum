package forum

import (
	"net/http"
)

func Auth(funcNext http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Token")
		if err == nil {
			_, ok := GetUserIDFromToken(c.Value)
			if ok != nil {
				ClearSession(w)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				funcNext.ServeHTTP(w, r)
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

func AuthLG(funcNext http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Token")
		if err == nil {
			_, ok := GetUserIDFromToken(c.Value)
			if ok != nil {
				ClearSession(w)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			funcNext.ServeHTTP(w, r)
		}
	}
}

func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "Token",
		MaxAge: -1,
	})
}
