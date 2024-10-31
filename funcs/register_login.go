package forum

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type POST struct {
	ID       int
	Name     string
	Title    string
	Content  string
	Category []string
	Likes    int
	Dislikes int
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	uname, _ := r.Cookie("username")
	cookie := http.Cookie{
		Name:   "username",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	err := setLoginTime(0, uname.Value)
	if err != nil {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	_, err := r.Cookie("username")
	if err != http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = RegisterT.Execute(w, nil)
	if err != nil {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		return
	}
}

func RegisterIngo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	_, err := r.Cookie("username")
	if err != http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	uname := r.FormValue("uname")
	password := r.FormValue("password")

	if email == "" || uname == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err = RegisterT.Execute(w, "Invalid Inputs, Please fill all inputs")
		if err != nil {
			http.Error(w, "500 Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	err = InsertUserInfo(email, string(hashedPassword), uname)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		err = RegisterT.Execute(w, "user name or email already used")
		if err != nil {
			http.Error(w, "500 Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed1", http.StatusMethodNotAllowed)
		return
	}
	_, err := r.Cookie("username")
	if err != http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err = LoginT.Execute(w, nil)
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
	_, err := r.Cookie("username")
	if err != http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	if err := setLoginTime(1, uname); err != nil {
		err = LoginT.Execute(w, "already user is login")
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
		}
		return
	}

	newCookie := http.Cookie{
		Name:     "username",
		Value:    uname,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &newCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func GetUserInfoByLoginInfo(email_users string) (string, string, string, error) {
	query := `SELECT uname, password FROM users WHERE email = ?`
	var users, password string
	err := db.QueryRow(query, email_users).Scan(&users, &password)
	if err == nil {
		return email_users, users, password, nil
	}
	return "", "", "", errors.New("not exists")
}

func InsertUserInfo(email, password, uname string) error {
	selector := `INSERT INTO users(password,uname,email) VALUES (?,?,?)`
	_, err := db.Exec(selector, password, uname, email)
	if err != nil {
		return err
	}
	return nil
}
