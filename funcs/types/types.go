package forum

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
	ImgBase64       string
}

type QueryOptions struct {
	UserID int
	PostID string
	Filter string
	Limit  int
	Offset int
}
type COMMENT struct {
	Id              int
	USER_ID         int
	Uname           string
	Content         string
	Likes           int
	Dislikes        int
	UserInteraction int
}

type Data struct {
	Post    POST
	COMMENT []COMMENT
}

type User struct {
	ID       int
	Password string
}
