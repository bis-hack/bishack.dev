package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// New creates a new instance of Client
func New(id, secret string) *Client {
	return &Client{id, secret, provider()}
}

// Signup ...
func (c *Client) Signup(
	username,
	password string,
	meta map[string]string,
) (*cip.SignUpOutput, error) {
	input := &cip.SignUpInput{}

	input.SetSecretHash(hash(username, c.ClientID, c.ClientSecret))
	input.SetClientId(c.ClientID)
	input.SetUsername(username)
	input.SetPassword(password)

	userAttributes := []*cip.AttributeType{}
	for k, v := range meta {
		a := &cip.AttributeType{}
		a.SetName(k)
		a.SetValue(v)
		userAttributes = append(userAttributes, a)
	}
	input.SetUserAttributes(userAttributes)

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
func (c *Client) AccountDetails(token string) *User {
	input := &cip.GetUserInput{}
	input.SetAccessToken(token)

	out, err := c.Provider.GetUser(input)
	if err != nil {
		log.Println("GetUser error", err.Error())
		return nil
	}

	return newUserFromAttributes(out.UserAttributes)
}

// GetUser ...
func (c *Client) GetUser(username string) *User {
	pid := os.Getenv("COGNITO_POOL_ID")

	input := &cip.AdminGetUserInput{}
	input.SetUserPoolId(pid)
	input.SetUsername(username)

	out, err := c.Provider.AdminGetUser(input)
	if err != nil {
		log.Println("AdminGetUser error:", err.Error())
		return nil
	}

	return newUserFromAttributes(out.UserAttributes)
}

// UpdateUser a user attribute
func (c *Client) UpdateUser(token string, attributes map[string]string) (*cip.UpdateUserAttributesOutput, error) {
	input := &cip.UpdateUserAttributesInput{}

	input.SetAccessToken(token)

	userAttributes := []*cip.AttributeType{}
	for k, v := range attributes {
		a := &cip.AttributeType{}
		a.SetName(k)
		a.SetValue(v)
		userAttributes = append(userAttributes, a)
	}

	input.SetUserAttributes(userAttributes)

	return c.Provider.UpdateUserAttributes(input)
}

//
// PRIVATE
//

// newUserFromOutput ...
func newUserFromAttributes(attrs []*cip.AttributeType) *User {
	user := &User{}

	// attr mappings
	am := map[string]string{}
	for _, a := range attrs {
		am[*a.Name] = *a.Value
	}

	user.ID = am["sub"]
	user.Bio = am["profile"]
	user.Name = am["name"]
	user.Email = am["email"]
	user.Country = am["locale"]
	user.Website = am["website"]
	user.Picture = am["picture"]
	user.Username = am["nickname"]

	return user
}

// provider returns a new cognito identity service
func provider() Provider {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	return cip.New(sess)
}

// hash that shit
func hash(username, id, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	_, _ = hash.Write([]byte(username + id))
	return string(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
}
