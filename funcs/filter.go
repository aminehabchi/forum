package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("category")
	created := r.FormValue("created")
	liked := r.FormValue("liked")
	filteredPosts := filterPosts(category, created, liked, r)

	_, err := r.Cookie("username")
	isLoggedIn := err == nil

	temp, _ := template.ParseFiles("templates/home.html")
	data := struct {
		Posts      []POST
		IsLoggedIn bool
		Categories []string
	}{
		Posts:      filteredPosts,
		IsLoggedIn: isLoggedIn,
		Categories: []string{"General", "Technology", "News", "Entertainment", "Hobbies", "Lifestyle"},
	}
	temp.Execute(w, data)
}

func filterPosts(category, created, liked string, r *http.Request) []POST {
	var filteredPosts []POST
	posts, e := GetPosts()
	if e != nil {
		fmt.Println(e)
	}
	user, _ := r.Cookie("username")
	for _, post := range posts {
		if category != "" && !elementExists(post.Category, category) {
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

	return filteredPosts
}

func elementExists(arr []string, elem string) bool {
    for _, v := range arr {
        if v == elem {
            return true
        }
    }
    return false
}