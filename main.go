package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"

	"bishack.dev/handler"
	mw "bishack.dev/middleware"

	// autoload env
	"github.com/aws/aws-xray-sdk-go/xray"
	_ "github.com/joho/godotenv/autoload"

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
	r.Post("/signup", handler.FinishSignup)
	r.Post("/login", handler.Login)

	// profile
	r.Get("/profile", handler.Profile)
	r.Post("/profile", handler.UpdateProfile)

	// security
	r.Get("/security", handler.Security)
	r.Post("/security", handler.ChangePassword)

	// like
	r.Put("/like/{id}", handler.ToggleLike)

	// slack
	r.Get("/slack-invite", handler.SlackInvite)

	// post
	r.Post("/update-post", handler.UpdatePost)
	r.Get("/edit/{id}", handler.EditPost)
	r.Get("/new", handler.New)
	r.Post("/new", handler.CreatePost)
	r.Get("/{username}/{id}", handler.GetPost)
	r.Get("/{username}", handler.GetUserPosts)

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
		xray.Handler(
			xray.NewFixedSegmentNamer("bishack.dev"),
			protect(mw.Context(mw.Token(mw.SessionUser(mw.AuthRedirects(r))))),
		),
	))
}
