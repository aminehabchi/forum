package forum

import (
	"encoding/json"
	data "forum/funcs/database"
	Error "forum/funcs/error"
	"net/http"
	"strconv"
)

func LoadMoreComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	offsetValue := r.FormValue("offset")

	offset, err := strconv.Atoi(offsetValue)
	if err != nil {
		Error.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	post_id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || post_id <= 0 {
		Error.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	user_id := 0
	if Cookie, err := r.Cookie("Token"); err == nil {
		user_id, _ = data.GetUserIDFromToken(Cookie.Value)
	}

	comments, err := data.GetComment(post_id, user_id, 3, offset)
	if err != nil {
		Error.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		Error.ErrorHandler(w, http.StatusInternalServerError)
	}
}
