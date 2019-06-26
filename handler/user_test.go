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
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserPosts(t *testing.T) {
	t.Run("user not found", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "likeService", l)
		context.Set(r, "session", s)

		u.On("GetUser", "").Return(nil)
		l.On("GetLikes", "").Return(nil, errors.New(""))

		GetUserPosts(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("likes with error", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "postService", p)
		context.Set(r, "likeService", l)
		context.Set(r, "session", s)

		u.On("GetUser", "").Return(&user.User{})
		p.On("GetUserPosts", "").Return([]*post.Post{
			{
				Title: "The quick brown test",
			},
		})
		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)
		l.On("GetLikes", "").Return(nil, errors.New(""))

		GetUserPosts(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("The quick brown test"), w.Body.String())
	})

	t.Run("likes ok", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)
		p := new(postMock)
		l := new(likeMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "postService", p)
		context.Set(r, "likeService", l)
		context.Set(r, "session", s)

		u.On("GetUser", "").Return(&user.User{})
		p.On("GetUserPosts", "").Return([]*post.Post{
			{
				Title: "The quick brown test",
			},
		})
		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)
		l.On("GetLikes", "").Return([]*like.Like{{}}, nil)

		GetUserPosts(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("The quick brown test"), w.Body.String())
	})
}

func TestProfile(t *testing.T) {
	t.Run("User nil", func(t *testing.T) {

		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/profile", nil)

		context.Set(r, "session", s)

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Profile(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("OK", func(t *testing.T) {

		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/profile", nil)

		context.Set(r, "session", s)
		context.Set(r, "user", map[string]string{})

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Profile(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Run("user nil", func(t *testing.T) {
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		context.Set(r, "session", s)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		UpdateProfile(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"token": "test",
		})
		u.On("UpdateUser", "test", mock.MatchedBy(func(args map[string]string) bool {
			return true
		})).Return(nil, errors.New("Invalid Token"))
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Invalid Token")

		UpdateProfile(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("ok", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"token": "test",
		})
		u.On("UpdateUser", "test", mock.MatchedBy(func(args map[string]string) bool {
			return true
		})).Return(&cip.UpdateUserAttributesOutput{}, nil)
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "success", "Profile Updated")

		UpdateProfile(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})
}
