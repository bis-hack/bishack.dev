package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"bishack.dev/utils"
	"bishack.dev/utils/session"
	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
)

const (
	userEndpoint = "https://api.github.com/user"
)

// FinishSignup ...
func FinishSignup(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	name := r.Form.Get("name")
	email := r.Form.Get("email")
	locale := r.Form.Get("locale")
	profile := r.Form.Get("profile")
	picture := r.Form.Get("picture")
	website := r.Form.Get("website")
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	u := context.Get(r, "userService").(interface {
		Signup(
			username,
			password string,
			meta map[string]string,
		) (*cip.SignUpOutput, error)
	})

	meta := map[string]string{
		"name":     name,
		"email":    email,
		"locale":   locale,
		"profile":  profile,
		"website":  website,
		"picture":  picture,
		"nickname": username,
	}

	_, err := u.Signup(username, password, meta)
	if err != nil {
		errMessage := "Could not sign you up. Try again!"

		if regexp.MustCompile("exists").MatchString(err.Error()) {
			errMessage = "Account already exists. Try to log in instead."
		}

		sess := context.Get(r, "session").(interface {
			SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
		})
		sess.SetFlash(w, r, "error", errMessage)
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	http.Redirect(w, r, "/verify?username="+username, http.StatusSeeOther)
}

// Verify ...
func Verify(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	username := r.URL.Query().Get("username")

	sess := context.Get(r, "session").(interface {
		SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	u := context.Get(r, "userService").(interface {
		Verify(username, code string) (*cip.ConfirmSignUpOutput, error)
	})

	if code != "" {
		_, err := u.Verify(username, code)

		if err != nil {
			sess.SetFlash(w, r, "error", "Verification failed. Try again!")
			http.Redirect(w, r, "/verify?username="+username, http.StatusSeeOther)
			return
		}

		sess.SetFlash(w, r, "success", "Account Verified!")
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
		return
	}

	flash := sess.GetFlash(w, r)

	// horray!
	if flash != nil && flash.Type == "success" {
		utils.Render(w, "main", "verified", map[string]interface{}{
			"Title": "Account Verified",
			"Flash": flash,
		})
		return
	}

	utils.Render(w, "main", "verify-form", map[string]interface{}{
		"Title":          "Verify",
		"Username":       username,
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
		resp, err := client.PostForm(utils.GithubEndpoint(code), url.Values{})
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

		http.Redirect(
			w,
			r,
			"/signup?access_token="+val.Get("access_token"),
			http.StatusSeeOther,
		)

		return
	}

	// check for access token after code verification
	accessToken := r.URL.Query().Get("access_token")
	if accessToken != "" {
		client := context.Get(r, "client").(interface {
			Do(r *http.Request) (*http.Response, error)
		})

		req, _ := http.NewRequest(
			http.MethodGet,
			userEndpoint,
			strings.NewReader(""),
		)
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
			sess.SetFlash(w, r, "error", "An error occurred!")
			http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
			return
		}

		utils.Render(w, "main", "signup-form", map[string]interface{}{
			"Title":          "Complete Signup",
			"GithubEndpoint": utils.GithubEndpoint(""),
			"GithubUser":     gu,
			"Flash":          sess.GetFlash(w, r),
			csrf.TemplateTag: csrf.TemplateField(r),
		})
		return
	}

	// otherwise, utils.Render signup page
	utils.Render(w, "main", "signup", map[string]interface{}{
		"Title":     "Sign Up",
		"Flash":     sess.GetFlash(w, r),
		"GithubURL": utils.GithubEndpoint(""),
	})
}

// Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	sess := context.Get(r, "session").(interface {
		DeleteUser(w http.ResponseWriter, r *http.Request)
	})
	sess.DeleteUser(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	sess := context.Get(r, "session").(interface {
		SetUser(w http.ResponseWriter, r *http.Request, username, token string)
		SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
	})

	_ = r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	u := context.Get(r, "userService").(interface {
		Login(username, password string) (*cip.InitiateAuthOutput, error)
	})

	out, err := u.Login(username, password)

	if err != nil {
		sess.SetFlash(w, r, "error", "Wrong username or password")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	token := out.AuthenticationResult.RefreshToken

	sess.SetUser(w, r, username, *token)
	sess.SetFlash(w, r, "success", "Welcome Back!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LoginForm ...
func LoginForm(w http.ResponseWriter, r *http.Request) {
	sess := context.Get(r, "session").(interface {
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	utils.Render(w, "main", "login-form", map[string]interface{}{
		"Title":          "User Login",
		"Flash":          sess.GetFlash(w, r),
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
