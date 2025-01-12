package forum

import (
	"database/sql"
	
	"net/http"
	"strings"
	"time"

	data "forum/funcs/database"
	Error "forum/funcs/error"
	types "forum/funcs/types"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := types.LoginT.Execute(w, nil)
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
		}
	case http.MethodPost:
		identifier := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if identifier == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := types.LoginT.Execute(w, "Please fill in all fields")
			if err != nil {
				Error.ErrorHandler(w, http.StatusInternalServerError)
			}
			return
		}

		user, err := data.GetUserInfoByLoginInfo(identifier)
		if err != nil {
			if err == sql.ErrNoRows {
				err := types.LoginT.Execute(w, "Invalid credentials")
				if err != nil {
					Error.ErrorHandler(w, http.StatusInternalServerError)
				}
				return
			}
			Error.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			err = types.LoginT.Execute(w, "Invalid credentials")
			if err != nil {
				Error.ErrorHandler(w, http.StatusInternalServerError)
			}
			return
		}
		uuidStr, err := data.GenereteTocken()
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		err = data.SetToken(uuidStr, user.ID)
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
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
		Error.ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}
