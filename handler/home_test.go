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

func TestHome(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		s := new(sessionMock)
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "postService", p)
		context.Set(r, "likeService", l)
		context.Set(r, "session", s)

		p.On("GetPosts").Return(nil)
		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Home(w, r)

		assert.Regexp(t, regexp.MustCompile("Log In"), w.Body.String())
		s.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		s := new(sessionMock)
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "user", &user.User{
			Username: "tibur",
		})

		context.Set(r, "postService", p)
		context.Set(r, "likeService", l)
		context.Set(r, "session", s)

		p.On("GetPosts").Return(nil)
		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Home(w, r)

		assert.Regexp(t, regexp.MustCompile("@tibur"), w.Body.String())

		s.AssertExpectations(t)
	})

	t.Run("likes with error", func(t *testing.T) {
		s := new(sessionMock)
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "user", &user.User{
			Username: "tibur",
		})

		context.Set(r, "postService", p)
		context.Set(r, "likeService", l)
		context.Set(r, "session", s)

		p.On("GetPosts").Return([]*post.Post{
			{ID: "test"},
		})
		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)
		l.On("GetLikes", "test").Return(nil, errors.New(""))

		Home(w, r)

		assert.Regexp(t, regexp.MustCompile("@tibur"), w.Body.String())

		s.AssertExpectations(t)
	})

	t.Run("likes ok", func(t *testing.T) {
		s := new(sessionMock)
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "user", &user.User{
			Username: "tibur",
		})

		context.Set(r, "postService", p)
		context.Set(r, "likeService", l)
		context.Set(r, "session", s)

		p.On("GetPosts").Return([]*post.Post{
			{ID: "test"},
		})
		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)
		l.On("GetLikes", "test").Return([]*like.Like{{}}, nil)

		Home(w, r)

		assert.Regexp(t, regexp.MustCompile("@tibur"), w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/xxx", nil)

	NotFound(w, r)

	assert.Regexp(t, regexp.MustCompile("Not Found"), w.Body.String())
}
