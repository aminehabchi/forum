package forum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Posting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	err := PostingT.Execute(w, nil)
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
}

func PostInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
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

	err = insertPost(id, title, content, strings.Join(category, " "))
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func insertPost(id int, title, content, category string) error {
	selector := `INSERT INTO posts(title,content,category,user_id) VALUES (?,?,?,?)`
	_, err := db.Exec(selector, title, content, category, id)
	if err != nil {
		return err
	}
	return nil
}

func CategoryFilter(categories []string) bool {
	allCategories := []string{"General", "Technology", "News", "Entertainment", "Hobbies", "Lifestyle"}
	for _, v := range categories {
		if !ElementExists(allCategories, v) {
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

	offset, _ := strconv.Atoi(r.FormValue("offset"))
	category := r.FormValue("category")
	created := r.FormValue("created")
	liked := r.FormValue("liked")

	// fmt.Println(offset, "c",category, "cr",created, "l",liked)

	filteredPosts, err := filterPosts(category, created, liked, r, offset, 3)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredPosts)
}
