package forum

import (
	"fmt"
	"net/http"
)

func IsPostLikedByUser(postID int, name string) bool {
	var existingInteraction int
	db.QueryRow("SELECT interaction FROM interactions WHERE username = ? AND post_id = ?", name, postID).Scan(&existingInteraction)
	return existingInteraction == 1
}

func HandleLikeDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	c, _ := r.Cookie("username")

	username := c.Value
	postID := r.FormValue("post_id")
	action := r.FormValue("action")
	types := r.FormValue("type")
	commentid := r.FormValue("commentid")
	if (types != "post" && types != "comment") || (action != "dislike" && action != "like") {
		// badrequest
	}
	addInteractions(username, commentid, action, types)

	if types == "post" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		link := fmt.Sprintf("/Comment?post_id=%v", postID)
		http.Redirect(w, r, link, http.StatusSeeOther)
	}
}

func addInteractions(username, postID, action, types string) {
	interaction := 0
	err := db.QueryRow("SELECT interaction FROM interactions where type = ? and post_id= ? and username= ?", types, postID, username).Scan(&interaction)
	if err == nil {
		if interaction == 1 && action == "like" {
			interaction = 0
		} else if interaction != 1 && action == "like" {
			interaction = 1
		} else if interaction != -1 && action == "dislike" {
			interaction = -1
		} else if interaction == -1 && action == "dislike" {
			interaction = 0
		}
		_, err = db.Exec("UPDATE interactions SET interaction=? where type = ? and post_id= ? and username= ?", interaction, types, postID, username)
		fmt.Println(err)
	} else {
		selector := `INSERT INTO interactions(username,post_id,type,interaction) VALUES (?,?,?,?)`
		if action == "like" {
			_, err := db.Exec(selector, username, postID, types, 1)
			fmt.Println(err)
		} else {
			_, err := db.Exec(selector, username, postID, types, -1)
			fmt.Println(err)
		}
	}
}
