package forum

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func Posting(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/Posting.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	// categories := []string{"General", "Technology", "Sports", "Entertainment"}

	temp.Execute(w, nil)
}

func PostInfo(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("username")

	// r.PostForm("categories")

	temp, err := template.ParseFiles("templates/Posting.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	uname := c.Value
	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.Form["categories"]

	if title == "" || content == "" || len(category) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		temp.Execute(w, "All fields are required, Please fill them")
		return
	}

	if !CategoryFilter(category) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("xxxxx")
		temp.Execute(w, "Invalid categorie, Please write valid catgerie")
		return
	}

	err = insertPost(uname, title, content, strings.Join(category, " "))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
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
	for _, v := range categories {
		if !ElementExists(v) {
			return false
		}
	}
	return true
}
