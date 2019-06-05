package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// CognitoProvider ...
type CognitoProvider interface {
	SignUp(*cip.SignUpInput) (*cip.SignUpOutput, error)
	ConfirmSignUp(*cip.ConfirmSignUpInput) (*cip.ConfirmSignUpOutput, error)
	InitiateAuth(*cip.InitiateAuthInput) (*cip.InitiateAuthOutput, error)
	GetUser(*cip.GetUserInput) (*cip.GetUserOutput, error)
}

// Client main struct
type Client struct {
	ClientID     string
	ClientSecret string
	Provider     CognitoProvider
}

// New creates a new instance of Client
func New(id, secret string) *Client {
	return &Client{id, secret, provider()}
}

// Signup ...
func (c *Client) Signup(
	username,
	password,
	email string,
) (*cip.SignUpOutput, error) {
	input := &cip.SignUpInput{}

	input.SetSecretHash(hash(username, c.ClientID, c.ClientSecret))
	input.SetClientId(c.ClientID)
	input.SetUsername(username)
	input.SetPassword(password)

	emailAttribute := &cip.AttributeType{}
	emailAttribute.SetName("email")
	emailAttribute.SetValue(email)

	input.SetUserAttributes([]*cip.AttributeType{
		emailAttribute,
	})

	return c.Provider.SignUp(input)
}

// Verify ...
func (c *Client) Verify(
	username,
	code string,
) (*cip.ConfirmSignUpOutput, error) {
	input := &cip.ConfirmSignUpInput{}

	input.SetSecretHash(hash(username, c.ClientID, c.ClientSecret))
	input.SetClientId(c.ClientID)
	input.SetUsername(username)
	input.SetConfirmationCode(code)

	return c.Provider.ConfirmSignUp(input)
}

// Login ...
func (c *Client) Login(
	username,
	password string,
) (*cip.InitiateAuthOutput, error) {
	input := &cip.InitiateAuthInput{}

	input.SetClientId(c.ClientID)

	secretHash := hash(username, c.ClientID, c.ClientSecret)
	input.SetAuthFlow(cip.AuthFlowTypeUserPasswordAuth)
	input.SetAuthParameters(map[string]*string{
		"USERNAME":    &username,
		"PASSWORD":    &password,
		"SECRET_HASH": &secretHash,
	})

	return c.Provider.InitiateAuth(input)
}

// AccountDetails ...
func (c *Client) AccountDetails(token string) (*cip.GetUserOutput, error) {
	input := &cip.GetUserInput{}
	input.SetAccessToken(token)
	return c.Provider.GetUser(input)
}

//
// PRIVATE
//

// provider returns a new cognito identity service
func provider() CognitoProvider {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	return cip.New(sess)
}

// hash that shit
func hash(username, id, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(username + id))
	return string(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
}
