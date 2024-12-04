package forum

import (
	"net/http"
)

func HandleLikeDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	user_id, ok := CheckIfCookieValid(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	action := r.FormValue("action")
	types := r.FormValue("type")
	commentid := r.FormValue("commentid")

	if (types != "post" && types != "comment") || (action != "dislike" && action != "like") {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	err := addInteractions(user_id, commentid, action, types)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

func addInteractions(user_id int, postID, action, types string) error {
	interaction := 0
	var err error
	if types == "post" {
		err = db.QueryRow("SELECT interaction FROM post_interactions where post_id= ? and user_id= ?", postID, user_id).Scan(&interaction)
	} else {
		err = db.QueryRow("SELECT interaction FROM comment_interactions where comment_id= ? and user_id= ?", postID, user_id).Scan(&interaction)
	}

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
		if types == "post" {
			_, err = db.Exec("UPDATE post_interactions SET interaction=? where post_id= ? and user_id= ?", interaction, postID, user_id)
		} else {
			_, err = db.Exec("UPDATE comment_interactions SET interaction=? where comment_id= ? and user_id= ?", interaction, postID, user_id)
		}
		if err != nil {
			return err
		}
	} else {
		var selector string
		if types == "post" {
			selector = `INSERT INTO post_interactions(user_id,post_id,interaction) VALUES (?,?,?)`
		} else {
			selector = `INSERT INTO comment_interactions(user_id,comment_id,interaction) VALUES (?,?,?)`
		}
		if action == "like" {
			_, err := db.Exec(selector, user_id, postID, 1)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(selector, user_id, postID, -1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
