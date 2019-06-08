package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"bishack.dev/utils/session"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
)

const (
	oauthEndpoint = "https://github.com/login/oauth"
	userEndpoint  = "https://api.github.com/user"
)

// FinishSignup ...
func FinishSignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.Form.Get("email")
	locale := r.Form.Get("locale")
	picture := r.Form.Get("picture")
	website := r.Form.Get("website")
	nickname := r.Form.Get("login")
	password := r.Form.Get("password")

	u := context.Get(r, "userService").(interface {
		Signup(username, password string, meta map[string]string) (*cip.SignUpOutput, error)
	})

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
			errMessage = "Account already exists. You can log in if you want to"
		}

		sess := context.Get(r, "session").(interface {
			SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
		})
		sess.SetFlash(w, r, "error", errMessage)
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	http.Redirect(w, r, "/verify?email="+email, http.StatusSeeOther)
}

// Verify ...
func Verify(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	email := r.URL.Query().Get("email")

	sess := context.Get(r, "session").(interface {
		SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	u := context.Get(r, "userService").(interface {
		Verify(username, code string) (*cip.ConfirmSignUpOutput, error)
	})

	if code != "" {
		_, err := u.Verify(email, code)

		if err != nil {
			sess.SetFlash(w, r, "error", "Verification failed. Try again!")
			http.Redirect(w, r, "/verify?email="+email, http.StatusSeeOther)
			return
		}

		sess.SetFlash(w, r, "success", "Account Verified!")
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
		return
	}

	flash := sess.GetFlash(w, r)

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
	sess := context.Get(r, "session").(interface {
		SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	client := context.Get(r, "client").(interface {
		PostForm(url string, data url.Values) (*http.Response, error)
	})

	// check for oauth code from github
	if code != "" {
		resp, err := client.PostForm(githubEndpoint(code), url.Values{})
		if err != nil {
			sess.SetFlash(w, r, "error", "Invalid or expired code")
			http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
			return
		}

		b, _ := ioutil.ReadAll(resp.Body)
		val, _ := url.ParseQuery(string(b))

		token := val.Get("access_token")
		if token == "" {
			sess.SetFlash(w, r, "error", "Invalid or expired code")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/signup?access_token="+val.Get("access_token"), http.StatusSeeOther)
		return
	}

	// check for access token after code verification
	accessToken := r.URL.Query().Get("access_token")
	if accessToken != "" {
		client := context.Get(r, "client").(interface {
			Do(r *http.Request) (*http.Response, error)
		})

		req, _ := http.NewRequest(http.MethodGet, userEndpoint, strings.NewReader(""))
		req.Header.Set("Authorization", "token "+accessToken)

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			// flash me baby!
			sess.SetFlash(w, r, "error", "Invalid or expired token!")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		gu := githubUser{}
		err = json.NewDecoder(resp.Body).Decode(&gu)
		if err != nil {
			sess.SetFlash(w, r, "error", "An error occured!")
			http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
			return
		}

		render(w, "main", "signup-form", map[string]interface{}{
			"Title":          "Complete Signup",
			"GithubEndpoint": githubEndpoint(""),
			"GithubUser":     gu,
			"Flash":          sess.GetFlash(w, r),
			csrf.TemplateTag: csrf.TemplateField(r),
		})
		return
	}

	// otherwise, render signup page
	render(w, "main", "signup", map[string]interface{}{
		"Title":     "Sign Up",
		"Flash":     sess.GetFlash(w, r),
		"GithubURL": githubEndpoint(""),
	})
}

// Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	sess := context.Get(r, "session").(interface {
		DeleteUser(w http.ResponseWriter, r *http.Request)
	})
	sess.DeleteUser(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	sess := context.Get(r, "session").(interface {
		SetUser(w http.ResponseWriter, r *http.Request, email, token string)
		SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
	})

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	u := context.Get(r, "userService").(interface {
		Login(username, password string) (*cip.InitiateAuthOutput, error)
	})

	out, err := u.Login(email, password)

	if err != nil {
		sess.SetFlash(w, r, "error", "Wrong email or password")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	token := out.AuthenticationResult.AccessToken

	sess.SetUser(w, r, email, *token)
	sess.SetFlash(w, r, "success", "Welcome Back!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LoginForm ...
func LoginForm(w http.ResponseWriter, r *http.Request) {
	sess := context.Get(r, "session").(interface {
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	render(w, "main", "login-form", map[string]interface{}{
		"Title":          "User Login",
		"Flash":          sess.GetFlash(w, r),
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
