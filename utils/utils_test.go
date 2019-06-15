package utils

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	_ "bishack.dev/testing"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Render(w, "xxx", "yyy", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		Render(w, "main", "notfound", nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGithubEndpoint(t *testing.T) {
	t.Run("has code", func(t *testing.T) {
		ep := GithubEndpoint("123")
		assert.Regexp(t, regexp.MustCompile("access_token"), ep)
	})

	t.Run("no code", func(t *testing.T) {
		ep := GithubEndpoint("123")
		assert.Equal(t, regexp.MustCompile("access_token").MatchString(ep), true)
	})
}
