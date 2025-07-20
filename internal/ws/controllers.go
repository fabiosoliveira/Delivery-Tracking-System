package ws

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/ws/internal/adapter"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/ws/internal/domain"
	"github.com/gorilla/websocket"
)

type controllers struct {
	db *sql.DB
}

func newControllers(db *sql.DB) *controllers {
	return &controllers{
		db: db,
	}
}

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

func (c *controllers) HandleConnection(w http.ResponseWriter, r *http.Request) {

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

		deliveryDao := adapter.NewDeliveryDAOSqlite(c.db)
		sendLocation := domain.NewSendLocation(deliveryDao)

		err = sendLocation.Execute(&domain.SendLocationInput{
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

func (c *controllers) HandleClientConnection(w http.ResponseWriter, r *http.Request) {
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
