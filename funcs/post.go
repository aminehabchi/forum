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
	Like, err := db.Query("SELECT interaction FROM interactions where post_id=? and type= ? and interaction=?", post_id, types, inter)
	if err != nil {
		return 0
	}
	i := 0
	for Like.Next() {
		i++
	}
	return i
}
