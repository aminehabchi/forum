package forum

import (
	"net/http"
)

func Auth(funcNext http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Token")
		if err == nil {
			_, ok := GetUserNameFromToken(c.Value)
			if ok != nil {
				cookie := http.Cookie{
					Name:   "Token",
					MaxAge: -1,
				}
				http.SetCookie(w, &cookie)
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
			_, ok := GetUserNameFromToken(c.Value)
			if ok != nil {
				cookie := http.Cookie{
					Name:   "Token",
					MaxAge: -1,
				}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			funcNext.ServeHTTP(w, r)
		}
	}
}
