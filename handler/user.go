package handler

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"

	"bishack.dev/services/like"
	"bishack.dev/services/post"
	"bishack.dev/services/user"
	"bishack.dev/utils"
	"bishack.dev/utils/session"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"

	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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

	posts := ps.GetUserPosts(username)

	ls := context.Get(r, "likeService").(interface {
		GetLikes(id string) ([]*like.Like, error)
	})
	// populate likes/comments count
	var wg sync.WaitGroup
	wg.Add(len(posts))
	for _, p := range posts {
		go func(p *post.Post) {
			defer wg.Done()

			results, err := ls.GetLikes(p.ID)
			if err != nil {
				log.Println("GetLikes error", err.Error())
				return
			}

			p.LikesCount = int64(len(results))
		}(p)
	}
	wg.Wait()

	title := fmt.Sprintf("%s's Posts", user.Name)
	utils.Render(w, "main", "user-page", map[string]interface{}{
		"Title":  title,
		"Flash":  sess.GetFlash(w, r),
		"Posts":  posts,
		"Author": user,
		"User":   context.Get(r, "user"),
	})
}

// Profile ...
func Profile(w http.ResponseWriter, r *http.Request) {

	sess := context.Get(r, "session").(interface {
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	// get user details from context
	user := context.Get(r, "user")

	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.Render(w, "main", "profile-form", map[string]interface{}{
		"Title":          "Edit User Profile",
		"Flash":          sess.GetFlash(w, r),
		"User":           user,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// UpdateProfile ...
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	sess := context.Get(r, "session").(interface {
		SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
		GetUser(r *http.Request) map[string]string
	})

	su := sess.GetUser(r)
	if su == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_ = r.ParseForm()

	args := map[string]string{}
	args["name"] = r.FormValue("name")
	args["locale"] = r.FormValue("locale")
	profile := r.FormValue("profile")
	args["profile"] = profile[:int(math.Min(128, float64(len(profile))))] // Maximum of 128 chars
	args["website"] = r.FormValue("website")

	us := context.Get(r, "userService").(interface {
		UpdateUser(token string, attrs map[string]string) (*cip.UpdateUserAttributesOutput, error)
	})

	token := su["token"]
	if _, err := us.UpdateUser(token, args); err != nil {
		sess.SetFlash(w, r, "error", err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sess.SetFlash(w, r, "success", "Profile Updated")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
