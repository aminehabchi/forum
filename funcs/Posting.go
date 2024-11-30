package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PostingData struct {
	Categories []string
	Error      string
}

var defaultCategories = []string{"General", "News", "Entertainment", "Hobbies", "Lifestyle", "Technology"}

func Posting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := PostingT.Execute(w, PostingData{Categories: defaultCategories})
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
		}
	case http.MethodPost:
		var err error
		c, _ := r.Cookie("Token")
		id, _ := GetUserIDFromToken(c.Value)

		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		category := r.Form["categories"]

		if title == "" || content == "" || len(category) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			data := PostingData{
				Categories: defaultCategories,
				Error:      "All fields are required, Please fill them",
			}
			err = PostingT.Execute(w, data)
			if err != nil {
				http.Error(w, "Could not load template", http.StatusInternalServerError)
			}
			return
		}

		if !CategoryFilter(category) {
			w.WriteHeader(http.StatusBadRequest)
			data := PostingData{
				Categories: defaultCategories,
				Error:      "Invalid categorie, Please write valid catgerie",
			}
			err = PostingT.Execute(w, data)
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
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

		formattedCategory := strings.ToUpper(string(category[0])) + strings.ToLower(string(category[1:]))
		selector = `INSERT INTO post_categories(post_id,category) VALUES (?,?)`
		_, _ = db.Exec(selector, idPost, formattedCategory)
	}
	return nil
}

func CategoryFilter(categories []string) bool {
	for _, v := range categories {
		if !allCategories[strings.ToLower(v)] {
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

	offsetValue := r.FormValue("offset")
	filterType := r.FormValue("type")

	fmt.Println("filter ->", filterType)

	offset, err := strconv.Atoi(offsetValue)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := r.Cookie("Token")

	var user_id int
	if err == nil {
		user_id, _ = GetUserIDFromToken(user.Value)
	}

	opts := QueryOptions{
		UserID: user_id,
		Limit:  3,
		Offset: offset,
		Filter: filterType,
	}

	query, args := BuildPostQuery(opts)

	posts, err := GetPosts(user_id, query, args...)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
