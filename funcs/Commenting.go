package forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type COMMENT struct {
	Id       int
	USER_ID  int
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
	user_id, _ := GetUserNameFromToken(c.Value)
	content := r.FormValue("Content")

	post_id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
		link := fmt.Sprintf("/Comment?post_id=%v", post_id)
		http.Redirect(w, r, link, http.StatusSeeOther)
		return
	}

	err = insertComment(post_id, user_id, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	link := fmt.Sprintf("/Comment?post_id=%v", post_id)
	http.Redirect(w, r, link, http.StatusSeeOther)
}

func Comment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	post_id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	post, err := getPost(post_id)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println(post_id)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	Comments, err := GetComment(post_id)
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

func insertComment(postid, user_id int, content string) error {
	selector := `INSERT INTO comments(post_id,user_id,content) VALUES (?,?,?)`
	_, err := db.Exec(selector, postid, user_id, content)
	if err != nil {
		return err
	}
	return nil
}

func GetComment(id int) ([]COMMENT, error) {
	rows, err := db.Query("SELECT comments.id,comments.post_id, users.uname, comments.content FROM comments JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? ORDER BY comments.id DESC", id)
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
		comment.Likes = getCommentLikeDisLike(id, 1)
		comment.Dislikes = getCommentLikeDisLike(id, -1)
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
	query := `SELECT posts.id, posts.title, posts.content, posts.category,users.uname FROM posts JOIN users ON posts.user_id = users.id WHERE posts.id = ?`
	var post POST
	var str string
	err := db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &str, &post.Name)
	if err != nil {
		fmt.Println(err)
		return POST{}, err
	}
	post.Category = strings.Split(str, " ")
	return post, nil
}
