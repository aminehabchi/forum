package forum

import (
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {

}

func IsPostLikedByUser(postID, user_id int) bool {
	var existingInteraction int
	db.QueryRow("SELECT interaction FROM post_interactions WHERE user_id = ? AND post_id = ?", user_id, postID).Scan(&existingInteraction)
	return existingInteraction == 1
}
