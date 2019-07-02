package middleware

import (
	"net/http"

	"bishack.dev/services/user"
	"github.com/stretchr/testify/mock"
)

type userServiceMock struct {
	mock.Mock
}

func (o *userServiceMock) AccountDetails(token string) *user.User {
	args := o.Called(token)

	resp := args.Get(0)
	if resp == nil {
		return nil
	}

	return resp.(*user.User)
}

func (o *userServiceMock) GetToken(username, token string) (string, error) {
	args := o.Called(username, token)
	return args.String(0), args.Error(1)
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
