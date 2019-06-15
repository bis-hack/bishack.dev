package handler

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	_ "bishack.dev/testing"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHome(t *testing.T) {

	t.Run("normal", func(t *testing.T) {
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "session", s)

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
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "user", map[string]string{
			"email":    "test@user.com",
			"nickname": "tibur",
		})
		context.Set(r, "session", s)

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

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
