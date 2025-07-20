package main

import (
	"context"
	"encoding/json"
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

	deliveryRepository := database.NewDeliveryRepositorySqlite(db)
	createDelivery := delivery.NewCreateDelivery(deliveryRepository, userRepository)
	listDelivery := delivery.NewListDelivery(deliveryRepository, userRepository)
	getDeliveryHistory := delivery.NewGetDeliveryHistory(deliveryRepository)

	// companyDao := database.NewCompanyDataAccessObject(db)

	mux := http.NewServeMux()

	auth.Module(mux, db)

	// mux.HandleFunc("GET /auth/signup", func(w http.ResponseWriter, r *http.Request) {

	// 	tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/signup.gohtml"))

	// 	tpl.ExecuteTemplate(w, "layout", nil)

	// })

	// mux.HandleFunc("POST /auth/signup", func(w http.ResponseWriter, r *http.Request) {
	// 	err := r.ParseForm()
	// 	if err != nil {
	// 		TrowError(err, w, r)
	// 		return
	// 	}

	// 	name := r.FormValue("name")
	// 	email := r.FormValue("email")
	// 	password := r.FormValue("password")

	// 	signUpInput := &auth.SignUpInput{
	// 		Name:     name,
	// 		Email:    email,
	// 		Password: password,
	// 	}

	// 	err = signUp.Execute(signUpInput)
	// 	if err != nil {
	// 		TrowError(err, w, r)
	// 		return
	// 	}

	// 	http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
	// })

	// mux.HandleFunc("GET /auth/signin", func(w http.ResponseWriter, r *http.Request) {

	// 	tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/signin.gohtml"))

	// 	tpl.ExecuteTemplate(w, "layout", nil)

	// })

	// mux.HandleFunc("POST /auth/signin", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("POST /auth/signin")

	// 	err := r.ParseForm()
	// 	if err != nil {
	// 		TrowError(err, w, r)
	// 		return
	// 	}

	// 	email := r.FormValue("email")
	// 	password := r.FormValue("password")

	// 	signInInput := &auth.SignInInput{
	// 		Email:    email,
	// 		Password: password,
	// 	}

	// 	output, err := signIn.Execute(signInInput)
	// 	if err != nil {
	// 		TrowError(err, w, r)
	// 		return
	// 	}

	// 	value := map[string]string{
	// 		"UserId":   *output.UserId,
	// 		"UserType": *output.UserType,
	// 	}

	// 	if encoded, err := cookies.S.Encode("userCookie", value); err == nil {
	// 		cookie := &http.Cookie{
	// 			Name:     "userCookie",
	// 			Value:    encoded,
	// 			Path:     "/",
	// 			MaxAge:   3600,
	// 			HttpOnly: true,
	// 			Secure:   true,
	// 			SameSite: http.SameSiteLaxMode,
	// 		}
	// 		http.SetCookie(w, cookie)
	// 	}

	// 	log.Println("Redirect to /drivers")
	// 	http.Redirect(w, r, "/drivers", http.StatusSeeOther)
	// })

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

	mux.HandleFunc("GET /delivery", func(w http.ResponseWriter, r *http.Request) {
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

		drivers, err := listDrivers.Execute(id)
		if err != nil {
			TrowError(err, w, r)
			return
		}

		deliveries, err := listDelivery.Execute(id)
		if err != nil {
			TrowError(err, w, r)
			return
		}

		data := struct {
			Drivers    []driver.ListDriversOutput
			Deliveries []delivery.ListDeliveryOutput
		}{
			Drivers:    drivers,
			Deliveries: deliveries,
		}

		tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/delivery.gohtml", "template/delivery-row.gohtml"))

		tpl.ExecuteTemplate(w, "layout", data)
	})

	mux.HandleFunc("GET /delivery/{id}/localization", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /delivery/{id}/localization")

		tpl := template.Must(template.ParseFiles("template/map-localization.gohtml"))

		tpl.Execute(w, nil)
	})

	mux.HandleFunc("POST /delivery", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("userCookie")
		if err != nil {
			log.Println("Error retrieving cookie:", err)
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		value := make(map[string]string)

		err = cookies.S.Decode("userCookie", cookie.Value, &value)
		if err != nil {
			log.Println("Error decoding cookie:", err)
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		companyId, err := strconv.Atoi(value["UserId"])
		if err != nil {
			log.Println("Error converting UserId to int:", err)
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		err = r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recipient := r.FormValue("recipient")
		address := r.FormValue("address")
		driver_id := r.FormValue("driver_id")

		driverId, err := strconv.Atoi(driver_id)
		if err != nil {
			log.Println("Error converting DriverId to int:", err)
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		log.Println("Received data - Driver ID:", driverId, "Recipient:", recipient, "Address:", address)

		createDeliveryInput := &delivery.CreateDeliveryInput{
			CompanyId: uint(companyId),
			DriverId:  uint(driverId),
			Recipient: recipient,
			Address:   address,
		}

		log.Println("createDeliveryInput:", createDeliveryInput)

		err = createDelivery.Execute(createDeliveryInput)
		if err != nil {
			log.Println("Error executing CreateDelivery:", err)
			TrowError(err, w, r)
			return
		}

		http.Redirect(w, r, "/delivery", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /app-delivery", ws.AppDelivery)
	mux.HandleFunc("POST /app-delivery/login", ws.LoginAppDelivery)
	mux.HandleFunc("/ws", ws.HandleConnection)
	mux.HandleFunc("/ws/", ws.HandleClientConnection)

	mux.HandleFunc("GET /api/delivery/{id}/history", func(w http.ResponseWriter, r *http.Request) {
		deliveryId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
			return
		}

		history, err := getDeliveryHistory.Execute(deliveryId)
		if err != nil {
			http.Error(w, "Error fetching delivery history", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(history)
	})

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
