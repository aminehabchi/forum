package forum

import (
	
	data "forum/funcs/database"
	Error "forum/funcs/error"
	types "forum/funcs/types"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type PostingData struct {
	Categories []string
	Error      string
}

var DefaultCategories = []string{"General", "News", "Entertainment", "Hobbies", "Lifestyle", "Technology"}

func Posting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := types.PostingT.Execute(w, PostingData{Categories: DefaultCategories})
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
		}
	case http.MethodPost:
		var err error
		c, _ := r.Cookie("Token")
		id, _ := data.GetUserIDFromToken(c.Value)

		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		category := r.Form["categories"]

		file, header, err := r.FormFile("file")
		imageExists := true

		if err != nil {
			if err.Error() == "http: no such file" {
				imageExists = false
			} else {
				http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
				return
			}
		}
		var name string
		if imageExists {
			defer file.Close()
			name, err = data.GenereteTocken()
			if err != nil {
				Error.ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			extensions := strings.Split(header.Filename, ".")[1]
			name = name + "." + extensions
			err = saveImg(file, "./images/"+name)
			if err != nil {
				http.Error(w, "Unable to save the image", http.StatusInternalServerError)
				return
			}
		}

		if title == "" || (content == "" && !imageExists) || len(category) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			data := PostingData{
				Categories: DefaultCategories,
				Error:      "All fields are required, Please fill them",
			}
			err = types.PostingT.Execute(w, data)
			if err != nil {
				Error.ErrorHandler(w, http.StatusInternalServerError)
			}
			return
		}

		if !CategoryFilter(category) {
			w.WriteHeader(http.StatusBadRequest)
			data := PostingData{
				Categories: DefaultCategories,
				Error:      "Invalid categorie, Please write valid categorie",
			}
			err = types.PostingT.Execute(w, data)
			if err != nil {
				Error.ErrorHandler(w, http.StatusInternalServerError)
			}
			return
		}

		err = data.InsertPost(id, title, content, category, name)
		if err != nil {
			Error.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		Error.ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}
func saveImg(file multipart.File, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, file)
	if err != nil {
		return err
	}

	return nil
}

func CategoryFilter(categories []string) bool {
	for _, v := range categories {
		if !data.AllCategories[strings.ToLower(v)] {
			return false
		}
	}
	return true
}
