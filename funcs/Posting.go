package forum

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func PostInfo(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("username")
	
	// r.PostForm("categories")

	uname := c.Value
	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.Form["categories"]
	
	err := insertPost(uname, title, content, strings.Join(category, " "))
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
func Posting(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/Posting.html")
	// categories := []string{"General", "Technology", "Sports", "Entertainment"}
	temp.Execute(w, nil)
}

func insertPost(uname, title, content ,category string) error {
	selector := `INSERT INTO posts(uname,title,content,category) VALUES (?,?,?,?)`
	_, err := db.Exec(selector, uname, title, content,category)
	if err != nil {
		return err
	}
	return nil
}
