package api

import (
	"fmt"
	"net/http"
	"github.com/Rouch3362/url-shortener/cmd/db"
)

// an API server struct for accessing storage object and running its methods
type APIServer struct {
	Addr 	string
	DB		*db.Storage
}

// a method for APIServer to run the server
func (a *APIServer) Run() {
	router := http.NewServeMux()
	
	
	createUrlRoute := http.HandlerFunc(a.createUrlsHandler)
	router.Handle("/urls" , AuthMiddleware(createUrlRoute))
	router.HandleFunc("/users", a.createUserHandler)
	router.HandleFunc("/login", a.LoginHandler)

	refreshTokenRoute := http.HandlerFunc(a.RefreshTokenHandler)
	router.Handle("/refresh-token", AuthMiddleware(refreshTokenRoute))

	fmt.Println("Server is Running on port 8000")
	http.ListenAndServe(a.Addr, router)
}