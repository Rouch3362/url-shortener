package db

import (
	"database/sql"
	"fmt"

	// required for using postgres driver
	"github.com/Rouch3362/url-shortener/cmd"
	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func ConnectionToDB() (*Storage , error){
	// loading env variables
	dbUser := cmd.ReadEnvVar("DB_USER")
	dbPass := cmd.ReadEnvVar("DB_PASS")
	dbName := cmd.ReadEnvVar("DB_NAME")


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

// initializing database tables
func (s *Storage) InitDB() error {
	err := s.createUserTable()

	if err != nil {
		return err
	}

	err = s.createUrlsTable()

	if err != nil {
		return err
	}

	err = s.createTokenTable()


	return err
}

