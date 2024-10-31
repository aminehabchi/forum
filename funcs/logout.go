package forum

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	c, _ := r.Cookie("Token")
	uname := TokenMap[c.Value][1]
	delete(TokenMap, c.Value)
	cookie := http.Cookie{
		Name:   "Token",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	err := setLoginTime(0, uname)
	if err != nil {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
