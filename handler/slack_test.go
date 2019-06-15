package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSlackInvite(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		c := new(clientMock)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "client", c)

		c.On("Get", mock.MatchedBy(func(url string) bool {
			return true
		})).Return(nil, errors.New(""))

		SlackInvite(w, r)
		assert.JSONEq(t, `{"ok":false}`, w.Body.String())
	})

	t.Run("success", func(t *testing.T) {
		c := new(clientMock)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "client", c)

		resp := &http.Response{}
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"ok":true}`)))
		c.On("Get", mock.MatchedBy(func(url string) bool {
			return true
		})).Return(resp, nil)

		SlackInvite(w, r)
		assert.JSONEq(t, `{"ok":true}`, w.Body.String())
	})
}
