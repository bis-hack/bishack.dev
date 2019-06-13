package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
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

var (
	cognitoID     = os.Getenv("COGNITO_CLIENT_ID")
	cognitoSecret = os.Getenv("COGNITO_CLIENT_SECRET")
	rxEnv         = regexp.MustCompile("`(?i)stag|prod")
)

func main() {
	// csrf
	csrfSecure := true
	isLive := rxEnv.MatchString(os.Getenv("UP_STAGE"))

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

	// localhost
	if !isLive {
		// set secure to false
		csrfSecure = false

		// launch nerdy stuff(pprof) server
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}

	protect := csrf.Protect([]byte(os.Getenv("CSRF_KEY")), csrf.Secure(csrfSecure))

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(
		port,
		protect(ctxMw(authMw(r))),
	))
}

// context middleware
func ctxMw(h http.Handler) http.Handler {
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

		h.ServeHTTP(w, r)
	})
}

// auth middleware
// checks for authenticated user and the path names
// i.e: authenticated /login should redirect to /
func authMw(h http.Handler) http.Handler {
	rx := regexp.MustCompile(`(?i)^/(signup|login|verify)`)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := context.Get(r, "session").(interface {
			GetUser(r *http.Request) map[string]string
		})

		user := sess.GetUser(r)

		if r.URL.Path != "/" && rx.MatchString(r.URL.Path) && user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}
