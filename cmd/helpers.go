package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func ReadEnvVar(varName string) string {
	env, err := godotenv.Read(".env")

	if err != nil{
		log.Fatal(err)
	}

	value := env[varName]

	return value 
}



func JsonGenerator(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Add("Content-Type" , "application/json")

	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(value)

	if err != nil {
		log.Fatal(err)
	}
}