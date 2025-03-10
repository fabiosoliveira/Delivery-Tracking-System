package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/auth"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/cookies"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/infra/database"
)

// The secret key used for AES-GCM encryption must be exactly 32 bytes long.
var secretKey = []byte("12345678901234567890123456789012")

func main() {

	db := database.DB
	userRepository := database.NewUserRepositorySqlite(db)
	signUp := auth.NewSignUp(userRepository)
	signIn := auth.NewSignIn(userRepository)

	// companyDao := database.NewCompanyDataAccessObject(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth/signup", func(w http.ResponseWriter, r *http.Request) {

		tpl := template.Must(template.ParseFiles("template/signup.gohtml"))

		tpl.Execute(w, nil)

	})

	mux.HandleFunc("POST /auth/signup", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			TrowError(err, w, r)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		signUpInput := &auth.SignUpInput{
			Name:     name,
			Email:    email,
			Password: password,
		}

		err = signUp.Execute(signUpInput)
		if err != nil {
			TrowError(err, w, r)
			return
		}

		http.Redirect(w, r, "/auth/signup", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /auth/signin", func(w http.ResponseWriter, r *http.Request) {

		tpl := template.Must(template.ParseFiles("template/signin.gohtml"))

		tpl.Execute(w, nil)

	})

	mux.HandleFunc("POST /auth/signin", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			TrowError(err, w, r)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		signInInput := &auth.SignInInput{
			Email:    email,
			Password: password,
		}

		output, err := signIn.Execute(signInInput)
		if err != nil {
			TrowError(err, w, r)
			return
		}

		b, err := json.Marshal(output)
		if err != nil {
			TrowError(err, w, r)
			return
		}

		cookie := &http.Cookie{
			Name:     "userCookie",
			Value:    string(b),
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		if err := cookies.WriteEncrypted(w, cookie, secretKey); err != nil {
			TrowError(err, w, r)
			return
		}

		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
	})

	// mux.HandleFunc("POST /conta/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println(r.PathValue("_method"))
	// 	id, err := strconv.Atoi(r.PathValue("id"))
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}

	// 	fmt.Println(id)

	// 	err = r.ParseForm()
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}

	// 	name := r.FormValue("name")
	// 	email := r.FormValue("email")
	// 	password := r.FormValue("password")

	// 	fmt.Println(name, email, password)

	// 	http.Redirect(w, r, "/auth/signup", http.StatusSeeOther)
	// })

	log.Println("Servidor Http iniciado na porta 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}

func TrowError(err error, w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("template/error.gohtml"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	tpl.Execute(w, struct {
		Error string
		Url   string
	}{err.Error(), r.URL.String()})
}
