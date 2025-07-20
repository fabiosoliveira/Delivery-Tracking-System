package auth

import (
	"database/sql"
	"net/http"
)

func Register(mux *http.ServeMux, db *sql.DB) {
	controllers := newControllers(db)

	mux.HandleFunc("GET /auth/signup", controllers.getSignup)
	mux.HandleFunc("POST /auth/signup", controllers.postSignup)
	mux.HandleFunc("GET /auth/signin", controllers.getSignin)
	mux.HandleFunc("POST /auth/signin", controllers.postSignin)
}
