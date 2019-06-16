package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"bishack.dev/services/post"
	"bishack.dev/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
)

// New ...
func New(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user = user.(map[string]string)
	utils.Render(w, "main", "new-form", map[string]interface{}{
		"Title":          "Create New Post",
		"User":           user,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// CreatePost ...
func CreatePost(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	attr := map[string]interface{}{}

	// cast to int
	publish, _ := strconv.Atoi(r.PostForm.Get("publish"))
	attr["publish"] = publish

	attr["title"] = r.PostForm.Get("title")
	attr["userId"] = r.PostForm.Get("userId")
	attr["content"] = r.PostForm.Get("content")
	attr["userPic"] = r.PostForm.Get("userPic")
	attr["username"] = r.PostForm.Get("username")

	ps := context.Get(r, "postService").(interface {
		Create(params map[string]interface{}) *post.Post
	})

	p := ps.Create(attr)
	if p == nil {
		log.Println("error")
		http.Redirect(w, r, "/new", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/p/%s", p.ID), http.StatusSeeOther)
}

// GetPost ...
func GetPost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	ps := context.Get(r, "postService").(interface {
		Get(id string) *post.Post
	})

	post := ps.Get(id)
	if post == nil {
		// not found
		utils.Render(w, "error", "notfound", map[string]interface{}{
			"Title": "Not Found",
		})

		return
	}

	user := context.Get(r, "user")

	utils.Render(w, "main", "post", map[string]interface{}{
		"Title": post.Title,
		"Post":  post,
		"User":  user,
	})
}
