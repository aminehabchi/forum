package forum

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed1", http.StatusMethodNotAllowed)
		return
	}

	err := LoginT.Execute(w, nil)
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
}

func LoginInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	_, uname, correctPassword, err := GetUserInfoByLoginInfo(email)
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
	err = setLoginTime(uuidStr, uname)
	if err != nil {
		http.Error(w, "error in token", http.StatusInternalServerError)
		return
	}
	newCookie := http.Cookie{
		Name:     "Token",
		Value:    uuidStr,
		Expires:  time.Now().Add(1 * time.Minute),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, &newCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetUserInfoByLoginInfo(email_users string) (string, string, string, error) {
	query := `SELECT uname, password FROM users WHERE email = ?`
	var users, password string
	err := db.QueryRow(query, email_users).Scan(&users, &password)
	if err == nil {
		return email_users, users, password, nil
	}
	return "", "", "", errors.New("not exists")
}

func setLoginTime(token, uname string) error {
	query := "UPDATE users SET token=? WHERE uname=?"
	_, err := db.Exec(query, token, uname)

	return err
}
