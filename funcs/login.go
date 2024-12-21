package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Password string
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := LoginT.Execute(w, nil)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
		}
	case http.MethodPost:
		identifier := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if identifier == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := LoginT.Execute(w, "Please fill in all fields")
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
			}
			return
		}

		user, err := GetUserInfoByLoginInfo(identifier)
		if err != nil {
			if err == sql.ErrNoRows {
				err := LoginT.Execute(w, "Invalid credentials")
				if err != nil {
					ErrorHandler(w, http.StatusInternalServerError)
				}
				return
			}
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			err = LoginT.Execute(w, "Invalid credentials")
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
			}
			return
		}
		uuidStr, err := GenereteTocken()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		err = setToken(uuidStr, user.ID)
		if err != nil {
			fmt.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		newCookie := http.Cookie{
			Name:     "Token",
			Value:    uuidStr,
			Expires:  time.Now().Add(1 * time.Hour),
			HttpOnly: true,
		}

		http.SetCookie(w, &newCookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}

func GetUserInfoByLoginInfo(identifier string) (*User, error) {
	query := `SELECT id, password FROM users WHERE email = ? OR uname = ?`
	user := &User{}
	err := db.QueryRow(query, identifier, identifier).Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func setToken(token string, id int) error {
	query := "UPDATE tokens SET token=?,created_at=CURRENT_TIMESTAMP WHERE user_id=?"
	_, err := db.Exec(query, token, id)
	if err != nil {
		return err
	}
	return nil
}
