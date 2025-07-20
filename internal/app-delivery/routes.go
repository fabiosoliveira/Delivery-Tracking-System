package appdelivery

import "net/http"

func registerRoutes(mux *http.ServeMux, controllers *controllers) {

	mux.HandleFunc("GET /app-delivery", controllers.AppDelivery)
	mux.HandleFunc("POST /app-delivery/login", controllers.LoginAppDelivery)
}
