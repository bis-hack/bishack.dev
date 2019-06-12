package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	_ "bishack.dev/testing"
	"bishack.dev/utils/session"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginForm(t *testing.T) {
	s := new(sessionMock)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	context.Set(r, "session", s)

	s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
		return true
	}), mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return(nil)

	LoginForm(w, r)

	assert.Regexp(t, regexp.MustCompile("User Login"), w.Body.String())
	assert.Regexp(t, regexp.MustCompile("login-form"), w.Body.String())
	s.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	t.Run("wrong username or password", func(t *testing.T) {
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/login", nil)

		form := url.Values{}
		form.Add("email", "")
		form.Add("password", "")
		r.PostForm = form

		context.Set(r, "userService", m)
		context.Set(r, "session", s)
		m.On("Login", "", "").Return(nil, errors.New("wrong username/password"))
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Wrong email or password")

		Login(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		m.AssertExpectations(t)
		s.AssertExpectations(t)
	})

	t.Run("login success", func(t *testing.T) {
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/login", nil)

		form := url.Values{}
		form.Add("email", "test@user.com")
		form.Add("password", "test")
		r.PostForm = form

		context.Set(r, "userService", m)
		context.Set(r, "session", s)

		result := &cip.AuthenticationResultType{}
		result.SetAccessToken("beepboop")
		out := &cip.InitiateAuthOutput{}
		out.SetAuthenticationResult(result)

		m.On("Login", "test@user.com", "test").Return(out, nil)
		s.On("SetUser", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "test@user.com", "beepboop")
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "success", "Welcome Back!")

		Login(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)

		m.AssertExpectations(t)
		s.AssertExpectations(t)
	})
}

func TestLogout(t *testing.T) {
	s := new(sessionMock)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/login", nil)

	context.Set(r, "session", s)

	s.On("DeleteUser", mock.MatchedBy(func(w http.ResponseWriter) bool {
		return true
	}), mock.MatchedBy(func(r *http.Request) bool {
		return true
	}))

	Logout(w, r)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	s.AssertExpectations(t)
}

func TestVerify(t *testing.T) {
	t.Run("with invalid code", func(t *testing.T) {
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/verify?code=111&email=test@user.com", nil)

		context.Set(r, "userService", m)
		context.Set(r, "session", s)

		m.On("Verify", "test@user.com", "111").Return(nil, errors.New(""))
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Verification failed. Try again!")

		Verify(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)

		m.AssertExpectations(t)
		s.AssertExpectations(t)
	})

	t.Run("with valid code", func(t *testing.T) {
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/verify?code=111&email=test@user.com", nil)

		context.Set(r, "userService", m)
		context.Set(r, "session", s)

		out := &cip.ConfirmSignUpOutput{}
		m.On("Verify", "test@user.com", "111").Return(out, nil)
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "success", "Account Verified!")

		Verify(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)

		m.AssertExpectations(t)
		s.AssertExpectations(t)
	})

	t.Run("verified", func(t *testing.T) {
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/verify", nil)

		context.Set(r, "userService", m)
		context.Set(r, "session", s)

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(&session.Flash{
			"success",
			"Account Verified!",
		})

		Verify(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("Account Verified!"), w.Body.String())

		s.AssertExpectations(t)
	})

	t.Run("form", func(t *testing.T) {
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/verify", nil)

		context.Set(r, "userService", m)
		context.Set(r, "session", s)

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Verify(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("Verify"), w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestFinishSignup(t *testing.T) {
	t.Run("error signup", func(t *testing.T) {
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("email", "test@user.com")
		form.Add("password", "beepboop")
		r.PostForm = form

		context.Set(r, "userService", m)
		context.Set(r, "session", s)

		m.On("Signup", "test@user.com", "beepboop", mock.MatchedBy(func(m map[string]string) bool {
			return true
		})).Return(nil, errors.New(""))
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Could not sign you up. Try again!")

		FinishSignup(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)

		m.AssertExpectations(t)
		s.AssertExpectations(t)
	})

	t.Run("error signup: account exists", func(t *testing.T) {
		m := new(userServiceMock)
		s := new(sessionMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("email", "test@user.com")
		form.Add("password", "beepboop")
		r.PostForm = form

		context.Set(r, "userService", m)
		context.Set(r, "session", s)

		m.On("Signup", "test@user.com", "beepboop", mock.MatchedBy(func(m map[string]string) bool {
			return true
		})).Return(nil, errors.New("exists"))
		s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		}), "error", "Account already exists. Try to log in instead.")

		FinishSignup(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)

		m.AssertExpectations(t)
		s.AssertExpectations(t)
	})

	t.Run("signup success", func(t *testing.T) {
		m := new(userServiceMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/signup", nil)

		context.Set(r, "userService", m)

		m.On("Signup", "", "", mock.MatchedBy(func(m map[string]string) bool {
			return true
		})).Return(nil, nil)

		FinishSignup(w, r)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		m.AssertExpectations(t)
	})
}

func TestSignup(t *testing.T) {
	t.Run("oauth code", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			s := new(sessionMock)
			c := new(clientMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/signup?code=123", nil)

			context.Set(r, "session", s)
			context.Set(r, "client", c)

			c.On("PostForm", mock.MatchedBy(func(url string) bool {
				return true
			}), mock.MatchedBy(func(data url.Values) bool {
				return true
			})).Return(nil, errors.New("error"))
			s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
				return true
			}), mock.MatchedBy(func(r *http.Request) bool {
				return true
			}), "error", "Invalid or expired code")

			Signup(w, r)

			assert.Equal(t, http.StatusSeeOther, w.Code)
			s.AssertExpectations(t)
			c.AssertExpectations(t)
		})

		t.Run("token error", func(t *testing.T) {
			s := new(sessionMock)
			c := new(clientMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/signup?code=123", nil)

			context.Set(r, "session", s)
			context.Set(r, "client", c)

			resp := &http.Response{}
			resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`?beep=boop`)))
			c.On("PostForm", mock.MatchedBy(func(url string) bool {
				return true
			}), mock.MatchedBy(func(data url.Values) bool {
				return true
			})).Return(resp, nil)
			s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
				return true
			}), mock.MatchedBy(func(r *http.Request) bool {
				return true
			}), "error", "Invalid or expired code")

			Signup(w, r)

			assert.Equal(t, http.StatusSeeOther, w.Code)
			s.AssertExpectations(t)
			c.AssertExpectations(t)
		})

		t.Run("code ok", func(t *testing.T) {
			s := new(sessionMock)
			c := new(clientMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/signup?code=123", nil)

			context.Set(r, "session", s)
			context.Set(r, "client", c)

			resp := &http.Response{}
			resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`access_token=123`)))
			c.On("PostForm", mock.MatchedBy(func(url string) bool {
				return true
			}), mock.MatchedBy(func(data url.Values) bool {
				return true
			})).Return(resp, nil)

			Signup(w, r)

			assert.Equal(t, http.StatusSeeOther, w.Code)
			s.AssertExpectations(t)
			c.AssertExpectations(t)
		})
	})

	t.Run("access_token", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			s := new(sessionMock)
			c := new(clientMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/signup?access_token=123", nil)

			context.Set(r, "session", s)
			context.Set(r, "client", c)

			c.On("Do", mock.MatchedBy(func(r *http.Request) bool {
				return true
			})).Return(nil, errors.New("error"))
			s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
				return true
			}), mock.MatchedBy(func(r *http.Request) bool {
				return true
			}), "error", "Invalid or expired token!")

			Signup(w, r)

			assert.Equal(t, http.StatusSeeOther, w.Code)
			s.AssertExpectations(t)
			c.AssertExpectations(t)
		})

		t.Run("json error", func(t *testing.T) {
			s := new(sessionMock)
			c := new(clientMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/signup?access_token=123", nil)

			context.Set(r, "session", s)
			context.Set(r, "client", c)

			resp := &http.Response{}
			resp.StatusCode = http.StatusOK
			resp.Header = http.Header{}
			resp.Header.Set("content-type", "application/json")
			resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`hello\n`)))

			c.On("Do", mock.MatchedBy(func(r *http.Request) bool {
				return true
			})).Return(resp, nil)
			s.On("SetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
				return true
			}), mock.MatchedBy(func(r *http.Request) bool {
				return true
			}), "error", "An error occured!")

			Signup(w, r)

			assert.Equal(t, http.StatusSeeOther, w.Code)
			c.AssertExpectations(t)
			s.AssertExpectations(t)
		})

		t.Run("ok", func(t *testing.T) {
			s := new(sessionMock)
			c := new(clientMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/signup?access_token=123", nil)

			context.Set(r, "session", s)
			context.Set(r, "client", c)

			resp := &http.Response{}
			resp.StatusCode = http.StatusOK
			resp.Header = http.Header{}
			resp.Header.Set("content-type", "application/json")
			resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"email":"test@user.com"}`)))

			c.On("Do", mock.MatchedBy(func(r *http.Request) bool {
				return true
			})).Return(resp, nil)
			s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
				return true
			}), mock.MatchedBy(func(r *http.Request) bool {
				return true
			})).Return(nil)

			Signup(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			c.AssertExpectations(t)
			s.AssertExpectations(t)
		})
	})

	t.Run("default", func(t *testing.T) {
		s := new(sessionMock)
		c := new(clientMock)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/signup", nil)

		context.Set(r, "session", s)
		context.Set(r, "client", c)

		s.On("GetFlash", mock.MatchedBy(func(w http.ResponseWriter) bool {
			return true
		}), mock.MatchedBy(func(r *http.Request) bool {
			return true
		})).Return(nil)

		Signup(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Regexp(t, regexp.MustCompile("(?i)connect"), w.Body.String())
	})
}
