package main

import (
	"log"

	"github.com/Rouch3362/url-shortener/api"
	"github.com/Rouch3362/url-shortener/cmd/db"
)

// main function to run the whole app
func main() {
	storage , err := db.ConnectionToDB()

	if err != nil {
		log.Fatal(err)
	}

	err = storage.InitDB()
	
	if err != nil {
		log.Fatal(err)
	}

	apiServer := api.APIServer{Addr: ":8000", DB: storage}

	apiServer.Run()
}