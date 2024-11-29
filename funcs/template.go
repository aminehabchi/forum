package forum

import (
	"fmt"
	"html/template"
)

var (
	CommentT  *template.Template
	HomeT     *template.Template
	LoginT    *template.Template
	PostingT  *template.Template
	RegisterT *template.Template

	templatesPaths = map[string]string{
		"register": "templates/register.html",
		"home":     "templates/home.html",
		"login":    "templates/login.html",
		"posting":  "templates/Posting.html",
		"comment":  "templates/comment.html",
	}
)

func ParseFiles() error {
	for name, path := range templatesPaths {
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %v", name, err)
		}

		switch name {
		case "register":
			RegisterT = tmpl
		case "home":
			HomeT = tmpl
		case "login":
			LoginT = tmpl
		case "posting":
			PostingT = tmpl
		case "comment":
			CommentT = tmpl
		}
	}

	return nil
}
