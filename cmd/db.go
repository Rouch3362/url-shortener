package cmd

import (
	"database/sql"
	"fmt"
	// required for using postgres driver
	_ "github.com/lib/pq"
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

// initializing database tables
func (s *Storage) InitDB() error {
	err := s.createUserTable()

	if err != nil {
		return err
	}

	err = s.createUrlsTable()

	return err
}


// creating table for storing urls
func (s *Storage) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id 			SERIAL PRIMARY KEY UNIQUE,
		username	VARCHAR(100) NOT NULL UNIQUE,
		password 	VARCHAR(100) NOT NULL UNIQUE,
		create_at 	timestamp	NOT NULL
	)`


	_ , err := s.DB.Exec(query)

	if err != nil {
		return err
	}

	// making the most requested column and index for accessing it faster
	idxQuery := "CREATE INDEX IF NOT EXISTS users_index ON users(username)"

	_ , err = s.DB.Query(idxQuery)

	return err
}
// creating table for storing users
func (s *Storage) createUrlsTable() error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		id 			SERIAL PRIMARY KEY UNIQUE,
		user_id	    INT REFERENCES users ON DELETE CASCADE NOT NULL,
		short_url 	VARCHAR(100) NOT NULL,
		long_url 	TEXT NOT NULL,
		clicks 		INT NOT NULL,	
		create_at 	timestamp	NOT NULL
	)`


	_ , err := s.DB.Exec(query)

	if err != nil {
		return err
	}

	// making the most requested column and index for accessing it faster
	idxQuery := "CREATE INDEX IF NOT EXISTS urls_index ON urls(short_url)"

	_ , err = s.DB.Query(idxQuery)

	return err
}