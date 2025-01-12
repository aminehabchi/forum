package forum

import (
	"database/sql"
	"encoding/json"
	data "forum/funcs/database"
	Error "forum/funcs/error"
	types "forum/funcs/types"
	"log"
	"net/http"
	"strings"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	filter := strings.ToLower(r.URL.Query().Get("type"))

	if filter != "" && !data.AllCategories[strings.ToLower(filter)] &&
		filter != "created" && filter != "liked" {
		Error.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	userID := 0
	if cookie, err := r.Cookie("Token"); err == nil {
		userID, _ = data.GetUserIDFromToken(cookie.Value)
	}

	if (filter == "created" || filter == "liked") && userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	opts := types.QueryOptions{
		UserID: userID,
		Filter: filter,
	}

	query, args := data.BuildPostQuery(opts)
	posts, err := data.GetPosts(userID, query, args...)
	if err != nil && err != sql.ErrNoRows {
		Error.ErrorHandler(w, http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		Error.ErrorHandler(w, http.StatusInternalServerError)
	}
}
