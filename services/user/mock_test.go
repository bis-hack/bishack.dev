package user

import (
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/mock"
)

type MockedUserService struct {
	mock.Mock
}

func (m *MockedUserService) SignUp(in *cip.SignUpInput) (*cip.SignUpOutput, error) {
	args := m.Called(in)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.SignUpOutput), args.Error(1)
}
func (m *MockedUserService) ConfirmSignUp(in *cip.ConfirmSignUpInput) (*cip.ConfirmSignUpOutput, error) {
	args := m.Called(in)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.ConfirmSignUpOutput), args.Error(1)
}
func (m *MockedUserService) InitiateAuth(in *cip.InitiateAuthInput) (*cip.InitiateAuthOutput, error) {
	args := m.Called(in)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.InitiateAuthOutput), args.Error(1)
}
func (m *MockedUserService) GetUser(in *cip.GetUserInput) (*cip.GetUserOutput, error) {
	args := m.Called(in)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.GetUserOutput), args.Error(1)
}
func (m *MockedUserService) AdminGetUser(in *cip.AdminGetUserInput) (*cip.AdminGetUserOutput, error) {
	args := m.Called(in)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.AdminGetUserOutput), args.Error(1)
}

func (m *MockedUserService) UpdateUserAttributes(in *cip.UpdateUserAttributesInput) (*cip.UpdateUserAttributesOutput, error) {
	args := m.Called(in)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*cip.UpdateUserAttributesOutput), args.Error(1)
}
