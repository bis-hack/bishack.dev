package user

import cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

// CognitoProvider ...
type CognitoProvider interface {
	SignUp(*cip.SignUpInput) (*cip.SignUpOutput, error)
	ConfirmSignUp(*cip.ConfirmSignUpInput) (*cip.ConfirmSignUpOutput, error)
	InitiateAuth(*cip.InitiateAuthInput) (*cip.InitiateAuthOutput, error)
	GetUser(*cip.GetUserInput) (*cip.GetUserOutput, error)
	AdminGetUser(*cip.AdminGetUserInput) (*cip.AdminGetUserOutput, error)
}

// Client main struct
type Client struct {
	ClientID     string
	ClientSecret string
	Provider     CognitoProvider
}

// User ...
type User struct {
	ID       string
	Bio      string
	Name     string
	Email    string
	Country  string
	Website  string
	Picture  string
	Username string
}
