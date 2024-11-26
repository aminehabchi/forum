package forum

import (
	"database/sql"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.FormValue("category")
	created := r.FormValue("created")
	liked := r.FormValue("liked")

	filteredPosts, err := filterPosts(category, created, liked, r)
	if err == sql.ErrNoRows {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		return
	}
	_, err = r.Cookie("Token")
	isLoggedIn := err == nil

	data := struct {
		Posts      []POST
		IsLoggedIn bool
		Categories []string
	}{
		Posts:      filteredPosts,
		IsLoggedIn: isLoggedIn,
		Categories: []string{"General", "Technology", "News", "Entertainment", "Hobbies", "Lifestyle"},
	}
	err = HomeT.Execute(w, data)
	if err != nil {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		return
	}
}

func filterPosts(category, created, liked string, r *http.Request) ([]POST, error) {
	var filteredPosts []POST

	user, err := r.Cookie("Token")
	var user_id int
	if err == nil {
		user_id, _ = GetUserNameFromToken(user.Value)
	}

	posts, e := GetPosts(user_id)
	if e != nil {
		return []POST{}, e
	}
	for _, post := range posts {
		if category != "" && !ElementExists(post.Category, category) {
			continue
		}

		if created == "on" {
			if post.USER_ID != user_id {
				continue
			}
		}

		if liked == "on" {
			if !IsPostLikedByUser(post.ID, user_id) {
				continue
			}
		}

		filteredPosts = append(filteredPosts, post)
	}

	return filteredPosts, nil
}

func ElementExists(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func IsPostLikedByUser(postID, user_id int) bool {
	var existingInteraction int
	db.QueryRow("SELECT interaction FROM post_interactions WHERE user_id = ? AND post_id = ?", user_id, postID).Scan(&existingInteraction)
	return existingInteraction == 1
}
