package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()

	router.HandleFunc("/urls" , createUrlsHandler).Methods("POST")

	fmt.Println("Server is Running on port 8000")
	http.ListenAndServe(":8000", router)
}