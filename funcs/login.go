package forum

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var TokenMap = make(map[string][2]string)

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
	err = setLoginTime(1, uname)
	if err != nil {
		if err.Error() == "already login" {
			LoginT.Execute(w, "user already login")
		} else {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
		}
		return
	}
	uuidStr, err := GenereteTocken()
	if err != nil {
		http.Error(w, "err in token", http.StatusInternalServerError)
		return
	}
	TokenMap[uuidStr] = [2]string{email, uname}
	newCookie := http.Cookie{
		Name:     "Token",
		Value:    uuidStr,
		Expires:  time.Now().Add(24 * time.Hour),
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

func setLoginTime(bl int, uname string) error {
	isActive := 0
	query1 := `SELECT is_active FROM users WHERE uname=?`
	err := db.QueryRow(query1, uname).Scan(&isActive)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		fmt.Println("select is_active err", err)
		return err
	}
	if bl == 1 && isActive == 1 {
		return errors.New("already login")
	}
	query := "UPDATE users SET is_active=? WHERE uname=?"
	_, err = db.Exec(query, bl, uname)
	if err != nil {
		fmt.Println("err update", err)
		return err
	}
	return nil
}
