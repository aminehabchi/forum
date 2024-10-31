package forum

import (
	"net/http"
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
	c, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	uname := c.Value
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

	err = insertPost(uname, title, content, strings.Join(category, " "))
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func insertPost(uname, title, content, category string) error {
	selector := `INSERT INTO posts(uname,title,content,category) VALUES (?,?,?,?)`
	_, err := db.Exec(selector, uname, title, content, category)
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
