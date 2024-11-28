package forum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

var allCategories map[string]bool

func Posting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := PostingT.Execute(w, nil)
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
		}
	case http.MethodPost:
		var err error
		c, _ := r.Cookie("Token")

		id, _ := GetUserNameFromToken(c.Value)
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.Form["categories"]

		if title == "" || content == "" || len(category) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			err = PostingT.Execute(w, "All fields are required, Please fill them")
			if err != nil {
				http.Error(w, "Could not load template", http.StatusInternalServerError)
			}
			return
		}

		if !CategoryFilter(category) {
			w.WriteHeader(http.StatusBadRequest)
			err = PostingT.Execute(w, "Invalid categorie, Please write valid catgerie")
			if err != nil {
				http.Error(w, "Could not load template", http.StatusInternalServerError)
			}
			return
		}

		err = insertPost(id, title, content, category)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	default:
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
	}
}

func insertPost(id int, title, content string, categories []string) error {
	selector := `INSERT INTO posts(title,content,user_id) VALUES (?,?,?)`
	a, err := db.Exec(selector, title, content, id)
	if err != nil {
		return err
	}
	idPost, _ := a.LastInsertId()

	for _, category := range categories {
		selector = `INSERT INTO post_categories(post_id,category) VALUES (?,?)`
		_, _ = db.Exec(selector, idPost, category)
	}
	return nil
}

func CategoryFilter(categories []string) bool {
	for _, v := range categories {
		if !allCategories[v] {
			return false
		}
	}
	return true
}

func LoadMorePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	offset := r.FormValue("offset")

	user, err := r.Cookie("Token")

	var user_id int
	if err == nil {
		user_id, _ = GetUserNameFromToken(user.Value)
	}

	quert := "SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content,users.uname FROM posts JOIN users ON posts.user_id = users.id ORDER BY posts.id DESC LIMIT 3 OFFSET " + offset

	posts, err := Get_Posts(user_id, quert)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
