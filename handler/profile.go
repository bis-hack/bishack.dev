package handler

import (
	"net/http"

	"bishack.dev/utils"
	"bishack.dev/utils/session"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
)

// ProfileForm ...
func ProfileForm(w http.ResponseWriter, r *http.Request) {

	sess := context.Get(r, "session").(interface {
		GetFlash(w http.ResponseWriter, r *http.Request) *session.Flash
	})

	// get user details from context and cast it as map[string]string if
	// not nil
	user := context.Get(r, "user")

	if user != nil {
		user = user.(map[string]string)
	}

	utils.Render(w, "main", "profile-form", map[string]interface{}{
		"Title":          "Edit User Profile",
		"Flash":          sess.GetFlash(w, r),
		"User":           user,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// UpdateProfile ...
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	sess := context.Get(r, "session").(interface {
		SetUser(w http.ResponseWriter, r *http.Request, email, token string)
		SetFlash(w http.ResponseWriter, r *http.Request, t, v string)
	})

	// get user details from context and cast it as map[string]string if
	// not nil
	user := context.Get(r, "user")

	if user != nil {
		user = user.(map[string]string)
	}

	if err := r.ParseForm(); err != nil {
		return
	}

	email := r.Form.Get("email")

	oldPassword := r.Form.Get("old_password")
	newPassword := r.Form.Get("new_password")
	confirmNewPassword := r.Form.Get("confirm_new_password")

	if email == "" {
		sess.SetFlash(w, r, "error", "Email is required")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if newPassword != confirmNewPassword {
		sess.SetFlash(w, r, "error", "Passwords dont match")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if newPassword == oldPassword {
		sess.SetFlash(w, r, "error", "Please dont use the same password")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// TODO:: call API to update user
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
