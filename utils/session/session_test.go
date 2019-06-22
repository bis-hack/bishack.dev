package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlash(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	c := New()

	c.SetFlash(w, r, "error", "error")

	flash := c.GetFlash(w, r)
	assert.NotNil(t, flash)
	assert.Equal(t, flash.Type, "error")
	assert.Equal(t, flash.Value, "error")

	flash = c.GetFlash(w, r)
	assert.Nil(t, flash)
}

func TestUser(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	c := New()

	t.Run("email or token nil", func(t *testing.T) {
		u := c.GetUser(r)
		assert.Nil(t, u)
	})

	t.Run("ok", func(t *testing.T) {
		c.SetUser(w, r, "test", "tuku")
		u := c.GetUser(r)
		assert.NotNil(t, u)
		assert.NotEmpty(t, u["username"])
		assert.NotEmpty(t, u["token"])
		assert.Equal(t, "test", u["username"])
		assert.Equal(t, "tuku", u["token"])
	})

	t.Run("delete", func(t *testing.T) {
		c.SetUser(w, r, "test", "tuku")
		c.DeleteUser(w, r)
		u := c.GetUser(r)
		assert.Nil(t, u)
	})
}
