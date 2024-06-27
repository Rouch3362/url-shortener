package main

import (
	"fmt"
	"log"
)

func main() {
	db , err := NewDB()

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":8000" , *db)

	runErr := server.Run() 
	
	if runErr != nil {
		log.Fatal(err)
	}

	fmt.Println("server listening on port:8000")
}