package forum

import (
	"database/sql"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, 500)
		return
	}
	if r.URL.Path != "/" {
		ErrorHandler(w, 400)
		return
	}

	userID := 0
	if c, err := r.Cookie("Token"); err == nil {
		userID, _ = GetUserIDFromToken(c.Value)
	}

	opts := QueryOptions{
		UserID: userID,
		Limit:  3,
		Offset: 0,
	}

	query, args := BuildPostQuery(opts)

	posts, err := GetPosts(userID, query, args...)
	if err != nil && err != sql.ErrNoRows {
		ErrorHandler(w, 500)
		return
	}

	data := struct {
		Posts      []POST
		IsLoggedIn bool
		Categories []string
	}{
		Posts:      posts,
		IsLoggedIn: userID > 0,
		Categories: defaultCategories,
	}

	if err = HomeT.Execute(w, data); err != nil {
		ErrorHandler(w, 500)
	}
}
