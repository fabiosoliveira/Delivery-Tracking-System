package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/application/company"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/database"
)

func main() {

	db := database.DB
	companyRepository := database.NewCompanyRepositorySqlite(db)
	cadastrarCompany := company.NewCadastrarCompany(companyRepository)

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

		company := company.NewCompany(name, email, password)

		err = cadastrarCompany.Execute(company)
		if err != nil {
			TrowError(err, w, r)
			return
		}

		http.Redirect(w, r, "/auth/signup", http.StatusSeeOther)
	})

	mux.HandleFunc("POST /conta/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.PathValue("_method"))
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(id)

		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		fmt.Println(name, email, password)

		http.Redirect(w, r, "/auth/signup", http.StatusSeeOther)
	})

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
