package handler

import (
	"log"
	"net/http"
	"sync"

	"bishack.dev/services/like"
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

	posts := ps.GetPosts()

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

	utils.Render(w, "main", "home", map[string]interface{}{
		"Title": "Bisdak Tech Community",
		"Flash": sess.GetFlash(w, r),
		"User":  u,
		"Posts": posts,
	})
}

// NotFound ...
func NotFound(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "error", "notfound", map[string]interface{}{
		"Title":       "404 - Not Found",
		"Description": "The page you're looking for could not be found",
	})
}
