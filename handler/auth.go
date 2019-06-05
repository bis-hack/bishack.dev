package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"bishack.dev/services/user"
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
	password := r.Form.Get("password")

	u := user.New(id, secret)
	_, err := u.Signup(email, password, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		render(w, "main", "verified", struct {
			Title string
		}{"User Verified"})
		return
	}

	render(w, "main", "verify-form", struct {
		Title string
		Email string
	}{"Verify", email})
}

// Signup ...
func Signup(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code != "" {
		resp, err := http.PostForm(githubEndpoint(code), url.Values{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		b, _ := ioutil.ReadAll(resp.Body)
		val, _ := url.ParseQuery(string(b))

		http.Redirect(w, r, "/signup?access_token="+val.Get("access_token"), http.StatusPermanentRedirect)
		return
	}

	accessToken := r.URL.Query().Get("access_token")
	if accessToken != "" {
		cl := http.Client{}

		req, _ := http.NewRequest(http.MethodGet, userEndpoint, strings.NewReader(""))
		req.Header.Set("Authorization", "token "+accessToken)

		resp, err := cl.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}

		type user struct {
			Email     string
			Name      string
			Login     string
			AvatarURL string `json:"avatar_url"`
			Location  string
			Website   string `json:"blog"`
		}

		gu := user{}
		err = json.NewDecoder(resp.Body).Decode(&gu)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		render(w, "main", "signup-form", struct {
			Title      string
			GithubURL  string
			GithubUser user
		}{"Complete Signup", githubEndpoint(""), gu})
		return
	}

	render(w, "main", "signup", struct {
		Title     string
		GithubURL string
	}{"Sign Up", githubEndpoint("")})
}
