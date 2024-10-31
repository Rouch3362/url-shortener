package main

import (
	"log"

	"github.com/Rouch3362/url-shortener/api"
	"github.com/Rouch3362/url-shortener/cmd"
)

// main function to run the whole app
func main() {
	storage , err := cmd.ConnectionToDB()

	if err != nil {
		log.Fatal(err)
	}

	err = storage.InitDB()
	
	if err != nil {
		log.Fatal(err)
	}

	api.Run()
}