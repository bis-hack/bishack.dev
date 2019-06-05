package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	oauthEndpoint = "https://github.com/login/oauth"
	userEndpoint  = "https://api.github.com/user"
)

// FinishSignup ...
func FinishSignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// log for now
	log.Println(r.Form)
	// redirect lang sa
	http.Redirect(w, r, "/", 301)
}

// Signup ...
func Signup(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code != "" {
		resp, err := http.PostForm(githubEndpoint(code), url.Values{})
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		b, _ := ioutil.ReadAll(resp.Body)
		val, _ := url.ParseQuery(string(b))

		http.Redirect(w, r, "/signup?access_token="+val.Get("access_token"), 301)
		return
	}

	accessToken := r.URL.Query().Get("access_token")
	if accessToken != "" {
		cl := http.Client{}

		req, _ := http.NewRequest(http.MethodGet, userEndpoint, strings.NewReader(""))
		req.Header.Set("Authorization", "token "+accessToken)

		resp, err := cl.Do(req)
		if err != nil || resp.StatusCode != 200 {
			http.Redirect(w, r, "/", 301)
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
			http.Error(w, err.Error(), 400)
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
