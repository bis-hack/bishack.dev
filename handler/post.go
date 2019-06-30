package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"bishack.dev/services/like"
	"bishack.dev/services/post"
	"bishack.dev/services/user"
	"bishack.dev/utils"
	"bishack.dev/utils/session"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
)

// New ...
func New(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user")
	if u == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.Render(w, "main", "new-form", map[string]interface{}{
		"Title":          "Create New Post",
		"User":           u,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	id := r.FormValue("id")
	cover := r.FormValue("cover")
	created, _ := strconv.Atoi(r.FormValue("created"))
	content := r.FormValue("content")

	ps := context.Get(r, "postService").(interface {
		UpdatePost(string, string, string, int64) error
	})

	sess := context.Get(r, "session").(interface {
		SetFlash(http.ResponseWriter, *http.Request, string, string)
	})

	err := ps.UpdatePost(id, cover, content, int64(created))
	if err != nil {
		sess.SetFlash(w, r, "error", "An error occurred. Try again.")
	} else {
		sess.SetFlash(w, r, "success", "Changes saved successfully!")
	}

	http.Redirect(w, r, "/edit/"+id, http.StatusSeeOther)
}

// EditPost ...
func EditPost(w http.ResponseWriter, r *http.Request) {
	uc := context.Get(r, "user")
	if uc == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u := uc.(*user.User)

	ps := context.Get(r, "postService").(interface {
		GetPost(username, id string) *post.Post
	})

	username := u.Username
	id := r.URL.Query().Get(":id")

	post := ps.GetPost(username, id)
	if post == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sess := context.Get(r, "session").(interface {
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	flash := sess.GetFlash(w, r)

	utils.Render(w, "main", "edit-form", map[string]interface{}{
		"Title":          post.Title,
		"User":           u,
		"Post":           post,
		"Flash":          flash,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// CreatePost ...
func CreatePost(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	attr := map[string]interface{}{}

	// cast to int
	publish, _ := strconv.Atoi(r.PostForm.Get("publish"))
	content := r.PostForm.Get("content")
	attr["publish"] = publish

	attr["title"] = r.PostForm.Get("title")
	attr["cover"] = r.PostForm.Get("cover")
	attr["author"] = r.PostForm.Get("author")
	attr["content"] = content
	attr["userPic"] = r.PostForm.Get("userPic")
	attr["username"] = r.PostForm.Get("username")
	attr["readingTime"] = computeReadingTime(content)

	ps := context.Get(r, "postService").(interface {
		CreatePost(params map[string]interface{}) *post.Post
	})

	p := ps.CreatePost(attr)
	if p == nil {
		log.Println("error")
		http.Redirect(w, r, "/new", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", p.Username, p.ID), http.StatusSeeOther)
}

// GetPost ...
func GetPost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	username := r.URL.Query().Get(":username")

	ps := context.Get(r, "postService").(interface {
		GetPost(username, id string) *post.Post
	})

	post := ps.GetPost(username, id)
	if post == nil {
		// not found
		utils.Render(w, "error", "notfound", map[string]interface{}{
			"Title": "Not Found",
		})

		return
	}

	post.ReadingTime = computeReadingTime(post.Content)

	ls := context.Get(r, "likeService").(interface {
		GetLike(id, username string) (*like.Like, error)
		GetLikes(id string) ([]*like.Like, error)
	})

	var u *user.User
	liker := false
	uc := context.Get(r, "user")
	if uc != nil {
		u = uc.(*user.User)

		_, err := ls.GetLike(post.ID, u.Username)
		if err == nil {
			liker = true
		}

	}

	likes, err := ls.GetLikes(post.ID)
	if err == nil {
		post.LikesCount = int64(len(likes))
	}

	chunks := strings.Split(post.Content, "\r\n\r\n")
	description := chunks[0]
	if len(chunks) >= 2 {
		description = chunks[1]
		description = strings.Join(
			regexp.MustCompile(`[a-zA-Z0-9\.\-\s\/\:,;]+`).FindAllString(description, -1),
			" ",
		)
	}

	utils.Render(w, "main", "post", map[string]interface{}{
		"Title":          post.Title,
		"Post":           post,
		"Description":    description,
		"User":           u,
		"Cover":          post.Cover,
		"Liker":          liker,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// ToggleLike ...
func ToggleLike(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	username := ""

	uc := context.Get(r, "user")
	if uc != nil {
		u := uc.(*user.User)
		username = u.Username
	}

	ls := context.Get(r, "likeService").(interface {
		ToggleLike(id, username string) error
	})

	err := ls.ToggleLike(id, username)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "ok")
}

func computeReadingTime(content string) int {
	const avgWPM = 265 // 265 wpm
	wordCount := len(content)

	return wordCount / avgWPM
}
