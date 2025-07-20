package auth

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/auth/internal/adapter"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/auth/internal/domain"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/internal/utils"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/cookies"
)

type controllers struct {
	db *sql.DB
}

func newControllers(db *sql.DB) *controllers {
	return &controllers{
		db: db,
	}
}

func (c *controllers) getSignup(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/signup.gohtml"))

	tpl.ExecuteTemplate(w, "layout", nil)
}

func (c *controllers) postSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.TrowError(err, w, r)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	signUpInput := &domain.SignUpInput{
		Name:     name,
		Email:    email,
		Password: password,
	}

	signUp := domain.NewSignUp(adapter.NewUserRepositorySqlite(c.db))
	err = signUp.Execute(signUpInput)
	if err != nil {
		utils.TrowError(err, w, r)
		return
	}

	http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
}

func (c *controllers) getSignin(w http.ResponseWriter, r *http.Request) {

	tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/signin.gohtml"))

	tpl.ExecuteTemplate(w, "layout", nil)

}

func (c *controllers) postSignin(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /auth/signin")

	err := r.ParseForm()
	if err != nil {
		utils.TrowError(err, w, r)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	signInInput := &domain.SignInInput{
		Email:    email,
		Password: password,
	}

	signIn := domain.NewSignIn(adapter.NewUserRepositorySqlite(c.db))
	output, err := signIn.Execute(signInInput)
	if err != nil {
		utils.TrowError(err, w, r)
		return
	}

	value := map[string]string{
		"UserId":   output.UserId,
		"UserType": output.UserType,
	}

	if encoded, err := cookies.S.Encode("userCookie", value); err == nil {
		cookie := &http.Cookie{
			Name:     "userCookie",
			Value:    encoded,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
	}

	log.Println("Redirect to /drivers")
	http.Redirect(w, r, "/drivers", http.StatusSeeOther)
}
