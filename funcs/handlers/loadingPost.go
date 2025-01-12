package forum

import (
	"database/sql"
	"encoding/json"

	data "forum/funcs/database"
	Error "forum/funcs/error"
	types "forum/funcs/types"
	"net/http"
	"strconv"
)

func LoadMorePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	offsetValue := r.FormValue("offset")
	filterType := r.FormValue("type")

	offset, err := strconv.Atoi(offsetValue)
	if err != nil {
		Error.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	user, err := r.Cookie("Token")

	var user_id int
	if err == nil {
		user_id, _ = data.GetUserIDFromToken(user.Value)
	}

	opts := types.QueryOptions{
		UserID: user_id,
		Limit:  4,
		Offset: offset,
		Filter: filterType,
	}

	query, args := data.BuildPostQuery(opts)

	posts, err := data.GetPosts(user_id, query, args...)
	if err != nil && err != sql.ErrNoRows {
		Error.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		Error.ErrorHandler(w, http.StatusInternalServerError)
	}
}
