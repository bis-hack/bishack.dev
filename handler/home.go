package handler

import (
	"net/http"

	"bishack.dev/services/post"
	"bishack.dev/utils"
	"bishack.dev/utils/session"
	"github.com/gorilla/context"
)

// Home ...
func Home(w http.ResponseWriter, r *http.Request) {
	sess := context.Get(r, "session").(interface {
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	// get user details from context and cast it as map[string]string if
	// not nil
	u := context.Get(r, "user")

	ps := context.Get(r, "postService").(interface {
		GetPosts() []*post.Post
	})

	utils.Render(w, "main", "home", map[string]interface{}{
		"Title": "Bisdak Tech Community",
		"Flash": sess.GetFlash(w, r),
		"User":  u,
		"Posts": ps.GetPosts(),
	})
}

// NotFound ...
func NotFound(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "error", "notfound", map[string]interface{}{
		"Title":       "404 - Not Found",
		"Description": "The page you're looking for could not be found",
	})
}
