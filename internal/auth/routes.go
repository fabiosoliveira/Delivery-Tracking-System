package auth

import "net/http"

func registerRoutes(mux *http.ServeMux, controllers *controllers) {

	mux.HandleFunc("GET /auth/signup", controllers.getSignup)
	mux.HandleFunc("POST /auth/signup", controllers.postSignup)
	mux.HandleFunc("GET /auth/signin", controllers.getSignin)
	mux.HandleFunc("POST /auth/signin", controllers.postSignin)
}
