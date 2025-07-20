package ws

import (
	"database/sql"
	"net/http"
)

func Register(mux *http.ServeMux, db *sql.DB) {
	controllers := newControllers(db)

	mux.HandleFunc("/ws", controllers.HandleConnection)
	mux.HandleFunc("/ws/", controllers.HandleClientConnection)
}
