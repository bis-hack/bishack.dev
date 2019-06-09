package handler

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	_ "bishack.dev/testing"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHome(t *testing.T) {

	t.Run("normal", func(t *testing.T) {
		// init new mocker
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", m)
		context.Set(r, "session", s)

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)
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
		// init new mocker
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		context.Set(r, "userService", m)
		context.Set(r, "session", s)
		nickname := &cip.AttributeType{}
		nickname.SetName("nickname")
		nickname.SetValue("tibur")

		out := &cip.GetUserOutput{}
		out.SetUserAttributes([]*cip.AttributeType{nickname})

		s.On("GetUser", mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(map[string]string{
			"email": "test@user.com",
			"token": "beepboop",
		})

		m.On("AccountDetails", mock.MatchedBy(func(s string) bool {
			return true
		})).Return(out, nil)
		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Home(w, r)
		assert.Regexp(t, regexp.MustCompile("Log Out"), w.Body.String())
		assert.Regexp(t, regexp.MustCompile("@tibur"), w.Body.String())

		m.AssertExpectations(t)
		s.AssertExpectations(t)
	})
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/xxx", nil)

	NotFound(w, r)

	assert.Regexp(t, regexp.MustCompile("Not Found"), w.Body.String())
}
