package delivery

import (
	"database/sql"
	"net/http"
)

func Register(mux *http.ServeMux, db *sql.DB) {
	controllers := newControllers(db)

	mux.HandleFunc("GET /delivery", controllers.getDelivery)

	mux.HandleFunc("GET /delivery/{id}/localization", controllers.getDeliveryLocalization)

	mux.HandleFunc("POST /delivery", controllers.postDelivery)

	mux.HandleFunc("GET /api/delivery/{id}/history", controllers.getDeliveryHistory)
}
