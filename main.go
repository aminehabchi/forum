package main

import (
	"fmt"
	"net/http"

	forum "forum/funcs"
)

func main() {
	err := forum.ParseFiles()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = forum.CreateDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/static/", forum.StaticFileHandler)

	http.HandleFunc("/", forum.Home)

	http.HandleFunc("/register", forum.AuthLG(forum.Register))
	http.HandleFunc("/login", forum.AuthLG(forum.Login))
	http.HandleFunc("/logout", forum.Auth(forum.Logout))

	http.HandleFunc("/Posting", forum.Auth(forum.Posting))
	http.HandleFunc("/load-more-posts", forum.LoadMorePosts)

	http.HandleFunc("/Commenting", forum.Commenting)
	http.HandleFunc("/load-more-comments", forum.LoadMoreComments)

	http.HandleFunc("/like-dislike", forum.Auth(forum.HandleLikeDislike))
	http.HandleFunc("/filter", forum.FilterHandler)

	http.HandleFunc("/error", forum.ErrorHandler)

	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
