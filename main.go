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
	err = forum.Createbase()
	if err != nil {
		fmt.Println(err)
		return
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", forum.Home)

	http.HandleFunc("/register", forum.Register)
	http.HandleFunc("/registerIngo", forum.RegisterIngo)

	http.HandleFunc("/login", forum.Login)
	http.HandleFunc("/loginInfo", forum.LoginInfo)

	http.HandleFunc("/Posting", forum.Auth(forum.Posting))
	http.HandleFunc("/PostInfo", forum.Auth(forum.PostInfo))

	http.HandleFunc("/logout", forum.Logout)

	http.HandleFunc("/Comment", forum.Comment)
	http.HandleFunc("/Commenting", forum.Auth(forum.Commenting))

	http.HandleFunc("/like-dislike", forum.Auth(forum.HandleLikeDislike))
	http.HandleFunc("/filter", forum.FilterHandler)
	fmt.Println("http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}
