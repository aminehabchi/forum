package forum

import (
	"database/sql"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.FormValue("category")
	created := r.FormValue("created")
	liked := r.FormValue("liked")

	filteredPosts, err := filterPosts(category, created, liked, r)
	if err == sql.ErrNoRows {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		return
	}
	_, err = r.Cookie("username")
	isLoggedIn := err == nil

	data := struct {
		Posts      []POST
		IsLoggedIn bool
		Categories []string
	}{
		Posts:      filteredPosts,
		IsLoggedIn: isLoggedIn,
		Categories: []string{"General", "Technology", "News", "Entertainment", "Hobbies", "Lifestyle"},
	}
	err = HomeT.Execute(w, data)
	if err != nil {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		return
	}
}

func filterPosts(category, created, liked string, r *http.Request) ([]POST, error) {
	var filteredPosts []POST
	posts, e := GetPosts()
	if e != nil {
		return []POST{}, e
	}
	user, _ := r.Cookie("username")
	for _, post := range posts {
		if category != "" && !ElementExists(post.Category, category) {
			continue
		}

		if created == "on" {
			if post.Name != user.Value {
				continue
			}
		}

		if liked == "on" {
			if !IsPostLikedByUser(post.ID, user.Value) {
				continue
			}
		}

		filteredPosts = append(filteredPosts, post)
	}

	return filteredPosts, nil
}

func ElementExists(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}
