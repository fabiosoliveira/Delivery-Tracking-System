package driver

import (
	"net/http"
)

func registerRoutes(mux *http.ServeMux, controllers *controllers) {

	mux.HandleFunc("GET /drivers", controllers.getDrivers)

	mux.HandleFunc("POST /drivers", controllers.postDrivers)
}
