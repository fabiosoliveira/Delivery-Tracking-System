package ws

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/auth"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/delivery"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/cookies"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/infra/database"
	"github.com/gorilla/websocket"
)

type GPSData struct {
	Delivery  int64   `json:"delivery_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var (
	clients      = make(map[int64]*websocket.Conn)
	clientsMutex = sync.Mutex{}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	deliveryRepository := database.NewDeliveryRepositorySqlite(database.DB)
	sendLocation := delivery.NewSendLocation(deliveryRepository)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao atualizar conexão:", err)
		return
	}
	defer conn.Close()

	log.Println("Cliente conectado")

	for {
		var gpsData GPSData

		err := conn.ReadJSON(&gpsData)
		if err != nil {
			log.Println("Erro ao ler mensagem:", err)
			break
		}

		log.Printf("Coordenadas recebidas: Delivery: %d, Latitude: %f, Longitude: %f\n", gpsData.Delivery, gpsData.Latitude, gpsData.Longitude)

		clientsMutex.Lock()
		if clientConn, ok := clients[gpsData.Delivery]; ok {
			err := clientConn.WriteJSON(gpsData)
			if err != nil {
				log.Printf("Erro ao enviar dados para o cliente: %v", err)
				clientConn.Close()
				delete(clients, gpsData.Delivery)
			}
		}
		clientsMutex.Unlock()

		err = sendLocation.Execute(&delivery.SendLocationInput{
			Delivery:  gpsData.Delivery,
			Latitude:  gpsData.Latitude,
			Longitude: gpsData.Longitude,
		})
		if err != nil {
			log.Println("Erro ao enviar localização:", err)
			break
		}
	}
}

func HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	deliveryID, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao atualizar conexão do cliente:", err)
		return
	}
	defer conn.Close()

	clientsMutex.Lock()
	clients[deliveryID] = conn
	clientsMutex.Unlock()

	log.Printf("Cliente conectado para delivery %d", deliveryID)

	// Keep connection open to listen for disconnect
	for {
		if _, _, err := conn.NextReader(); err != nil {
			clientsMutex.Lock()
			delete(clients, deliveryID)
			clientsMutex.Unlock()
			log.Printf("Cliente desconectado para delivery %d", deliveryID)
			break
		}
	}
}

func AppDelivery(w http.ResponseWriter, r *http.Request) {
	homeTemplate := template.Must(template.ParseFiles("template/app-delivery.gohtml"))

	// homeTemplate.Execute(w, "wss://"+r.Host+"/ws")
	homeTemplate.Execute(w, "wss://hm3c050c-8080.brs.devtunnels.ms/ws")
}

func LoginAppDelivery(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /LoginAppDelivery")

	userRepository := database.NewUserRepositorySqlite(database.DB)
	signIn := auth.NewSignIn(userRepository)

	deliveryRepository := database.NewDeliveryRepositorySqlite(database.DB)
	listDelivery := NewListDelivery(deliveryRepository, userRepository)

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

	signInInput := &auth.SignInInput{
		Email:    userData.Email,
		Password: userData.Password,
	}

	output, err := signIn.Execute(signInInput)
	if err != nil {
		respondWithError(w, "Erro ao efetuar login: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Output:", output)

	if *output.UserType != "driver" {
		respondWithError(w, "Erro ao efetuar login: usuário não é um motorista", http.StatusBadRequest)
		return
	}

	driverId, err := strconv.Atoi(*output.UserId)
	if err != nil {
		respondWithError(w, "Erro ao efetuar login: "+err.Error(), http.StatusBadRequest)
		return
	}

	deliveries, err := listDelivery.Execute(driverId)
	if err != nil {
		respondWithError(w, "Erro ao efetuar login: "+err.Error(), http.StatusBadRequest)
		return
	}

	value := map[string]string{
		"UserId":   *output.UserId,
		"UserType": *output.UserType,
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
