package forum

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

type USER struct {
	Unknown bool
	Uname   string
	Message string
}

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
	cookie := http.Cookie{
		Name:   "username",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/register.html")
	temp.Execute(w, nil)
}

func RegisterIngo(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	uname := r.FormValue("uname")
	password := r.FormValue("password")
	if email == "" || uname == "" || password == "" {
		// badrequest
	}
	err := InsertUserInfo(email, password, uname)
	if err != nil {
		fmt.Println(err)
		temp, _ := template.ParseFiles("register.html")
		temp.Execute(w, "user name or email already used")
	} else {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/login.html")
	temp.Execute(w, nil)
}

func LoginInfo(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/login.html")
	email_uname := r.FormValue("email")
	password := r.FormValue("password")
	_, uname, correctPassword, err := checkInfo(email_uname)
	if err != nil {
		temp.Execute(w, "can t find user")
	} else if correctPassword != password {
		temp.Execute(w, "password incorrect")
	} else {
		c, err := r.Cookie("username")
		if err != nil || (err == nil && c.Value == uname) {
			cookie := http.Cookie{
				Name:  "username",
				Value: uname,
				// timeeeee
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else if err == nil || c.Value != uname {
			temp.Execute(w, "already user is login")
		}
	}
}

func checkInfo(email_uname string) (string, string, string, error) {
	email, uname, correctPassword, err := GetUserInfoByLoginInfo(email_uname)
	if err != nil {
		return "", "", "", err
	}
	return email, uname, correctPassword, nil
}

func GetUserInfoByLoginInfo(email_users string) (string, string, string, error) {
	query := `SELECT email, password FROM users WHERE uname = ?`
	var users, password string
	err := db.QueryRow(query, email_users).Scan(&users, &password)
	if err == nil {
		return email_users, users, password, nil
	}
	query = `SELECT uname, password FROM users WHERE email = ?`
	var email string
	err = db.QueryRow(query, email_users).Scan(&email, &password)
	if err == nil {
		return email, email_users, password, nil
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
