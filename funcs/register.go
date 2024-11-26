package forum

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	err := RegisterT.Execute(w, nil)
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
	email := r.FormValue("email")
	uname := r.FormValue("uname")
	password := r.FormValue("password")

	if email == "" || uname == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := RegisterT.Execute(w, "Invalid Inputs, Please fill all inputs")
		if err != nil {
			http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	err = InsertUserInfo(email, string(hashedPassword), uname)
	if err != nil && err != sql.ErrNoRows {
		err = RegisterT.Execute(w, "user name or email already used")
		if err != nil {
			http.Error(w, "500 Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func InsertUserInfo(email, password, uname string) error {
	selector := `INSERT INTO users(password,uname,email) VALUES (?,?,?)`
	_, err := db.Exec(selector, password, uname, email)
	if err != nil {
		return err
	}
	return nil
}
