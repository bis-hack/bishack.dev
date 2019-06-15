package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func TestContextMw(t *testing.T) {
	t.Run("should attach some shit to context", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/signup", nil)

		Context(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		us := context.Get(r, "userService")
		assert.NotNil(t, us)

		sess := context.Get(r, "session")
		assert.NotNil(t, sess)

		client := context.Get(r, "client")
		assert.NotNil(t, client)
	})
}
