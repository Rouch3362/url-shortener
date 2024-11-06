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
	
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	
	createUrlRoute := http.HandlerFunc(a.createUrlsHandler)
	subrouter.Handle("/urls", AuthMiddleware(createUrlRoute)).Methods("POST")
	
	deleteUrlRoute := http.HandlerFunc(a.deleteUrlHandler)
	subrouter.Handle("/urls/{id}", AuthMiddleware(deleteUrlRoute)).Methods("DELETE")
	subrouter.HandleFunc("/{id}", a.getUrlHandler).Methods("GET")
	subrouter.HandleFunc("/register", a.createUserHandler).Methods("POST")
	subrouter.HandleFunc("/users/{username}", a.GetUser).Methods("GET")

	deleteUserRoute := http.HandlerFunc(a.DeleteUser)

	subrouter.Handle("/users", AuthMiddleware(deleteUserRoute)).Methods("DELETE")
	// subrouter.HandleFunc("/whoami", a.).Methods("GET")
	subrouter.HandleFunc("/login", a.LoginHandler).Methods("POST")

	refreshTokenRoute := http.HandlerFunc(a.RefreshTokenHandler)
	subrouter.Handle("/refresh-token", AuthMiddleware(refreshTokenRoute)).Methods("POST")

	fmt.Println("Server is Running on port 8000")
	http.ListenAndServe(a.Addr, router)
}