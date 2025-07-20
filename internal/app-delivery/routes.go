package appdelivery

import (
	"database/sql"
	"net/http"
)

func Register(mux *http.ServeMux, db *sql.DB) {
	controllers := newControllers(db)

	mux.HandleFunc("GET /app-delivery", controllers.AppDelivery)
	mux.HandleFunc("POST /app-delivery/login", controllers.LoginAppDelivery)
}
