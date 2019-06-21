package user

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

		o, err := client.Signup("beep", "boop", map[string]string{
			"beep": "boop",
		})

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
	t.Run("valid token", func(t *testing.T) {
		to := new(MockedUserService)
		client := New("id", "secret")
		// change provider to our mocked object
		client.Provider = to
		to.On(
			"GetUser",
			mock.MatchedBy(func(in *cip.GetUserInput) bool {
				return true
			}),
		).Return(&cip.GetUserOutput{
			Username: aws.String("test"),
			UserAttributes: []*cip.AttributeType{
				{
					Name:  aws.String("email"),
					Value: aws.String("test@testing.com"),
				},
			},
		}, nil)

		u := client.AccountDetails("test")

		assert.Equal(t, "test@testing.com", u.Email)
		to.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		to := new(MockedUserService)
		client := New("id", "secret")
		// change provider to our mocked object
		client.Provider = to
		to.On(
			"GetUser",
			mock.MatchedBy(func(in *cip.GetUserInput) bool {
				return true
			}),
		).Return(nil, errors.New(""))

		u := client.AccountDetails("test")

		assert.Nil(t, u)
		to.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		to := new(MockedUserService)
		client := New("id", "secret")
		// change provider to our mocked object
		client.Provider = to
		to.On(
			"AdminGetUser",
			mock.MatchedBy(func(in *cip.AdminGetUserInput) bool {
				return true
			}),
		).Return(nil, errors.New(""))

		u := client.GetUser("test")

		assert.Nil(t, u)
		to.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		to := new(MockedUserService)
		client := New("id", "secret")
		// change provider to our mocked object
		client.Provider = to
		to.On(
			"AdminGetUser",
			mock.MatchedBy(func(in *cip.AdminGetUserInput) bool {
				return true
			}),
		).Return(&cip.AdminGetUserOutput{}, nil)

		u := client.GetUser("test")

		assert.NotNil(t, u)
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

func TestUpdateUserAttributes(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		to := new(MockedUserService)
		client := New("id", "secret")
		// change provider to our mocked object
		client.Provider = to
		to.On(
			"UpdateUserAttributes",
			mock.MatchedBy(func(in *cip.UpdateUserAttributesInput) bool {
				return true
			}),
		).Return(&cip.UpdateUserAttributesOutput{
			CodeDeliveryDetailsList: []*cip.CodeDeliveryDetailsType{
				&cip.CodeDeliveryDetailsType{
					AttributeName:  aws.String("email"),
					DeliveryMedium: aws.String("email"),
					Destination:    aws.String("richard@mailinator.com"),
				},
			},
		}, nil)

		attrs := map[string]string{
			"email": "richard@mail.co",
		}

		u := client.UpdateUser("legit_token", attrs)

		assert.NoError(t, u, "Error updating")
		to.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		to := new(MockedUserService)
		client := New("id", "secret")
		// change provider to our mocked object
		client.Provider = to
		to.On(
			"UpdateUserAttributes",
			mock.MatchedBy(func(in *cip.UpdateUserAttributesInput) bool {
				return true
			}),
		).Return(nil, errors.New(""))

		u := client.UpdateUser("token", map[string]string{
			"email": "chardyy.orofeo@gmail.com",
		})

		assert.Nil(t, u)
		to.AssertExpectations(t)
	})
}
