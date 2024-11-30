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

type QueryOptions struct {
	UserID int
	PostID string
	Filter string
	Limit  int
	Offset int
}

func GetPosts(userID int, query string, args ...interface{}) ([]POST, error) {
	rows, err := db.Query(query, args...)
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

		// Get comment count
		err = db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", p.ID).Scan(&p.NbComment)
		if err != nil {
			return nil, err
		}

		p.Likes = getPostLikeDisLike(p.ID, 1)
		p.Dislikes = getPostLikeDisLike(p.ID, -1)
		p.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")

		if userID > 0 {
			db.QueryRow("SELECT interaction FROM post_interactions WHERE user_id = ? AND post_id = ?", userID, p.ID).Scan(&p.UserInteraction)
		}

		// Get categories
		categories := getPostCategories(p.ID)
		p.Category = categories

		posts = append(posts, p)
	}
	return posts, nil
}

func BuildPostQuery(opts QueryOptions) (string, []interface{}) {
	baseQuery := `
        SELECT posts.id, posts.user_id, posts.title, posts.created_at, posts.content, users.uname 
        FROM posts
        JOIN users ON posts.user_id = users.id`

	var args []interface{}

	switch opts.Filter {
	case "":
		if opts.PostID != "" {
			baseQuery += " WHERE posts.id = ?"
			args = append(args, opts.PostID)
		}

	case "Created":
		baseQuery += " WHERE posts.user_id = ?"
		args = append(args, opts.UserID)

	case "Liked":
		baseQuery += `
            JOIN post_interactions ON post_interactions.post_id = posts.id
            WHERE post_interactions.user_id = ? AND post_interactions.interaction = 1`
		args = append(args, opts.UserID)
	default:
		baseQuery += `
            JOIN post_categories ON post_categories.post_id = posts.id
            WHERE post_categories.category = ?`
		args = append(args, opts.Filter)
	}

	baseQuery += " ORDER BY posts.id DESC"

	if opts.Limit > 0 {
		baseQuery += " LIMIT ? OFFSET ?"
		args = append(args, opts.Limit, opts.Offset)
	}

	return baseQuery, args
}

func getPostLikeDisLike(post_id, inter int) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM post_interactions WHERE post_id=? AND interaction=?", post_id, inter).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func getPostCategories(postID int) []string {
	rows, _ := db.Query("SELECT category FROM post_categories WHERE post_categories.post_id=?", postID)

	var categories []string
	for rows.Next() {
		var category string
		rows.Scan(&category)
		categories = append(categories, category)
	}

	return categories
}
