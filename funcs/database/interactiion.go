package forum

func AddInteractions(user_id int, postID, action, types string) error {
	interaction := 0
	var err error
	if types == "post" {
		err = Db.QueryRow("SELECT interaction FROM post_interactions where post_id= ? and user_id= ?", postID, user_id).Scan(&interaction)
	} else {
		err = Db.QueryRow("SELECT interaction FROM comment_interactions where comment_id= ? and user_id= ?", postID, user_id).Scan(&interaction)
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
			_, err = Db.Exec("UPDATE post_interactions SET interaction=? where post_id= ? and user_id= ?", interaction, postID, user_id)
		} else {
			_, err = Db.Exec("UPDATE comment_interactions SET interaction=? where comment_id= ? and user_id= ?", interaction, postID, user_id)
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
			_, err := Db.Exec(selector, user_id, postID, 1)
			if err != nil {
				return err
			}
		} else {
			_, err := Db.Exec(selector, user_id, postID, -1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
