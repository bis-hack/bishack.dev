package handler

import (
	"net/http"

	"bishack.dev/utils/session"
	"github.com/gorilla/context"
)

// Home ...
func Home(w http.ResponseWriter, r *http.Request) {
	var user map[string]string
	sess := context.Get(r, "session").(interface {
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	user = sessionUser(r)

	render(w, "main", "index", map[string]interface{}{
		"Title": "Home",
		"Flash": sess.GetFlash(w, r),
		"User":  user,
	})
}

// NotFound ...
func NotFound(w http.ResponseWriter, r *http.Request) {
	render(w, "main", "notfound", map[string]interface{}{
		"Title": "Not Found",
	})
}
