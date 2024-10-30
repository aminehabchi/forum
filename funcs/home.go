package forum

import (
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	posts, err := GetPosts()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting posts:", err)
		return
	}

	_, err = r.Cookie("username")
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

	err = temp.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}
