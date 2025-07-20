package appdelivery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app-delivery/internal/adapter"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app-delivery/internal/domain"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/internal/cookies"
)

type controllers struct {
	db *sql.DB
}

func newControllers(db *sql.DB) *controllers {
	return &controllers{
		db: db,
	}
}
func (c *controllers) AppDelivery(w http.ResponseWriter, r *http.Request) {
	homeTemplate := template.Must(template.ParseFiles("template/app-delivery.gohtml"))

	// homeTemplate.Execute(w, "wss://"+r.Host+"/ws")
	homeTemplate.Execute(w, "wss://hm3c050c-8080.brs.devtunnels.ms/ws")
}

func (c *controllers) LoginAppDelivery(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /LoginAppDelivery")

	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var userData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = json.Unmarshal(jsonData, &userData)
	if err != nil {
		fmt.Println("Error unmarshaling:", err)
		return
	}
	fmt.Println("Decoded User:", userData)

	signInInput := &domain.SignInInput{
		Email:    userData.Email,
		Password: userData.Password,
	}

	driverRepository := adapter.NewDriverRepositorySqlite(c.db)
	signIn := domain.NewSignIn(driverRepository)

	output, err := signIn.Execute(signInInput)
	if err != nil {
		respondWithError(w, "Erro ao efetuar login: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Output:", output)

	driverId, err := strconv.Atoi(output.UserId)
	if err != nil {
		respondWithError(w, "Erro ao efetuar login: "+err.Error(), http.StatusBadRequest)
		return
	}

	deliveryRepository := adapter.NewDeliveryRepositorySqlite(c.db)
	listDelivery := domain.NewListDelivery(deliveryRepository, driverRepository)

	deliveries, err := listDelivery.Execute(driverId)
	if err != nil {
		respondWithError(w, "Erro ao efetuar login: "+err.Error(), http.StatusBadRequest)
		return
	}

	value := map[string]string{
		"UserId":   output.UserId,
		"UserType": output.UserType,
	}

	if encoded, err := cookies.S.Encode("userCookieDriver", value); err == nil {
		cookie := &http.Cookie{
			Name:     "userCookieDriver",
			Value:    encoded,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deliveries)
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response := struct {
		Message    string `json:"message"`
		StatusCode int    `json:"statusCode"`
	}{
		Message:    message,
		StatusCode: statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
