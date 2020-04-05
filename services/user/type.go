package user

import cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

// Provider ...
type Provider interface {
	SignUp(*cip.SignUpInput) (*cip.SignUpOutput, error)
	ConfirmSignUp(*cip.ConfirmSignUpInput) (*cip.ConfirmSignUpOutput, error)
	InitiateAuth(*cip.InitiateAuthInput) (*cip.InitiateAuthOutput, error)
	GetUser(*cip.GetUserInput) (*cip.GetUserOutput, error)
	AdminGetUser(*cip.AdminGetUserInput) (*cip.AdminGetUserOutput, error)
	UpdateUserAttributes(*cip.UpdateUserAttributesInput) (*cip.UpdateUserAttributesOutput, error)
	ChangePassword(*cip.ChangePasswordInput) (*cip.ChangePasswordOutput, error)
}

// Client main struct
type Client struct {
	ClientID     string
	ClientSecret string
	Provider     Provider
}

// User ...
type User struct {
	ID       string
	Bio      string
	Name     string
	Email    string
	Location string
	Website  string
	Picture  string
	Username string
}
