package forum

import (
	"encoding/json"
	"net/http"
)

func HandleLikeDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
		return
	}
	c, _ := r.Cookie("Token")
	user_id, _ := GetUserIDFromToken(c.Value)
	action := r.FormValue("action")
	types := r.FormValue("type")
	commentid := r.FormValue("commentid")

	if (types != "post" && types != "comment") || (action != "dislike" && action != "like") {
		http.Error(w, "invalid Type or Action", http.StatusBadRequest)
		return
	}
	err := addInteractions(user_id, commentid, action, types)
	if err != nil {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		return
	}

	type LikeDislike struct {
		Like        int `json:"like"`
		Dislike     int `json:"dislike"`
		Interaction int `json:"interaction"`
	}
	var res LikeDislike
	if types == "post" {
		db.QueryRow("SELECT COUNT(*) FROM post_interactions WHERE post_id= ? AND interaction = ?", commentid, 1).Scan(&res.Like)
		db.QueryRow("SELECT COUNT(*) FROM post_interactions where post_id= ? and interaction = ?", commentid, -1).Scan(&res.Dislike)
		db.QueryRow("SELECT interaction FROM post_interactions WHERE user_id = ? AND post_id = ?", user_id, commentid).Scan(&res.Interaction)

	} else {
		db.QueryRow("SELECT COUNT(*) FROM comment_interactions WHERE comment_id= ? AND interaction = ?", commentid, 1).Scan(&res.Like)
		db.QueryRow("SELECT COUNT(*) FROM comment_interactions where comment_id= ? and interaction = ?", commentid, -1).Scan(&res.Dislike)
		db.QueryRow("SELECT interaction FROM comment_interactions WHERE user_id = ? AND comment_id = ?", user_id, commentid).Scan(&res.Interaction)

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
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
