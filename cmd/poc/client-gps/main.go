package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type GPSData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
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

		log.Printf("Coordenadas recebidas: Latitude: %f, Longitude: %f\n", gpsData.Latitude, gpsData.Longitude)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate := template.Must(template.ParseFiles("cmd/poc/client-gps/index.gohtml"))

	// homeTemplate.Execute(w, "wss://"+r.Host+"/ws")
	homeTemplate.Execute(w, "wss://l0fmwclb-8080.brs.devtunnels.ms/ws")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/ws", handleConnection)

	log.Println("Servidor WebSocket iniciado na porta 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
