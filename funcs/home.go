package forum

import (
	"database/sql"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "page not found", 404)
		return
	}

	c, err := r.Cookie("Token")
	isLoggedIn := err == nil
	var userID int
	if isLoggedIn {
		userID, _ = GetUserNameFromToken(c.Value)
	}
	posts, err := Get_Posts(userID, "SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content, users.uname FROM posts JOIN users ON posts.user_id = users.id ORDER BY posts.id DESC LIMIT 3 OFFSET 0")
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
		return
	}
	data := struct {
		Posts      []POST
		IsLoggedIn bool
		Categories []string
	}{
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
		Categories: []string{"General", "Technology", "News", "Entertainment", "Hobbies", "Lifestyle"},
	}

	err = HomeT.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}
