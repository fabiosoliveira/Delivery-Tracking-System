package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/application/user"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/database"
)

func main() {

	db := database.DB
	userRepository := database.NewUserRepositorySqlite(db)
	criarConta := user.NewCriarConta(userRepository)

	userDao := database.NewUserDataAccessObject(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /conta", func(w http.ResponseWriter, r *http.Request) {
		users, err := userDao.ListarContas()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tpl := template.Must(template.ParseFiles("template/conta.gohtml", "template/conta-row.gohtml"))

		tpl.Execute(w, users)

	})

	mux.HandleFunc("POST /conta", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		userType := r.FormValue("type")

		userTypeID, err := strconv.Atoi(userType)
		if err != nil {
			http.Error(w, "Invalid user type", http.StatusBadRequest)
			return
		}

		user := user.NewUser(name, email, password, user.UserType(userTypeID))

		criarConta.Execute(user)

		http.Redirect(w, r, "/conta", http.StatusSeeOther)
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
		userType := r.FormValue("type")

		userTypeID, err := strconv.Atoi(userType)
		if err != nil {
			http.Error(w, "Invalid user type", http.StatusBadRequest)
			return
		}

		fmt.Println(name, email, password, user.UserType(userTypeID))

		http.Redirect(w, r, "/conta", http.StatusSeeOther)
	})

	log.Println("Servidor Http iniciado na porta 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
