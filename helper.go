package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func JsonGenerator(w http.ResponseWriter , statusCode int , value any) {
	w.Header().Add("Content-Type" , "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(value)

	if err != nil {
		log.Fatal(err)
	}

}

