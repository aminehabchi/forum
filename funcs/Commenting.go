package forum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type COMMENT struct {
	Id              int
	USER_ID         int
	Uname           string
	Content         string
	Likes           int
	Dislikes        int
	UserInteraction int
}

type data struct {
	Post     POST
	COMMENT  []COMMENT
	ErrorMsg string
}

func Commenting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		post_id, _ := strconv.Atoi(r.FormValue("post_id"))
		if post_id == 0 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		var user_id int
		Cookie, err := r.Cookie("Token")
		isLoggedIn := err == nil
		if isLoggedIn {
			user_id, _ = GetUserNameFromToken(Cookie.Value)
		}

		Posts, err := Get_Posts(user_id, `
		SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content, users.uname FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.id=`+r.FormValue("post_id"))
		if err == sql.ErrNoRows {
			http.Error(w, "bad request1", http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		var post POST
		post = Posts[0]
		Comments, err := GetComment(post_id, user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Data := data{Post: post, COMMENT: Comments}
		err = CommentT.Execute(w, Data)
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		c, err := r.Cookie("Token")
		if err != nil {
			w.WriteHeader(403)
			return
		}
		user_id, err := GetUserNameFromToken(c.Value)
		if err != nil {
			w.WriteHeader(403)
			return
		}

		post_id, _ := strconv.Atoi(r.FormValue("post_id"))
		content := r.FormValue("Content")
		if content == "" || post_id == 0 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		var comment_id int
		comment_id, err = insertComment(post_id, user_id, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type comment struct {
			Uname   string `json:"uname"`
			Content string `json:"content"`
			Id      int    `json:"id"`
		}

		var Comment comment
		db.QueryRow(`SELECT uname FROM users WHERE id = ?`, user_id).Scan(&Comment.Uname)
		Comment.Content = content
		Comment.Id = comment_id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Comment)
	default:
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
	}
}

func insertComment(postid, user_id int, content string) (int, error) {
	selector := `INSERT INTO comments(post_id,user_id,content) VALUES (?,?,?)`
	result, err := db.Exec(selector, postid, user_id, content)
	if err != nil {
		return -1, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(commentID), nil
}

func GetComment(id, userID int) ([]COMMENT, error) {
	rows, err := db.Query("SELECT comments.id,users.uname,comments.content FROM comments JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? ORDER BY comments.id DESC", id)
	if err != nil {
		return []COMMENT{}, err
	}
	defer rows.Close()

	var comments []COMMENT
	for rows.Next() {
		var content, uname string
		var comment_id int
		err := rows.Scan(&comment_id, &uname, &content)
		if err != nil {
			log.Fatal(err, "err2")
			return []COMMENT{}, err
		}
		comment := COMMENT{Id: comment_id, Uname: uname, Content: content}
		comment.Likes = getCommentLikeDisLike(comment_id, 1)
		comment.Dislikes = getCommentLikeDisLike(comment_id, -1)
		db.QueryRow("SELECT interaction FROM comment_interactions WHERE user_id = ? AND comment_id = ?", userID, comment_id).Scan(&comment.UserInteraction)
		comments = append(comments, comment)
	}

	// Check for any errors during the iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err, "err1")
		return []COMMENT{}, err
	}
	return comments, nil
}
