package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/securecookie"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/auth"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/delivery"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/driver"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/infra/database"
)

// Hash keys should be at least 32 bytes long
var hashKey = []byte("12345678901234567890123456789012")

// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var blockKey = []byte("1234567890123456")
var s = securecookie.New(hashKey, blockKey)

/**
Fábio
fabio@gmail.com
123456As#
**/

func main() {

	db := database.DB
	userRepository := database.NewUserRepositorySqlite(db)
	signUp := auth.NewSignUp(userRepository)
	signIn := auth.NewSignIn(userRepository)
	listDrivers := driver.NewListDrivers(userRepository)
	registerDriver := driver.NewRegister(userRepository)

	deliveryRepository := database.NewDeliveryRepositorySqlite(db)
	createDelivery := delivery.NewCreateDelivery(deliveryRepository, userRepository)
	listDelivery := delivery.NewListDelivery(deliveryRepository, userRepository)

	// companyDao := database.NewCompanyDataAccessObject(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth/signup", func(w http.ResponseWriter, r *http.Request) {

		tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/signup.gohtml"))

		tpl.ExecuteTemplate(w, "layout", nil)

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

		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /auth/signin", func(w http.ResponseWriter, r *http.Request) {

		tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/signin.gohtml"))

		tpl.ExecuteTemplate(w, "layout", nil)

	})

	mux.HandleFunc("POST /auth/signin", func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST /auth/signin")

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

		value := map[string]string{
			"UserId":   *output.UserId,
			"UserType": *output.UserType,
		}

		if encoded, err := s.Encode("userCookie", value); err == nil {
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
	})

	mux.HandleFunc("GET /drivers", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /drivers")
		cookie, err := r.Cookie("userCookie")
		if err != nil {
			log.Println("GET /drivers: cookie not found")
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		value := make(map[string]string)

		err = s.Decode("userCookie", cookie.Value, &value)
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

		err = s.Decode("userCookie", cookie.Value, &value)
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

		err = s.Decode("userCookie", cookie.Value, &value)
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

	mux.HandleFunc("POST /delivery", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("userCookie")
		if err != nil {
			log.Println("Error retrieving cookie:", err)
			http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
			return
		}

		value := make(map[string]string)

		err = s.Decode("userCookie", cookie.Value, &value)
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
