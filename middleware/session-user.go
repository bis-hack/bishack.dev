package middleware

import (
	"net/http"

	"bishack.dev/services/user"
	"github.com/gorilla/context"
)

// SessionUser middleware checks for the `user` session and if it exists
// it will try to fetch the user details from Cognito service and attach
// them to the context object.
func SessionUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			h.ServeHTTP(w, r)
		}()

		token := context.Get(r, "token")
		if token == nil {
			return
		}

		u := context.Get(r, "userService").(interface {
			AccountDetails(token string) *user.User
		})

		user := u.AccountDetails(token.(string))
		if user == nil {
			return
		}

		context.Set(r, "user", user)
	})
}

// Token queries Cognito for a fresh access token from the refresh token
// that is saved under the user session.
// Main reason for this decision is to solve Cognito's short-lived session
// lifespan.
func Token(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			h.ServeHTTP(w, r)
		}()

		ses := context.Get(r, "session").(interface {
			GetUser(r *http.Request) map[string]string
		})

		su := ses.GetUser(r)
		if su == nil {
			return
		}

		us := context.Get(r, "userService").(interface {
			GetToken(string, string) (string, error)
		})

		username := su["username"]
		token := su["token"]
		accessToken, err := us.GetToken(username, token)
		if err != nil {
			return
		}

		context.Set(r, "token", accessToken)
	})
}
