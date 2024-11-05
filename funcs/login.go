package forum

import (
	"errors"
	"fmt"
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
	err = setLoginTime(1, uuidStr, uname)
	if err != nil {
		if err.Error() == "already login" {
			LoginT.Execute(w, "user already login")
		} else {
			http.Error(w, "Could not load template1", http.StatusInternalServerError)
		}
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

func setLoginTime(bl int, token, uname string) error {
	isActive := 0
	loginTime := ""
	query := `SELECT is_active,tokenTime FROM users WHERE uname=?`
	err := db.QueryRow(query, uname).Scan(&isActive, &loginTime)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if bl == 1 {
		if isActive == 1 && loginTime != "00-00-0000" {
			t, _ := stringToTime(loginTime)
			tt := t.Add(1 * time.Minute)
			if tt.Before(time.Now()) {
				return errors.New("already login")
			}
		}
		query = "UPDATE users SET token=?,is_active=?,tokenTime=? WHERE uname=?"
		_, err = db.Exec(query, token, 1, timeToString(time.Now()), uname)
		if err != nil {
			return err
		}
	} else {
		query = "UPDATE users SET is_active=? WHERE uname=?"
		_, err = db.Exec(query, 0, uname)
		if err != nil {
			return err
		}
	}
	return nil
}

func timeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func stringToTime(s string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, s)
}
