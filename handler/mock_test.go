package handler

import (
	"net/http"
	"net/url"

	"bishack.dev/services/like"
	"bishack.dev/services/post"
	"bishack.dev/services/user"
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

func (o *userServiceMock) UpdateUser(token string, attrs map[string]string) (*cip.UpdateUserAttributesOutput, error) {
	args := o.Called(token, attrs)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.UpdateUserAttributesOutput), args.Error(1)
}

func (o *userServiceMock) GetUser(username string) *user.User {
	args := o.Called(username)

	resp := args.Get(0)
	if resp == nil {
		return nil
	}

	return resp.(*user.User)
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

func (p *postMock) CreatePost(vals map[string]interface{}) *post.Post {
	args := p.Called(vals)
	resp := args.Get(0)

	if resp == nil {
		return nil
	}

	return resp.(*post.Post)
}

func (p *postMock) GetPosts() []*post.Post {
	args := p.Called()
	resp := args.Get(0)

	if resp == nil {
		return nil
	}

	return resp.([]*post.Post)
}

func (p *postMock) GetUserPosts(username string) []*post.Post {
	args := p.Called(username)
	resp := args.Get(0)

	if resp == nil {
		return nil
	}

	return resp.([]*post.Post)
}

func (p *postMock) UpdatePost(id, cover, content string, created int64) error {
	args := p.Called(id, cover, content, created)
	_ = args.Get(0)
	return args.Error(0)
}

func (p *postMock) GetPost(username, id string) *post.Post {
	args := p.Called()
	resp := args.Get(0)

	if resp == nil {
		return nil
	}

	return resp.(*post.Post)
}

type likeMock struct {
	mock.Mock
}

func (l *likeMock) GetLikes(id string) ([]*like.Like, error) {
	args := l.Called(id)
	resp := args.Get(0)

	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.([]*like.Like), args.Error(1)
}

func (l *likeMock) GetLike(id, username string) (*like.Like, error) {
	args := l.Called(id, username)
	resp := args.Get(0)

	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*like.Like), args.Error(1)
}

func (l *likeMock) ToggleLike(id, username string) error {
	args := l.Called(id, username)
	return args.Error(0)
}
