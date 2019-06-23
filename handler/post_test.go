package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"bishack.dev/services/like"
	"bishack.dev/services/post"
	"bishack.dev/services/user"
	_ "bishack.dev/testing"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Run("not logged in", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/new", nil)

		New(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("logged in", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/new", nil)

		user := &user.User{
			Username: "test",
		}

		context.Set(r, "user", user)

		New(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile(`value="test"`), w.Body.String())
	})
}

func TestCreatePost(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		p := new(postMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/new", nil)

		context.Set(r, "postService", p)
		p.On("CreatePost", mock.MatchedBy(func(vals map[string]interface{}) bool {
			return true
		})).Return(nil)

		CreatePost(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		p.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		p := new(postMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/new", nil)

		context.Set(r, "postService", p)
		p.On("CreatePost", mock.MatchedBy(func(vals map[string]interface{}) bool {
			return true
		})).Return(&post.Post{
			Title:   "test",
			Content: "test",
			ID:      "test",
		})

		CreatePost(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		p.AssertExpectations(t)
	})
}

func TestGetPost(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "postService", p)
		context.Set(r, "likeService", l)

		p.On("GetPost", mock.MatchedBy(func(id string) bool {
			return true
		})).Return(nil)

		GetPost(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		p.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "postService", p)
		context.Set(r, "likeService", l)

		p.On("GetPost", mock.MatchedBy(func(id string) bool {
			return true
		})).Return(&post.Post{
			Title: "test",
			Content: `
			On January 21, JYP Entertainment announced they would be debuting
			a new girl group, being the first girl group from the label since
			Twice’s debut in 2015.[6][7] On the same day, the group’s official
			YouTube account was created and the label’s official channel shared
			a video trailer unveiling the five members.[8][9]
			`,
		})
		l.On("GetLikes", "").Return(nil, errors.New(""))

		GetPost(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("test"), w.Body.String())
	})

	t.Run("ok with likes", func(t *testing.T) {
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "postService", p)
		context.Set(r, "user", &user.User{Username: "test"})
		context.Set(r, "likeService", l)

		p.On("GetPost", mock.MatchedBy(func(id string) bool {
			return true
		})).Return(&post.Post{
			Title: "test",
			ID:    "test",
		})
		l.On("GetLike", "test", "test").Return(&like.Like{}, nil)
		l.On("GetLikes", "test").Return([]*like.Like{
			{},
		}, nil)

		GetPost(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("test"), w.Body.String())
	})
}

func TestToggleLike(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "likeService", l)

		l.On("ToggleLike", "", "").Return(errors.New(""))

		ToggleLike(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Regexp(t, regexp.MustCompile("error"), w.Body.String())

	})

	t.Run("ok", func(t *testing.T) {
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "user", &user.User{
			Username: "test",
		})

		context.Set(r, "likeService", l)

		l.On("ToggleLike", "", "test").Return(nil)

		ToggleLike(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("ok"), w.Body.String())

	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		p := new(postMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "postService", p)
		context.Set(r, "session", s)

		p.On("UpdatePost", "", "", "", int64(0)).Return(errors.New(""))
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "An error occurred. Try again.").Return(nil)

		UpdatePost(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("ok", func(t *testing.T) {
		p := new(postMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "postService", p)
		context.Set(r, "session", s)

		p.On("UpdatePost", "", "", "", int64(0)).Return(nil)
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "success", "Changes saved successfully!").Return(nil)

		UpdatePost(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})
}

func TestEditPost(t *testing.T) {
	t.Run("user not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		EditPost(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("post is nil", func(t *testing.T) {
		p := new(postMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/edit/test", nil)

		context.Set(r, "postService", p)
		context.Set(r, "user", &user.User{})

		p.On("GetPost", mock.MatchedBy(func(username string) bool {
			return true
		}), mock.MatchedBy(func(id string) bool {
			return true
		})).Return(nil)

		EditPost(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("ok", func(t *testing.T) {
		p := new(postMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/edit/test", nil)

		context.Set(r, "postService", p)
		context.Set(r, "session", s)
		context.Set(r, "user", &user.User{})

		p.On("GetPost", mock.MatchedBy(func(username string) bool {
			return true
		}), mock.MatchedBy(func(id string) bool {
			return true
		})).Return(&post.Post{})

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		EditPost(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("edit-form"), w.Body.String())
	})
}
