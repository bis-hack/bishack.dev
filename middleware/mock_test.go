package middleware

import (
	"net/http"

	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/mock"
)

type userServiceMock struct {
	mock.Mock
}

func (o *userServiceMock) AccountDetails(token string) (*cip.GetUserOutput, error) {
	args := o.Called(token)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.GetUserOutput), args.Error(1)
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
