package driver

import (
	"database/sql"
	"net/http"
)

func Module(mux *http.ServeMux, db *sql.DB) {
	controllers := newControllers(db)
	registerRoutes(mux, controllers)
}
