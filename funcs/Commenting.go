package forum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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
}

func Commenting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		post_id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil || post_id <= 0 {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		user_id := 0
		if Cookie, err := r.Cookie("Token"); err == nil {
			user_id, _ = GetUserIDFromToken(Cookie.Value)
		}

		opts := QueryOptions{
			UserID: user_id,
			PostID: r.FormValue("post_id"),
		}

		query, args := BuildPostQuery(opts)

		posts, err := GetPosts(user_id, query, args...)
		if err == sql.ErrNoRows || len(posts) == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		comments, err := GetComment(post_id, user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Data := data{
			Post:    posts[0],
			COMMENT: comments,
		}
		if err = CommentT.Execute(w, Data); err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		c, err := r.Cookie("Token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user_id, err := GetUserIDFromToken(c.Value)
		if err != nil {
			ClearSession(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		content := strings.TrimSpace(r.FormValue("Content"))
		if content == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Please enter a comment",
			})
			return
		}

		post_id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil || post_id <= 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var comment_id int
		comment_id, err = insertComment(post_id, user_id, content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		type comment struct {
			Uname   string `json:"uname"`
			Content string `json:"content"`
			Id      int    `json:"id"`
		}

		var Comment comment
		err = db.QueryRow(`SELECT uname FROM users WHERE id = ?`, user_id).Scan(&Comment.Uname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

func getCommentLikeDisLike(comment_id, inter int) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM comment_interactions WHERE comment_id=? AND interaction=?", comment_id, inter).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
