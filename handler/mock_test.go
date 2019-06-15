package handler

import (
	"net/http"
	"net/url"

	"bishack.dev/utils/session"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/mock"
)

type userServiceMock struct {
	mock.Mock
}

func (o *userServiceMock) Login(username, password string) (*cip.InitiateAuthOutput, error) {
	args := o.Called(username, password)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.InitiateAuthOutput), args.Error(1)
}

func (o *userServiceMock) Signup(
	username,
	password string,
	meta map[string]string,
) (*cip.SignUpOutput, error) {
	args := o.Called(username, password, meta)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.SignUpOutput), args.Error(1)
}

func (o *userServiceMock) AccountDetails(token string) (*cip.GetUserOutput, error) {
	args := o.Called(token)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.GetUserOutput), args.Error(1)
}

func (o *userServiceMock) Verify(
	username,
	code string,
) (*cip.ConfirmSignUpOutput, error) {
	args := o.Called(username, code)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.ConfirmSignUpOutput), args.Error(1)
}

type sessionMock struct {
	mock.Mock
}

func (o *sessionMock) GetUser(r *http.Request) map[string]string {
	args := o.Called(r)

	resp := args.Get(0)
	if resp == nil {
		return nil
	}

	return resp.(map[string]string)
}

func (o *sessionMock) SetUser(
	w http.ResponseWriter,
	r *http.Request,
	email,
	token string,
) {
	o.Called(w, r, email, token)
}

func (o *sessionMock) DeleteUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	o.Called(w, r)
}

func (o *sessionMock) SetFlash(
	w http.ResponseWriter,
	r *http.Request,
	t,
	v string,
) {
	o.Called(w, r, t, v)
}

func (o *sessionMock) GetFlash(
	w http.ResponseWriter,
	r *http.Request,
) *session.Flash {
	args := o.Called(w, r)

	resp := args.Get(0)
	if resp == nil {
		return nil
	}

	return resp.(*session.Flash)
}

type clientMock struct {
	mock.Mock
}

func (c *clientMock) PostForm(url string, data url.Values) (*http.Response, error) {
	args := c.Called(url, data)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*http.Response), args.Error(1)
}

func (c *clientMock) Do(r *http.Request) (*http.Response, error) {
	args := c.Called(r)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*http.Response), args.Error(1)
}

func (c *clientMock) Get(url string) (*http.Response, error) {
	args := c.Called(url)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*http.Response), args.Error(1)
}

type postMock struct {
	mock.Mock
}

func (p *postMock) GetCount() int64 {
	args := p.Called()
	resp := args.Get(0)
	return resp.(int64)
}
