package api

import (
	"fmt"
	"net/http"
	"github.com/Rouch3362/url-shortener/cmd/db"
	"github.com/gorilla/mux"
)

// an API server struct for accessing storage object and running its methods
type APIServer struct {
	Addr 	string
	DB		*db.Storage
}

// a method for APIServer to run the server
func (a *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/urls" , a.createUrlsHandler).Methods("POST")
	router.HandleFunc("/users", a.createUserHandler).Methods("POST")
	router.HandleFunc("/login", a.LoginHandler).Methods("POST")

	fmt.Println("Server is Running on port 8000")
	http.ListenAndServe(a.Addr, router)
}