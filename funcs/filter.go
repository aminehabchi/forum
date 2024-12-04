package forum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	filter := strings.ToLower(r.URL.Query().Get("type"))

	if filter != "" && !allCategories[strings.ToLower(filter)] &&
		filter != "created" && filter != "liked" {
			ErrorHandler(w,http.StatusBadRequest)
			return
		}

	userID := 0
	if cookie, err := r.Cookie("Token"); err == nil {
		userID, _ = GetUserIDFromToken(cookie.Value)
	}

	if (filter == "created" || filter == "liked") && userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	opts := QueryOptions{
		UserID: userID,
		Filter: filter,
	}

	query, args := BuildPostQuery(opts)
	posts, err := GetPosts(userID, query, args...)
	if err != nil && err != sql.ErrNoRows {
		ErrorHandler(w, http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
	}
}
