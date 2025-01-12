package forum

import (
	"encoding/base64"
	"fmt"

	Data "forum/funcs/types"
	"io"
	"os"
	"strings"
	"time"
)

func GetPosts(userID int, query string, args ...interface{}) ([]Data.POST, error) {
	rows, err := Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Data.POST
	for rows.Next() {
		var p Data.POST
		var timeCreated time.Time
		err := rows.Scan(&p.ID, &p.USER_ID, &p.Title, &timeCreated, &p.Content, &p.ImgBase64, &p.Name)
		if err != nil {
			return nil, err
		}

		// Get comment count
		err = Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", p.ID).Scan(&p.NbComment)
		if err != nil {
			return nil, err
		}

		p.Likes = getPostLikeDisLike(p.ID, 1)
		p.Dislikes = getPostLikeDisLike(p.ID, -1)
		p.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")

		if userID > 0 {
			Db.QueryRow("SELECT interaction FROM post_interactions WHERE user_id = ? AND post_id = ?", userID, p.ID).Scan(&p.UserInteraction)
		}

		categories := getPostCategories(p.ID)
		p.Category = categories
		p.ImgBase64, _ = EncodeImg("./images/" + p.ImgBase64)
		posts = append(posts, p)

	}
	return posts, nil
}

func BuildPostQuery(opts Data.QueryOptions) (string, []interface{}) {
	baseQuery := `
        SELECT posts.id, posts.user_id, posts.title, posts.created_at, posts.content,posts.img, users.uname 
        FROM posts
        JOIN users ON posts.user_id = users.id`

	var args []interface{}

	switch opts.Filter {
	case "":
		if opts.PostID != "" {
			baseQuery += " WHERE posts.id = ?"
			args = append(args, opts.PostID)
		}

	case "created":
		baseQuery += " WHERE posts.user_id = ?"
		args = append(args, opts.UserID)

	case "liked":
		baseQuery += `
            JOIN post_interactions ON post_interactions.post_id = posts.id
            WHERE post_interactions.user_id = ? AND post_interactions.interaction = 1`
		args = append(args, opts.UserID)
	default:
		formattedCategory := strings.ToUpper(string(opts.Filter[0])) + strings.ToLower(opts.Filter[1:])
		baseQuery += `
            JOIN post_categories ON post_categories.post_id = posts.id
            WHERE post_categories.category = ?`
		args = append(args, formattedCategory)
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
	err := Db.QueryRow("SELECT COUNT(*) FROM post_interactions WHERE post_id=? AND interaction=?", post_id, inter).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func getPostCategories(postID int) []string {
	rows, _ := Db.Query("SELECT category FROM post_categories WHERE post_categories.post_id=?", postID)

	var categories []string
	for rows.Next() {
		var category string
		rows.Scan(&category)
		categories = append(categories, category)
	}

	return categories
}

func InsertPost(id int, title, content string, categories []string, imgName string) error {
	selector := `INSERT INTO posts(title,content,user_id,img) VALUES (?,?,?,?)`
	a, err := Db.Exec(selector, title, content, id, imgName)
	if err != nil {
		return err
	}
	idPost, _ := a.LastInsertId()

	for _, category := range categories {

		formattedCategory := strings.ToUpper(string(category[0])) + strings.ToLower(string(category[1:]))
		selector = `INSERT INTO post_categories(post_id,category) VALUES (?,?)`
		_, _ = Db.Exec(selector, idPost, formattedCategory)
	}
	return nil
}

func EncodeImg(imagePath string) (string, error) {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// Read the image file into a byte slice
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// Check if the image file is empty
	if len(imageBytes) == 0 {
		return "", fmt.Errorf("image file is empty")
	}

	// Encode the image bytes to base64
	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)
	return encodedImage, nil
}
