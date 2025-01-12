package main

import (
	"fmt"
	"net/http"
	forum "forum/funcs"
	data "forum/funcs/database"
	handlers "forum/funcs/handlers"
	types "forum/funcs/types"
)

func main() {
	err := types.ParseFiles()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = data.CreateDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/static/", forum.StaticFileHandler)

	http.HandleFunc("/", handlers.Home)

	http.HandleFunc("/register", handlers.AuthLG(handlers.Register))
	http.HandleFunc("/login", handlers.AuthLG(handlers.Login))
	http.HandleFunc("/logout", handlers.Auth(handlers.Logout))

	http.HandleFunc("/Posting", handlers.Auth(handlers.Posting))
	http.HandleFunc("/load-more-posts", handlers.LoadMorePosts)

	http.HandleFunc("/Commenting", handlers.Commenting)
	http.HandleFunc("/load-more-comments", handlers.LoadMoreComments)

	http.HandleFunc("/like-dislike", handlers.HandleLikeDislike)
	http.HandleFunc("/filter", handlers.FilterHandler)

	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
