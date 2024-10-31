package forum

import "html/template"

var (
	CommentT  *template.Template
	HomeT     *template.Template
	LoginT    *template.Template
	PostingT  *template.Template
	RegisterT *template.Template
)

func ParseFiles() error {
	var err error
	RegisterT, err = template.ParseFiles("templates/register.html")
	if err != nil {
		return err
	}
	HomeT, err = template.ParseFiles("templates/home.html")
	if err != nil {
		return err
	}
	LoginT, err = template.ParseFiles("templates/login.html")
	if err != nil {
		return err
	}
	PostingT, err = template.ParseFiles("templates/Posting.html")
	if err != nil {
		return err
	}
	CommentT, err = template.ParseFiles("templates/comment.html")
	if err != nil {
		return err
	}
	return nil
}
