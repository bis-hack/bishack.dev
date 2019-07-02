package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bishack.dev/services/user"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSessionUserMw(t *testing.T) {
	t.Run("token not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		SessionUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		u := context.Get(r, "user")
		assert.Nil(t, u)
	})

	t.Run("error account details", func(t *testing.T) {
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "token", "test")

		u.On("AccountDetails", "test").Return(nil, errors.New(""))

		SessionUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		user := context.Get(r, "user")
		assert.Nil(t, user)
		u.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "token", "test")

		resp := &user.User{}
		u.On("AccountDetails", "test").Return(resp, nil)

		SessionUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		user := context.Get(r, "user")
		assert.NotNil(t, user)
		u.AssertExpectations(t)
	})
}

func TestTokenMw(t *testing.T) {
	t.Run("session user not found", func(t *testing.T) {
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "session", s)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Token(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		token := context.Get(r, "token")
		assert.Nil(t, token)
	})

	t.Run("error", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"username": "test",
			"token":    "test",
		})
		u.On("GetToken", "test", "test").Return("", errors.New(""))

		Token(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		token := context.Get(r, "token")
		assert.Nil(t, token)
	})

	t.Run("ok", func(t *testing.T) {
		s := new(sessionMock)
		u := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"username": "test",
			"token":    "test",
		})
		u.On("GetToken", "test", "test").Return("test", nil)

		Token(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		token := context.Get(r, "token")
		assert.Equal(t, "test", token.(string))
	})
}
