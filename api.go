package main

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type APIServer struct {
	Addr string
	DB   Storage
}

func NewApiServer(addr string, db Storage) *APIServer {
	apiInstance := &APIServer{
		Addr: addr,
		DB:   db,
	}

	return apiInstance
}

func (a *APIServer) Run() error {
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	subRouter.HandleFunc("/hello" , SayHello)

	err := http.ListenAndServe(a.Addr , router)


	return err

}


func SayHello(w http.ResponseWriter , r *http.Request) {
	helloStruct := struct{Message string}{"hello world"}

	JsonGenerator(w , http.StatusOK , helloStruct)

} 




