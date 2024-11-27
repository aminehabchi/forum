package forum

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := LoginT.Execute(w, nil)
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
		}
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, id, correctPassword, err := GetUserInfoByLoginInfo(email)
		if err != nil {
			err = LoginT.Execute(w, "email not found")
			if err != nil {
				http.Error(w, "Could not load template", http.StatusInternalServerError)
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(correctPassword), []byte(password)); err != nil {
			err = LoginT.Execute(w, "password incorrect")
			if err != nil {
				http.Error(w, "Could not load template", http.StatusInternalServerError)
			}
			return
		}
		uuidStr, err := GenereteTocken()
		if err != nil {
			http.Error(w, "err in token", http.StatusInternalServerError)
			return
		}
		err = setLoginTime(uuidStr, id)
		if err != nil {
			http.Error(w, "error in token", http.StatusInternalServerError)
			return
		}
		newCookie := http.Cookie{
			Name:     "Token",
			Value:    uuidStr,
			Expires:  time.Now().Add(1 * time.Hour),
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(w, &newCookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
	}
}

func GetUserInfoByLoginInfo(email_users string) (string, int, string, error) {
	query := `SELECT id, password FROM users WHERE email = ?`
	var password string
	var id int
	err := db.QueryRow(query, email_users).Scan(&id, &password)
	if err == nil {
		return email_users, id, password, nil
	}
	return "", -1, "", errors.New("not exists")
}

func setLoginTime(token string, id int) error {
	query := "UPDATE users SET token=? WHERE id=?"
	_, err := db.Exec(query, token, id)

	return err
}
