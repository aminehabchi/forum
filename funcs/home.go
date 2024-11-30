package forum

import (
	"database/sql"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "page not found", 404)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}
