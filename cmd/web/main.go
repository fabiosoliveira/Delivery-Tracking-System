package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	appDelivery "github.com/fabiosoliveira/Delivery-Tracking-System/internal/app-delivery"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/auth"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/delivery"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/driver"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/infra/database"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/ws"
)

/**
Fábio
fabio@gmail.com
123456As#
**/

func main() {

	db := database.DB

	mux := http.NewServeMux()

	auth.Register(mux, db)
	delivery.Register(mux, db)
	driver.Register(mux, db)
	appDelivery.Register(mux, db)
	ws.Register(mux, db)

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
