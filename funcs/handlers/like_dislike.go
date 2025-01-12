package forum

import (
	data "forum/funcs/database"
	Error "forum/funcs/error"
	"net/http"
)

func HandleLikeDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	user_id, ok := CheckIfCookieValid(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	action := r.FormValue("action")
	types := r.FormValue("type")
	commentid := r.FormValue("commentid")

	if (types != "post" && types != "comment") || (action != "dislike" && action != "like") {
		Error.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	err := data.AddInteractions(user_id, commentid, action, types)
	if err != nil {
		Error.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}
