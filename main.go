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

func main() {
	csrfSecure := false
	if regexp.MustCompile("`(?i)stag|prod").MatchString(os.Getenv("UP_STAGE")) {
		csrfSecure = true
	}

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
	r.Get("/signup", handler.Signup)
	r.Get("/verify", handler.Verify)
	r.Post("/signup", handler.FinishSignup)

	// launch
	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(
		port,
		csrf.Protect([]byte(os.Getenv("CSRF_KEY")), csrf.Secure(csrfSecure))(r),
	))
}
