package handler

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"bishack.dev/services/post"
	"bishack.dev/services/user"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserPosts(t *testing.T) {
	t.Run("user not found", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "session", s)

		u.On("GetUser", "").Return(nil)

		GetUserPosts(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ok", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)
		p := new(postMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "postService", p)
		context.Set(r, "session", s)

		u.On("GetUser", "").Return(&user.User{})
		p.On("GetUserPosts", "").Return([]*post.Post{
			&post.Post{
				Title: "The quick brown test",
			},
		})
		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		GetUserPosts(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("The quick brown test"), w.Body.String())
	})
}
