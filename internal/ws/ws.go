package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/delivery"
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
