package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/pat"
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
		c.SetUser(w, r, "test@user.com", "tuku")
		u := c.GetUser(r)
		assert.NotNil(t, u)
		assert.NotEmpty(t, u["email"])
		assert.NotEmpty(t, u["token"])
		assert.Equal(t, "test@user.com", u["email"])
		assert.Equal(t, "tuku", u["token"])
	})

	t.Run("delete", func(t *testing.T) {
		c.SetUser(w, r, "test@user.com", "tuku")
		c.DeleteUser(w, r)
		u := c.GetUser(r)
		assert.Nil(t, u)
	})
}

func getRouter() *pat.Router {
	return pat.New()
}
