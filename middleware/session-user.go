package middleware

import (
	"net/http"

	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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

		ss := context.Get(r, "session").(interface {
			GetUser(r *http.Request) map[string]string
		})

		su := ss.GetUser(r)
		if su == nil {
			return
		}

		u := context.Get(r, "userService").(interface {
			AccountDetails(token string) (*cip.GetUserOutput, error)
		})

		token := su["token"]
		o, err := u.AccountDetails(token)
		if err != nil {
			return
		}

		if len(o.UserAttributes) == 0 {
			return
		}

		ua := map[string]string{}
		for _, v := range o.UserAttributes {
			ua[*v.Name] = *v.Value
		}

		context.Set(r, "user", ua)
	})
}
