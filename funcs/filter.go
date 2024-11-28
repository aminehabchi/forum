package forum

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.FormValue("type")
	fmt.Println(filter)
	if !allCategories[filter] {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	Cookie, _ := r.Cookie("Token")
	user_id, _ := GetUserNameFromToken(Cookie.Value)

	var posts []POST
	var query string
	var err error
	switch filter {
	case "":
		query = `
		SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content, users.uname FROM posts
		JOIN users ON posts.user_id = users.id 
		ORDER BY posts.id DESC 
		`
		posts, err = Get_Posts(user_id, query)
	case "Created":
		query = `
			SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content, users.uname FROM posts
			JOIN users ON posts.user_id = users.id
			WHERE  posts.user_id=` + strconv.Itoa(user_id) + `
			ORDER BY posts.id DESC 
			`
		posts, err = Get_Posts(user_id, query)
	case "Liked":
		query = `
			SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content, users.uname FROM posts
			JOIN users ON posts.user_id = users.id
			JOIN post_interactions ON post_interactions.post_id = posts.id
			WHERE  users.id=` + strconv.Itoa(user_id) + ` AND post_interactions.interaction=1
			ORDER BY posts.id DESC;
			`
		posts, err = Get_Posts(user_id, query)

	default:
		query = `
		SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content, users.uname FROM posts
		JOIN users ON posts.user_id = users.id
		JOIN post_categories ON post_categories.post_id = posts.id
		WHERE  post_categories.category='` + filter + `'
		ORDER BY posts.id DESC 
		`
		posts, err = Get_Posts(user_id, query)
	}
	fmt.Println(err)
	for _, v := range posts {
		fmt.Println(v)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func IsPostLikedByUser(postID, user_id int) bool {
	var existingInteraction int
	db.QueryRow("SELECT interaction FROM post_interactions WHERE user_id = ? AND post_id = ?", user_id, postID).Scan(&existingInteraction)
	return existingInteraction == 1
}
