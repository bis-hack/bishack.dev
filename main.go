package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"bishack.dev/handler"
	"github.com/gorilla/csrf"
	"github.com/gorilla/pat"
)

// PUBLICFOLDER ...
const PUBLICFOLDER = "public"

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
		protect(r),
	))
}
