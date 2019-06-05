package user

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
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

func TestSignUp(t *testing.T) {
	to := new(MockedUserService)
	client := New("id", "secret")
	// change provider to our mocked object
	client.Provider = to
	t.Run("pass proper inputs", func(t *testing.T) {
		to.On(
			"SignUp",
			mock.MatchedBy(func(in *cip.SignUpInput) bool {
				if *in.Username == "beep" && *in.Password == "boop" {
					return true
				}
				return false
			}),
		).Return(&cip.SignUpOutput{UserSub: aws.String("beep")}, nil)

		o, err := client.Signup("beep", "boop", "beep@boop.com")

		assert.Equal(t, err, nil)
		assert.Equal(t, "beep", *o.UserSub)
		to.AssertExpectations(t)
	})
}

func TestVerify(t *testing.T) {
	to := new(MockedUserService)
	client := New("id", "secret")
	// change provider to our mocked object
	client.Provider = to
	t.Run("pass proper inputs", func(t *testing.T) {
		to.On(
			"ConfirmSignUp",
			mock.MatchedBy(func(in *cip.ConfirmSignUpInput) bool {
				if *in.Username == "beep" && *in.ConfirmationCode == "boop" {
					return true
				}
				return false
			}),
		).Return(&cip.ConfirmSignUpOutput{}, nil)

		_, err := client.Verify("beep", "boop")

		assert.Equal(t, err, nil)
		to.AssertExpectations(t)
	})
}

func TestAccountDetails(t *testing.T) {
	to := new(MockedUserService)
	client := New("id", "secret")
	// change provider to our mocked object
	client.Provider = to
	t.Run("valid token", func(t *testing.T) {
		to.On(
			"GetUser",
			mock.MatchedBy(func(in *cip.GetUserInput) bool {
				if *in.AccessToken == "wadiwasi" {
					return true
				}
				return false
			}),
		).Return(&cip.GetUserOutput{
			Username: aws.String("beep"),
		}, nil)

		o, err := client.AccountDetails("wadiwasi")

		assert.Equal(t, err, nil)
		assert.Equal(t, "beep", *o.Username)
		to.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	to := new(MockedUserService)
	client := New("id", "secret")
	// change provider to our mocked object
	client.Provider = to
	t.Run("valid token", func(t *testing.T) {
		to.On(
			"InitiateAuth",
			mock.MatchedBy(func(in *cip.InitiateAuthInput) bool {
				if v, exists := in.AuthParameters["USERNAME"]; exists && *v == "beep" {
					return true
				}
				return false
			}),
		).Return(&cip.InitiateAuthOutput{
			AuthenticationResult: &cip.AuthenticationResultType{
				AccessToken: aws.String("wadiwasitoken"),
			},
		}, nil)

		o, err := client.Login("beep", "boop")

		assert.Equal(t, err, nil)
		assert.Equal(t, "wadiwasitoken", *o.AuthenticationResult.AccessToken)
		to.AssertExpectations(t)
	})
}
