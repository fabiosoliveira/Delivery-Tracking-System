package delivery

import "net/http"

func registerRoutes(mux *http.ServeMux, controllers *controllers) {

	mux.HandleFunc("GET /delivery", controllers.getDelivery)

	mux.HandleFunc("GET /delivery/{id}/localization", controllers.getDeliveryLocalization)

	mux.HandleFunc("POST /delivery", controllers.postDelivery)

	mux.HandleFunc("GET /api/delivery/{id}/history", controllers.getDeliveryHistory)
}
