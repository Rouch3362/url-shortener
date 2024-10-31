package cmd

import (
	"log"

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