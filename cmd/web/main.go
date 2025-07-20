package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/auth"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/delivery"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/driver"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/ws"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/infra/database"
)

/**
Fábio
fabio@gmail.com
123456As#
**/

func main() {

	db := database.DB

	mux := http.NewServeMux()

	auth.Module(mux, db)
	delivery.Module(mux, db)
	driver.Module(mux, db)

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
