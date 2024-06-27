package main

import (
	"database/sql"
)

type Storage struct {
	DB *sql.DB
}




func NewDB() (*Storage, error) {
	connectionStr := "user=postgres dbname=postgres password=amirali3362 sslmode=disable"
	
	db , err := sql.Open("postgres" , connectionStr)

	if err != nil {
		return nil , err
	}

	storageInstance := &Storage{
		DB: db,
	}

	return storageInstance , nil
}


func (s *Storage) Init() error{
	err := s.CreateUsersTable()
	err  = s.CreateUrlsTable()
	
	return err
}



func (s *Storage) CreateUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
		id			SERIAL PRIMARY KEY UNIQUE,
		username	VARCHAR(100) NOT NULL,
		password	VARCHAR(100) NOT NULL,
		created_at	timestamp 	 NOT NULL
	)`

	_ , err := s.DB.Exec(query)

	return err
}


func (s *Storage) CreateUrlsTable() error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		id			SERIAL PRIMARY KEY UNIQUE,
		user_id		INT REFERENCES users NOT NULL,
		old_url		TEXT		 NOT NULL,
		new_url 	VARCHAR(200) NOT NULL,
		created_at 	timestamp	 NOT NULL
	)`

	_ , err := s.DB.Exec(query)


	return err
}