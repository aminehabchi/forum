package forum

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "Token",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
