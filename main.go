package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"

	"bishack.dev/handler"
	"bishack.dev/services/user"
	"bishack.dev/utils/session"

	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/pat"
)

// PUBLICFOLDER ...
const PUBLICFOLDER = "public"

var (
	cognitoID     = os.Getenv("COGNITO_CLIENT_ID")
	cognitoSecret = os.Getenv("COGNITO_CLIENT_SECRET")
)

func main() {
	// init route
	r := pat.New()

	// not found
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			handler.Home(w, r)
		} else {
			handler.NotFound(w, r)
		}
	})

	// handlers et al

	// GET
	r.Get("/signup", handler.Signup)
	r.Get("/verify", handler.Verify)
	r.Get("/login", handler.LoginForm)
	r.Get("/logout", handler.Logout)
	r.Get("/slack-invite", handler.SlackInvite)

	// POST
	r.Post("/signup", handler.FinishSignup)
	r.Post("/login", handler.Login)

	// csrf
	csrfSecure := false
	if regexp.MustCompile("`(?i)stag|prod").MatchString(os.Getenv("UP_STAGE")) {
		csrfSecure = true
	}
	protect := csrf.Protect([]byte(os.Getenv("CSRF_KEY")), csrf.Secure(csrfSecure))

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(
		port,
		protect(svc(r)),
	))
}

func svc(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// user service
		u := user.New(cognitoID, cognitoSecret)
		context.Set(r, "userService", u)

		// session
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

		context.Set(r, "client", c)

		h.ServeHTTP(w, r)
	})
}
