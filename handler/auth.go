package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"bishack.dev/services/user"
	"bishack.dev/utils/session"
	"github.com/gorilla/csrf"
)

const (
	oauthEndpoint = "https://github.com/login/oauth"
	userEndpoint  = "https://api.github.com/user"
)

// FinishSignup ...
func FinishSignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := os.Getenv("COGNITO_CLIENT_ID")
	secret := os.Getenv("COGNITO_CLIENT_SECRET")

	email := r.Form.Get("email")
	locale := r.Form.Get("locale")
	picture := r.Form.Get("picture")
	website := r.Form.Get("website")
	nickname := r.Form.Get("login")
	password := r.Form.Get("password")

	u := user.New(id, secret)

	meta := map[string]string{
		"email":    email,
		"locale":   locale,
		"website":  website,
		"picture":  picture,
		"nickname": nickname,
	}

	_, err := u.Signup(email, password, meta)
	if err != nil {
		errMessage := "Could not sign you up. Try again!"

		if regexp.MustCompile("exists").MatchString(err.Error()) {
			errMessage = "Account exists already"
		}

		session.SetFlash(w, r, "error", errMessage)
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	http.Redirect(w, r, "/verify?email="+email, http.StatusSeeOther)
}

// Verify ...
func Verify(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	email := r.URL.Query().Get("email")

	if code != "" {
		id := os.Getenv("COGNITO_CLIENT_ID")
		secret := os.Getenv("COGNITO_CLIENT_SECRET")

		u := user.New(id, secret)
		_, err := u.Verify(email, code)

		if err != nil {
			log.Println(err.Error())
			session.SetFlash(w, r, "error", "Verification failed. Try again!")
			http.Redirect(w, r, "/verify?email="+email, http.StatusSeeOther)
			return
		}

		session.SetFlash(w, r, "success", "Account Verified!")
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
		return
	}

	flash := session.GetFlash(w, r)

	// horray!
	if flash != nil && flash.Type == "success" {
		render(w, "main", "verified", map[string]interface{}{
			"Title": "Account Verified",
			"Flash": flash,
		})
		return
	}

	render(w, "main", "verify-form", map[string]interface{}{
		"Title":          "Verify",
		"Email":          email,
		"Flash":          flash,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// Signup ...
func Signup(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// check for oauth code from github
	if code != "" {
		resp, err := http.PostForm(githubEndpoint(code), url.Values{})
		if err != nil {
			session.SetFlash(w, r, "error", "Invalid or expired code!")
			http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
			return
		}

		b, _ := ioutil.ReadAll(resp.Body)
		val, _ := url.ParseQuery(string(b))

		token := val.Get("access_token")
		if token == "" {
			session.SetFlash(w, r, "error", "Invalid or expired code")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/signup?access_token="+val.Get("access_token"), http.StatusSeeOther)
		return
	}

	// check for access token after code verification
	accessToken := r.URL.Query().Get("access_token")
	if accessToken != "" {
		cl := http.Client{}

		req, _ := http.NewRequest(http.MethodGet, userEndpoint, strings.NewReader(""))
		req.Header.Set("Authorization", "token "+accessToken)

		resp, err := cl.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			// flash me baby!
			session.SetFlash(w, r, "error", "Invalid or expired token!")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		gu := githubUser{}
		err = json.NewDecoder(resp.Body).Decode(&gu)
		if err != nil {
			session.SetFlash(w, r, "error", "An error occured!")
			http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
			return
		}

		render(w, "main", "signup-form", map[string]interface{}{
			"Title":          "Complete Signup",
			"GithubEndpoint": githubEndpoint(""),
			"GithubUser":     gu,
			"Flash":          session.GetFlash(w, r),
			csrf.TemplateTag: csrf.TemplateField(r),
		})
		return
	}

	// otherwise, render signup page
	render(w, "main", "signup", map[string]interface{}{
		"Title":     "Sign Up",
		"Flash":     session.GetFlash(w, r),
		"GithubURL": githubEndpoint(""),
	})
}
