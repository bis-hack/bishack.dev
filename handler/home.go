package handler

import (
	"net/http"

	"bishack.dev/utils/session"
)

// Home ...
func Home(w http.ResponseWriter, r *http.Request) {
	var user map[string]string

	creds := session.GetUser(r)
	if creds != nil {
		user = sessionUser(creds["token"])
	}

	render(w, "main", "index", map[string]interface{}{
		"Title": "Home",
		"Flash": session.GetFlash(w, r),
		"User":  user,
	})
}

// NotFound ...
func NotFound(w http.ResponseWriter, r *http.Request) {
	render(w, "main", "notfound", map[string]interface{}{
		"Title": "Not Found",
	})
}
