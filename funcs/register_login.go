package forum

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"
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
	uname, _ := r.Cookie("username")
	cookie := http.Cookie{
		Name:   "username",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	err := setLoginTime(0, uname.Value)
	if err != nil {
		fmt.Println(err)
	}
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
			err = setLoginTime(1, uname)
			if err != nil {
				temp.Execute(w, "already user is login")
			}else{
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			}
		} else if err == nil || c.Value != uname {
			temp.Execute(w, "already user is login")
		}
	}
}

func setLoginTime(bl int, uname string) error {
	isActive := 0
	query1 := `SELECT is_active FROM users WHERE uname=?`
	db.QueryRow(query1, uname).Scan(&isActive)

	if bl == 1 && isActive == 1 {
		return errors.New("already login")
	}
	currentTime := time.Now()
	query := "UPDATE users SET is_active=?, loginTime=? WHERE uname=?"
	_, err := db.Exec(query, bl, currentTime.Format("2006-01-02 15:04:05"), uname)
	if err != nil {
		return err
	}
	return nil
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
