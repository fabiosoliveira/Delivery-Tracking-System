package driver

import (
	"database/sql"
	"net/http"
)

func Register(mux *http.ServeMux, db *sql.DB) {
	controllers := newControllers(db)

	mux.HandleFunc("GET /drivers", controllers.getDrivers)

	mux.HandleFunc("POST /drivers", controllers.postDrivers)
}
