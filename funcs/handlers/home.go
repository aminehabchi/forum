package forum

import (
	"database/sql"
	data "forum/funcs/database"
	Error "forum/funcs/error"
	types "forum/funcs/types"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	if r.URL.Path != "/" {
		Error.ErrorHandler(w, http.StatusNotFound)
		return
	}

	userID := 0
	if c, err := r.Cookie("Token"); err == nil {
		userID, _ = data.GetUserIDFromToken(c.Value)
	}

	opts := types.QueryOptions{
		UserID: userID,
		Limit:  4,
		Offset: 0,
	}

	query, args := data.BuildPostQuery(opts)

	posts, err := data.GetPosts(userID, query, args...)
	if err != nil && err != sql.ErrNoRows {
		Error.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts      []types.POST
		IsLoggedIn bool
		Categories []string
	}{
		Posts:      posts,
		IsLoggedIn: userID > 0,
		Categories: DefaultCategories,
	}

	if err = types.HomeT.Execute(w, data); err != nil {
		Error.ErrorHandler(w, http.StatusInternalServerError)
	}
}
