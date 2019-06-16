package handler

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"bishack.dev/services/post"
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

		user := map[string]string{
			"nickname": "test",
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
		p.On("Create", mock.MatchedBy(func(vals map[string]interface{}) bool {
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
		p.On("Create", mock.MatchedBy(func(vals map[string]interface{}) bool {
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

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "postService", p)
		p.On("Get", mock.MatchedBy(func(id string) bool {
			return true
		})).Return(nil)

		GetPost(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		p.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		p := new(postMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/p/test", nil)

		context.Set(r, "postService", p)
		p.On("Get", mock.MatchedBy(func(id string) bool {
			return true
		})).Return(&post.Post{
			Title: "test",
		})

		GetPost(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("test"), w.Body.String())
		p.AssertExpectations(t)
	})
}
