package forum

import (
	"database/sql"
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

type data struct {
	Post     POST
	COMMENT  []COMMENT
	ErrorMsg string
}

func Commenting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	c, _ := r.Cookie("Token")
	uname := TokenMap[c.Value][1]
	content := r.FormValue("Content")

	id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := getPost(id)

	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
		Comments, err := GetComment(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		DATA := data{
			Post:     post,
			COMMENT:  Comments,
			ErrorMsg: "Comment cannot be empty. Please enter some content.",
		}
		err = CommentT.Execute(w, DATA)
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
		}
		return
	}

	err = insertComment(id, uname, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Comments, err := GetComment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	DATA := data{Post: post, COMMENT: Comments}
	err = CommentT.Execute(w, DATA)
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
	}
}

func Comment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	post, err := getPost(id)

	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	Comments, err := GetComment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Data := data{Post: post, COMMENT: Comments}
	err = CommentT.Execute(w, Data)
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
	}
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
