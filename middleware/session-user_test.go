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
	t.Run("session not found", func(t *testing.T) {
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/signup", nil)

		context.Set(r, "session", s)
		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		SessionUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		u := context.Get(r, "user")
		assert.Nil(t, u)
		s.AssertExpectations(t)
	})

	t.Run("error account details", func(t *testing.T) {
		u := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/signup", nil)

		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"token": "test",
		})

		u.On("AccountDetails", "test").Return(nil, errors.New(""))

		SessionUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		user := context.Get(r, "user")
		assert.Nil(t, user)
		u.AssertExpectations(t)
		s.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		u := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/signup", nil)

		context.Set(r, "session", s)
		context.Set(r, "userService", u)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"token": "test",
		})

		resp := &user.User{}
		u.On("AccountDetails", "test").Return(resp, nil)

		SessionUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		user := context.Get(r, "user")
		assert.NotNil(t, user)
		u.AssertExpectations(t)
		s.AssertExpectations(t)
	})
}
