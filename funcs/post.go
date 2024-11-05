package forum

import "strings"

func GetPosts() ([]POST, error) {
	rows, err := db.Query("SELECT id, uname, title, content, category FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []POST
	for rows.Next() {
		var p POST
		var str_categories string
		err := rows.Scan(&p.ID, &p.Name, &p.Title, &p.Content, &str_categories)
		if err != nil {
			return nil, err
		}
		p.Category = strings.Split(str_categories, " ")
		p.Likes = getLikeDisLike("post", p.ID, 1)
		p.Dislikes = getLikeDisLike("post", p.ID, -1)
		posts = append(posts, p)
	}
	return posts, nil
}

func getLikeDisLike(types string, post_id, inter int) int {
	// Use a count query to directly get the number of interactions
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM interactions WHERE post_id=? AND type=? AND interaction=?", post_id, types, inter).Scan(&count)
	if err != nil {
		// Optionally log the error
		return 0
	}
	return count
}
