package main

import (
	"fmt"
	"net/http"

	forum "forum/funcs"
)

func main() {
	err:=forum.Createbase()
	if err!=nil{
		fmt.Println(err)
		return 
	}
	
	fs := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("GET /", forum.Home)

	http.HandleFunc("GET /register", forum.Register)
	http.HandleFunc("POST /registerIngo", forum.RegisterIngo)

	http.HandleFunc("GET /login", forum.Login)
	http.HandleFunc("POST /loginInfo", forum.LoginInfo)

	http.HandleFunc("GET /Posting", forum.Auth(forum.Posting))
	http.HandleFunc("POST /PostInfo", forum.Auth(forum.PostInfo))

	http.HandleFunc("GET /logout", forum.Logout)

	http.HandleFunc("GET /Comment", forum.Comment)
	http.HandleFunc("POST /Commenting", forum.Auth(forum.Commenting))

	http.HandleFunc("POST /like-dislike", forum.Auth(forum.HandleLikeDislike))
	http.HandleFunc("GET /filter", forum.FilterHandler)
	fmt.Println("http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}
