package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/auth"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/delivery"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/driver"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/ws"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/cookies"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/infra/database"
)

/**
Fábio
fabio@gmail.com
123456As#
**/

func main() {

	db := database.DB
	userRepository := database.NewUserRepositorySqlite(db)
	// signUp := auth.NewSignUp(userRepository)
	// signIn := auth.NewSignIn(userRepository)
	listDrivers := driver.NewListDrivers(userRepository)
	registerDriver := driver.NewRegister(userRepository)

	// deliveryRepository := database.NewDeliveryRepositorySqlite(db)
	// createDelivery := delivery.NewCreateDelivery(deliveryRepository, userRepository)
	// listDelivery := delivery.NewListDelivery(deliveryRepository, userRepository)
	// getDeliveryHistory := delivery.NewGetDeliveryHistory(deliveryRepository)

	// companyDao := database.NewCompanyDataAccessObject(db)

	mux := http.NewServeMux()

	auth.Module(mux, db)
	delivery.Module(mux, db)

	mux.HandleFunc("GET /drivers", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /drivers")
		cookie, err := r.Cookie("userCookie")
		if err != nil {
			log.Println("GET /drivers: cookie not found")
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		value := make(map[string]string)

		err = cookies.S.Decode("userCookie", cookie.Value, &value)
		if err != nil {
			log.Println("GET /drivers: error decoding cookie", err)
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		id, err := strconv.Atoi(value["UserId"])
		if err != nil {
			log.Println("GET /drivers: error converting UserId to int", err)
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		drivers, err := listDrivers.Execute(id)
		if err != nil {
			TrowError(err, w, r)
			return

		}
		tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/drivers.gohtml", "template/drivers-row.gohtml"))

		tpl.ExecuteTemplate(w, "layout", drivers)
	})

	mux.HandleFunc("POST /drivers", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("userCookie")
		if err != nil {
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		value := make(map[string]string)

		err = cookies.S.Decode("userCookie", cookie.Value, &value)
		if err != nil {
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		id, err := strconv.Atoi(value["UserId"])
		if err != nil {
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		register := &driver.RegisterInput{
			Name:      name,
			Email:     email,
			Password:  password,
			CompanyId: uint(id),
		}

		err = registerDriver.Execute(register)
		if err != nil {
			TrowError(err, w, r)
			return
		}

		http.Redirect(w, r, "/drivers", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /app-delivery", ws.AppDelivery)
	mux.HandleFunc("POST /app-delivery/login", ws.LoginAppDelivery)
	mux.HandleFunc("/ws", ws.HandleConnection)
	mux.HandleFunc("/ws/", ws.HandleClientConnection)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Iniciar servidor em goroutine
	go func() {
		fmt.Println("Servidor iniciado em http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Erro ao iniciar o servidor:", err)
		}
	}()

	<-stop
	fmt.Println("\nEncerrando o servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Encerrar servidor HTTP
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Erro durante o shutdown do servidor:", err)
	}

	// Encerrar conexão com o banco de dados
	if err := db.Close(); err != nil {
		fmt.Println("Erro ao fechar conexão com banco de dados:", err)
	} else {
		fmt.Println("Conexão com banco de dados encerrada com sucesso.")
	}

	fmt.Println("Servidor finalizado com sucesso.")
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
