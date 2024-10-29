package forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Commenting(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("post_id"))
	content := r.FormValue("Content")
	c, _ := r.Cookie("username")
	uname := c.Value
	post, _ := getPost(id)
	
	err := insertComment(post.Title, uname, content)
	if err != nil {
		fmt.Println(err)
	}

	// post, _ := getPost(id)
	Comments := GetComment(id)
	temp, _ := template.ParseFiles("templates/comment.html")
	type data struct {
		Post    POST
		COMMENT []COMMENT
	}
	ff := data{Post: post, COMMENT: Comments}
	temp.Execute(w, ff)
}

func Comment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("post_id"))
	post, _ := getPost(id)
	Comments := GetComment(id)
	temp, _ := template.ParseFiles("templates/comment.html")
	type data struct {
		Post    POST
		COMMENT []COMMENT
	}
	ff := data{Post: post, COMMENT: Comments}
	temp.Execute(w, ff)
}

func insertComment(title, uname, content string) error {
	selector := `INSERT INTO comments(postTitle,uname,content) VALUES (?,?,?)`
	_, err := db.Exec(selector, title, uname, content)
	if err != nil {
		return err
	}
	return nil
}

type COMMENT struct {
	Id       int
	Uname    string
	Content  string
	Likes    int
	Dislikes int
}

func GetComment(id int) []COMMENT {
	rows, err := db.Query("SELECT id, uname, content FROM comments WHERE id = ?", id)
	if err != nil {
		log.Fatal(err, "  99999")
	}
	defer rows.Close()

	var comments []COMMENT
	for rows.Next() {
		var content, uname string
		var id int
		err := rows.Scan(&id, &uname, &content)
		if err != nil {
			log.Fatal(err, "err2")
		}
		comment := COMMENT{Id: id, Uname: uname, Content: content}
		comment.Likes = getLikeDisLike("comment", id, 1)
		comment.Dislikes = getLikeDisLike("comment", id, -1)
		comments = append(comments, comment)
	}

	// Check for any errors during the iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err, "err1")
	}
	return comments
}

func getPost(id int) (POST, error) {
	query := `SELECT uname, title, content, category FROM posts WHERE id = ?`
	var post POST
	var str string
	err := db.QueryRow(query, id).Scan(&post.Name,&post.Title ,&post.Content, &str )
	if err != nil {
		return POST{}, err
	}
	post.Category = strings.Split(str, " ")
	return post, nil
}
