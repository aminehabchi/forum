package forum

import (
	"database/sql"
	"encoding/json"
	data "forum/funcs/database"
	Error "forum/funcs/error"
	types "forum/funcs/types"
	"net/http"
	"strconv"
	"strings"
)

func Commenting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		post_id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil || post_id <= 0 {
			Error.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		user_id := 0
		if Cookie, err := r.Cookie("Token"); err == nil {
			user_id, _ = data.GetUserIDFromToken(Cookie.Value)
		}

		opts := types.QueryOptions{
			UserID: user_id,
			PostID: r.FormValue("post_id"),
		}

		query, args := data.BuildPostQuery(opts)

		posts, err := data.GetPosts(user_id, query, args...)
		if err == sql.ErrNoRows || len(posts) == 0 {
			Error.ErrorHandler(w, http.StatusNotFound)
			return
		}
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		comments, err := data.GetComment(post_id, user_id, 3, 0)
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		Data := types.Data{
			Post:    posts[0],
			COMMENT: comments,
		}
		if err = types.CommentT.Execute(w, Data); err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		c, err := r.Cookie("Token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user_id, err := data.GetUserIDFromToken(c.Value)
		if err != nil {
			ClearSession(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		content := strings.TrimSpace(r.FormValue("Content"))
		if content == "" {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "Please enter a comment",
			})

			if err != nil {
				Error.ErrorHandler(w, http.StatusInternalServerError)
			}
			return
		}

		post_id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil || post_id <= 0 {
			Error.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		var comment_id int
		comment_id, err = data.InsertComment(post_id, user_id, content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		type comment struct {
			Uname   string
			Content string
			Id      int
		}

		var Comment comment
		err = data.Db.QueryRow(`SELECT uname FROM users WHERE id = ?`, user_id).Scan(&Comment.Uname)
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		Comment.Content = content
		Comment.Id = comment_id
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(Comment)
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
		}
	default:
		Error.ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}
