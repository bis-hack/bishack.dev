package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func TestRedirectMw(t *testing.T) {
	t.Run("should redirect", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/signup", nil)

		context.Set(r, "user", map[string]string{})

		AuthRedirects(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("should render default", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/login", nil)

		AuthRedirects(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
