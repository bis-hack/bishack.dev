package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"bishack.dev/services/like"
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

func TestUserProfileForm(t *testing.T) {
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

		UpdateProfileForm(w, r)

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

		UpdateProfileForm(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUserUpdate(t *testing.T) {
	t.Run("User is nil", func(t *testing.T) {
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		context.Set(r, "session", s)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		UserUpdate(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("Email is empty", func(t *testing.T) {
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		context.Set(r, "session", s)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{})

		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Email is required")

		UserUpdate(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		s.AssertExpectations(t)
	})

	t.Run("User update error", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		form := url.Values{}
		form.Set("email", "test@mailinator.com")

		r.Form = form

		u.On("GetUser")

		um := map[string]string{
			"token": "test",
		}

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(um)

		u.On("UpdateUser", um["token"], um).Return(nil, errors.New(""))

		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Email is required")

		UserUpdate(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		s.AssertExpectations(t)
	})

	t.Run("User update ok", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		form := url.Values{}
		form.Set("email", "test@mailinator.com")

		r.Form = form

		u.On("GetUser")

		um := map[string]string{
			"token": "test",
		}

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(um)

		u.On("UpdateUser", um["token"], um).Return(nil, nil)

		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "success", "Email Successfully Updated")

		UserUpdate(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		s.AssertExpectations(t)
	})
}
