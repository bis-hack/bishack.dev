package handler

import (
	"fmt"
	"net/http"

	"bishack.dev/services/post"
	"bishack.dev/services/user"
	"bishack.dev/utils"
	"bishack.dev/utils/session"
	"github.com/gorilla/context"
)

// GetUserPosts ...
func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get(":username")

	sess := context.Get(r, "session").(interface {
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	us := context.Get(r, "userService").(interface {
		GetUser(username string) *user.User
	})

	user := us.GetUser(username)
	if user == nil {
		utils.Render(w, "error", "notfound", map[string]interface{}{
			"Title":       "404 - Not Found",
			"Description": "The page you're looking for could not be found",
		})
		return
	}

	ps := context.Get(r, "postService").(interface {
		GetUserPosts(username string) []*post.Post
	})

	title := fmt.Sprintf("%s's Posts", user.Name)
	utils.Render(w, "main", "user-page", map[string]interface{}{
		"Title": title,
		"Flash": sess.GetFlash(w, r),
		"Posts": ps.GetUserPosts(username),
		"User":  user,
	})
}
