package forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type COMMENT struct {
	Id       int
	Uname    string
	Content  string
	Postid   int
	Likes    int
	Dislikes int
}

func Commenting(w http.ResponseWriter, r *http.Request) {
	a := r.FormValue("post_id")
	id, _ := strconv.Atoi(r.FormValue("post_id"))
	content := r.FormValue("Content")
	c, _ := r.Cookie("username")
	uname := c.Value
	fmt.Println(id, a, content, c)
	post, _ := getPost(id)

	err := insertComment(id, uname, content)
	if err != nil {
		fmt.Println(err)
	}

	// post, _ := getPost(id)
	Comments, err := GetComment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp, _ := template.ParseFiles("templates/comment.html")
	type data struct {
		Post    POST
		COMMENT []COMMENT
	}
	ff := data{Post: post, COMMENT: Comments}
	temp.Execute(w, ff)
}

func Comment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post, err := getPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Comments, err := GetComment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp, err := template.ParseFiles("templates/comment.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	type data struct {
		Post    POST
		COMMENT []COMMENT
	}
	fmt.Println(post.ID)
	ff := data{Post: post, COMMENT: Comments}
	temp.Execute(w, ff)
}

func insertComment(postid int, uname, content string) error {
	selector := `INSERT INTO comments(post_id,uname,content) VALUES (?,?,?)`
	_, err := db.Exec(selector, postid, uname, content)
	if err != nil {
		return err
	}
	return nil
}

func GetComment(id int) ([]COMMENT, error) {
	rows, err := db.Query("SELECT id,post_id, uname, content FROM comments WHERE post_id = ?", id)
	if err != nil {
		log.Fatal(err, "99999")
		return []COMMENT{}, err
	}
	defer rows.Close()

	var comments []COMMENT
	for rows.Next() {
		var content, uname string
		var id, pid int
		err := rows.Scan(&id, &pid, &uname, &content)
		if err != nil {
			log.Fatal(err, "err2")
			return []COMMENT{}, err
		}
		comment := COMMENT{Id: id, Postid: pid, Uname: uname, Content: content}
		comment.Likes = getLikeDisLike("comment", id, 1)

		comment.Dislikes = getLikeDisLike("comment", id, -1)
		comments = append(comments, comment)
	}

	// Check for any errors during the iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err, "err1")
		return []COMMENT{}, err
	}
	return comments, nil
}

func getPost(id int) (POST, error) {
	query := `SELECT id,uname, title, content, category FROM posts WHERE id = ?`
	var post POST
	var str string
	err := db.QueryRow(query, id).Scan(&post.ID, &post.Name, &post.Title, &post.Content, &str)
	if err != nil {
		return POST{}, err
	}
	post.Category = strings.Split(str, " ")
	return post, nil
}
