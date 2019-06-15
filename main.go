package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"

	"bishack.dev/handler"
	mw "bishack.dev/middleware"

	"github.com/gorilla/csrf"
	"github.com/gorilla/pat"
)

func main() {
	// csrf
	csrfSecure := true

	// env
	rxEnv := regexp.MustCompile("`(?i)stag|prod")
	isLive := rxEnv.MatchString(os.Getenv("UP_STAGE"))

	// init route
	r := pat.New()

	// not found
	r.NotFoundHandler = http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.URL.Path == "/" {
			handler.Home(w, r)
		} else {
			handler.NotFound(w, r)
		}
	})

	// handlers et al

	// auth
	r.Get("/signup", handler.Signup)
	r.Get("/verify", handler.Verify)
	r.Get("/login", handler.LoginForm)
	r.Get("/logout", handler.Logout)
	r.Get("/slack-invite", handler.SlackInvite)
	r.Get("/profile", handler.ProfileForm)

	// POST
	r.Post("/signup", handler.FinishSignup)
	r.Post("/login", handler.Login)
	r.Post("/update", handler.UpdateProfile)

	// slack
	r.Get("/slack-invite", handler.SlackInvite)

	// post
	r.Get("/{username}/{id}", handler.GetPost)
	r.Get("/new", handler.New)
	r.Post("/new", handler.CreatePost)

	// user
	r.Get("/{username}", handler.GetUserPosts)

	// like
	r.Put("/like/{id}", handler.ToggleLike)

	// on local
	if !isLive {
		// set secure to false
		csrfSecure = false

		// launch nerdy stuff(pprof) server
		go func() {
			_ = http.ListenAndServe(":6060", nil)
		}()
	}

	protect := csrf.Protect([]byte(
		os.Getenv("CSRF_KEY")),
		csrf.Secure(csrfSecure),
	)

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(
		port,
		protect(mw.Context(mw.SessionUser(mw.AuthRedirects(r)))),
	))
}
