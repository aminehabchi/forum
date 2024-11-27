package forum

import (
	"time"
)

type POST struct {
	ID              int
	USER_ID         int
	Name            string
	Title           string
	CreatedAt       string
	Content         string
	Category        []string
	Likes           int
	Dislikes        int
	NbComment       int
	UserInteraction int
}

// SELECT posts.id, posts.user_id,posts.title,posts.created_at ,posts.content, posts.category,users.uname FROM posts JOIN users ON posts.user_id = users.id ORDER BY posts.id DESC LIMIT ? OFFSET ?", limit, offset
func Get_Posts(userID int, Query string) ([]POST, error) {
	rows, err := db.Query(Query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []POST
	for rows.Next() {
		var p POST

		var timeCreated time.Time
		err := rows.Scan(&p.ID, &p.USER_ID, &p.Title, &timeCreated, &p.Content, &p.Name)
		if err != nil {
			return nil, err
		}
		err = db.QueryRow("SELECT COUNT(*) FROM posts JOIN comments ON posts.id = comments.post_id WHERE posts.id = ? ", p.ID).Scan(&p.NbComment)

		if err != nil {
			return nil, err
		}
		p.Likes = getPostLikeDisLike(p.ID, 1)
		p.Dislikes = getPostLikeDisLike(p.ID, -1)
		p.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")

		db.QueryRow("SELECT interaction FROM post_interactions WHERE user_id = ? AND post_id = ?", userID, p.ID).Scan(&p.UserInteraction)

		posts = append(posts, p)
	}
	return posts, nil
}

func getPostLikeDisLike(post_id, inter int) int {
	// Use a count query to directly get the number of interactions
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM post_interactions WHERE post_id=? AND interaction=?", post_id, inter).Scan(&count)
	if err != nil {
		// Optionally log the error
		return 0
	}
	return count
}

func getCommentLikeDisLike(comment_id, inter int) int {
	// Use a count query to directly get the number of interactions
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM comment_interactions WHERE comment_id=? AND interaction=?", comment_id, inter).Scan(&count)
	if err != nil {
		// Optionally log the error
		return 0
	}
	return count
}
