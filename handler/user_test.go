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
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		UpdateProfile(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/update", nil)

		context.Set(r, "token", "test")
		context.Set(r, "session", s)
		context.Set(r, "userService", u)

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

		context.Set(r, "token", "test")
		context.Set(r, "session", s)
		context.Set(r, "userService", u)

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

func TestChangePassword(t *testing.T) {
	t.Run("show form", func(t *testing.T) {
		s := new(sessionMock)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/security", nil)

		context.Set(r, "session", s)
		context.Set(r, "user", map[string]string{})

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Security(w, r)

		assert.Regexp(t, regexp.MustCompile("security-form"), w.Body.String())
		s.AssertExpectations(t)
	})

	t.Run("user nil", func(t *testing.T) {
		s := new(sessionMock)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/security", nil)

		context.Set(r, "session", s)
		context.Set(r, "user", nil)

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Security(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("nil token", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/security", nil)

		ChangePassword(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("password mismatch", func(t *testing.T) {
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/security", nil)

		form := url.Values{}
		form.Add("old", "old_password")
		form.Add("new", "new_password")
		form.Add("confirm", "new_passw0rd")
		r.PostForm = form

		context.Set(r, "token", "test")
		context.Set(r, "session", s)

		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Password confirmation doesn't match the password")

		ChangePassword(w, r)
		s.AssertExpectations(t)
		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("incorrect password", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/security", nil)

		form := url.Values{}
		form.Add("old", "incorrect_password")
		form.Add("new", "new_password")
		form.Add("confirm", "new_password")
		r.PostForm = form

		context.Set(r, "token", "test")
		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		u.On("ChangePassword", "test", "incorrect_password", "new_password").Return(nil, errors.New("Incorrect Password"))
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Incorrect Password")

		ChangePassword(w, r)

		u.AssertExpectations(t)
		s.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/security", nil)

		form := url.Values{}
		form.Add("old", "pass")
		form.Add("new", "passs")
		form.Add("confirm", "passs")
		r.PostForm = form

		context.Set(r, "token", "test")
		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		u.On("ChangePassword", "test", "pass", "passs").Return(nil, errors.New("Invalid Password"))
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Invalid Password")

		ChangePassword(w, r)

		u.AssertExpectations(t)
		s.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/security", nil)

		form := url.Values{}
		form.Add("old", "old_password")
		form.Add("new", "new_password")
		form.Add("confirm", "new_password")
		r.PostForm = form

		context.Set(r, "token", "test")
		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		u.On("ChangePassword", "test", "old_password", "new_password").Return(&cip.ChangePasswordOutput{}, nil)
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "success", "Password Successfully Updated")

		ChangePassword(w, r)

		u.AssertExpectations(t)
		s.AssertExpectations(t)
	})
}
