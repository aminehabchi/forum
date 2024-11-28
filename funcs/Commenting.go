package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		post_id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		Cookie, _ := r.Cookie("Token")
		user_id, _ := GetUserNameFromToken(Cookie.Value)
		post, err := Get_Posts(user_id, `
		SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content, users.uname FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.user_id=`+strconv.Itoa(user_id)+` and posts.id=`+r.FormValue("post_id"))

		if err != nil && err != sql.ErrNoRows {
			fmt.Println(post_id)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		c, err := r.Cookie("Token")
		isLoggedIn := err == nil
		var userID int
		if isLoggedIn {
			userID, _ = GetUserNameFromToken(c.Value)
		}

		Comments, err := GetComment(post_id, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Data := data{Post: post[0], COMMENT: Comments}
		err = CommentT.Execute(w, Data)
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
		}
	case http.MethodPost:
		c, _ := r.Cookie("Token")
		user_id, _ := GetUserNameFromToken(c.Value)
		content := r.FormValue("Content")

		post_id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if content == "" {
			link := "/Commenting?post_id=" + strconv.Itoa(post_id)
			http.Redirect(w, r, link, http.StatusSeeOther)
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

		var res comment
		db.QueryRow(`SELECT uname FROM users WHERE id = ?`, user_id).Scan(&res.Uname)
		res.Content = content
		res.Id = comment_id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
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
