package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.FormValue("type")
	fmt.Println(filter)
	if !allCategories[filter] {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := 0
	if cookie, err := r.Cookie("Token"); err == nil {
		userID, _ = GetUserIDFromToken(cookie.Value)
	}

	opts := QueryOptions{
		UserID: userID,
		Filter: filter,
	}

	query, args := BuildPostQuery(opts)
	posts, err := GetPosts(userID, query, args...)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func IsPostLikedByUser(postID, user_id int) bool {
	var existingInteraction int
	db.QueryRow("SELECT interaction FROM post_interactions WHERE user_id = ? AND post_id = ?", user_id, postID).Scan(&existingInteraction)
	return existingInteraction == 1
}
