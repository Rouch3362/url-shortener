package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	db , err := NewDB()

	envErr := godotenv.Load(".env")

	if envErr != nil {
		log.Fatal(envErr)
	}

	PORT := os.Getenv("PORT")


	if err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	server := NewApiServer(PORT , *db)

	runErr := server.Run() 
	
	if runErr != nil {
		log.Fatal(err)
	}

	fmt.Println("server listening on port:8000")
}