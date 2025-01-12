package forum

import (
	Data "forum/funcs/database"
	"net/http"
)

func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "Token",
		MaxAge: -1,
	})
}

func Auth(funcNext http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, bl := CheckIfCookieValid(w, r); bl {
			funcNext(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

func AuthLG(funcNext http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, bl := CheckIfCookieValid(w, r); bl {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			funcNext(w, r)
		}
	}
}
func CheckIfCookieValid(w http.ResponseWriter, r *http.Request) (int, bool) {
	var userId int
	c, err := r.Cookie("Token")
	if err == nil {
		userId, err = Data.GetUserIDFromToken(c.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "Token",
				MaxAge: -1,
			})
			return userId, false
		} else {
			return userId, true
		}
	}
	return userId, false
}
