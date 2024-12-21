package forum

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := RegisterT.Execute(w, nil)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
		}
	case http.MethodPost:
		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		uname := strings.ToLower(strings.TrimSpace(r.FormValue("uname")))
		password := r.FormValue("password")

		// if err := RegisterValidation(email, uname, password); err != "" {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	err := RegisterT.Execute(w, err)
		// 	if err != nil {
		// 		ErrorHandler(w, http.StatusInternalServerError)
		// 	}
		// 	return
		// }

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		err = InsertUserInfo(email, string(hashedPassword), uname)
		if err != nil {
			errMsg := err.Error()
			switch {
			case strings.Contains(errMsg, "UNIQUE constraint failed: users.uname"):
				if execErr := RegisterT.Execute(w, "This username is already taken. Please choose a different one."); execErr != nil {
					log.Printf("Template execution error: %v", execErr)
					ErrorHandler(w, http.StatusInternalServerError)
				}
				return

			case strings.Contains(errMsg, "UNIQUE constraint failed: users.email"):
				if execErr := RegisterT.Execute(w, "This email address is already registered. Please use a different email."); execErr != nil {
					log.Printf("Template execution error: %v", execErr)
					ErrorHandler(w, http.StatusInternalServerError)
				}
				return

			default:
				log.Printf("Database error during registration: %v", err)
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}

func RegisterValidation(email, uname, password string) string {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if email == "" || !emailRegex.MatchString(email) {
		return "Please enter a valid email address"

	}
	if len(uname) < 3 || len(uname) > 30 {
		return "Username must be between 3 and 30 characters"

	}

	if len(password) < 8 {
		return "Password must be at least 8 characters long"

	}
	return ""
}

func InsertUserInfo(email, password, uname string) error {
	selector := `INSERT INTO users(password,uname,email) VALUES (?,?,?)`
	result, err := db.Exec(selector, password, uname, email)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	selector = `INSERT INTO tokens(user_id) VALUES (?)`
	_, err = db.Exec(selector, int(id))
	if err != nil {
		return err
	}
	return nil
}
