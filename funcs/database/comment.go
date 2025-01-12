package forum

import (
	Data "forum/funcs/types"
	"log"
)

func InsertComment(postid, user_id int, content string) (int, error) {
	selector := `INSERT INTO comments(post_id,user_id,content) VALUES (?,?,?)`
	result, err := Db.Exec(selector, postid, user_id, content)
	if err != nil {
		return -1, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(commentID), nil
}

func GetComment(id, userID, limit, offset int) ([]Data.COMMENT, error) {
	rows, err := Db.Query("SELECT comments.id,users.uname,comments.content FROM comments JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? ORDER BY comments.id DESC LIMIT ? OFFSET ?", id, limit, offset)
	if err != nil {
		return []Data.COMMENT{}, err
	}
	defer rows.Close()

	var comments []Data.COMMENT
	for rows.Next() {
		var content, uname string
		var comment_id int
		err := rows.Scan(&comment_id, &uname, &content)
		if err != nil {
			return []Data.COMMENT{}, err
		}
		comment := Data.COMMENT{Id: comment_id, Uname: uname, Content: content}
		comment.Likes = getCommentLikeDisLike(comment_id, 1)
		comment.Dislikes = getCommentLikeDisLike(comment_id, -1)
		Db.QueryRow("SELECT interaction FROM comment_interactions WHERE user_id = ? AND comment_id = ?", userID, comment_id).Scan(&comment.UserInteraction)
		comments = append(comments, comment)
	}

	// Check for any errors during the iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err, "err1")
		return []Data.COMMENT{}, err
	}
	return comments, nil
}

func getCommentLikeDisLike(comment_id, inter int) int {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM comment_interactions WHERE comment_id=? AND interaction=?", comment_id, inter).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
