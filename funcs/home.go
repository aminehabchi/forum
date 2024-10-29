package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/home.html")
	posts, e := GetPosts()
	if e != nil {
		fmt.Println(e)
	}

	_, err := r.Cookie("username")
	isLoggedIn := err == nil

	data := struct {
		Posts      []POST
		IsLoggedIn bool
		Categories []string
	}{
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
		Categories: []string{"General", "Technology", "News", "Entertainment", "Hobbies", "Lifestyle"},
	}

	temp.Execute(w, data)
}
