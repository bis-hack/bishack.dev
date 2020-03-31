package middleware

import (
	"net"
	"net/http"
	"os"
	"time"

	"bishack.dev/services/like"
	"bishack.dev/services/post"
	"bishack.dev/services/user"
	"bishack.dev/utils/session"
	"github.com/gorilla/context"

	// autoload env
	_ "github.com/joho/godotenv/autoload"
)

var (
	cognitoID        = os.Getenv("COGNITO_CLIENT_ID")
	cognitoSecret    = os.Getenv("COGNITO_CLIENT_SECRET")
	dynamoTablePosts = os.Getenv("DYNAMO_TABLE_POSTS")
	dynamoTableLikes = os.Getenv("DYNAMO_TABLE_LIKES")
	dynamoEndpoint   = os.Getenv("DYNAMO_ENDPOINT")
)

// Context middleware will inject services, helpers and other utility code
// to the context object
func Context(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// user service
		u := user.New(cognitoID, cognitoSecret)
		context.Set(r, "userService", u)

		// session helper
		s := session.New()
		context.Set(r, "session", s)

		var netTransport = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		}

		// support timeout and net transport.
		c := &http.Client{
			Timeout:   time.Second * 10,
			Transport: netTransport,
		}

		// http client
		context.Set(r, "client", c)

		p := post.New(dynamoTablePosts, dynamoEndpoint, nil)
		context.Set(r, "postService", p)

		l := like.New(dynamoTableLikes, dynamoEndpoint, nil)
		context.Set(r, "likeService", l)

		h.ServeHTTP(w, r)
	})
}
