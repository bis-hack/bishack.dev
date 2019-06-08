package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRender(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()
		render(w, "xxx", "yyy", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSessionUser(t *testing.T) {
	t.Run("error account", func(t *testing.T) {
		u := new(userServiceMock)
		s := new(sessionMock)

		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "session", s)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"token": "beepboop",
		})
		u.On("AccountDetails", "beepboop").Return(nil, errors.New(""))

		su := sessionUser(r)
		assert.Nil(t, su)
	})

	t.Run("error attribute", func(t *testing.T) {
		u := new(userServiceMock)
		s := new(sessionMock)

		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", u)
		context.Set(r, "session", s)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"token": "beepboop",
		})

		resp := &cip.GetUserOutput{}
		u.On("AccountDetails", "beepboop").Return(resp, nil)

		su := sessionUser(r)
		assert.Nil(t, su)
	})
}

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
