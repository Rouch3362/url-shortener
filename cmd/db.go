package cmd

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	DB *sql.DB
}

func ConnectionToDB() (*Storage , error){
	// loading env variables
	dbUser := ReadEnvVar("DB_USER")
	dbPass := ReadEnvVar("DB_PASS")
	dbName := ReadEnvVar("DB_NAME")


	connectionString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", dbUser , dbName , dbPass)

	// opening new connection to postgresql
	db , err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil , err
	}

	storageInstance := &Storage{
		DB: db,
	}


	return storageInstance , nil

}

