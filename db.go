package main

import (
	"database/sql"
	"fmt"
	"strings"

	"log"
	"net/http"

)


type DBCommands interface {
	CreateUserDB(*User) (*User , error)
	GetUserByUsernameDB(string) (*User , *Error)
	GetUserByIDDB(int) (*User , *Error)
}


type Storage struct {
	DB *sql.DB
}




func NewDB() (*Storage, error) {
	DB_USER := LoadEnvVariable("DB_USER")
	DB_NAME := LoadEnvVariable("DB_NAME")
	DB_PASS := LoadEnvVariable("DB_PASS")
	connectionStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DB_USER , DB_NAME , DB_PASS)
	
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
	if err != nil {
		return err
	}
	err  = s.CreateUrlsTable()
	
	return err
	
}



func (s *Storage) CreateUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
		id			SERIAL PRIMARY KEY UNIQUE,
		username	VARCHAR(100) NOT NULL UNIQUE,
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



func (s *Storage) GetUserByUsernameDB(username string) (*User , *Error) {
	query := "SELECT * FROM users WHERE username=$1"

	// created an instance for filling it with result from database
	user := User{}

	// QueryRow returns only one row and if we use scan after it it will return an error or nil
	// scan accepts destination for returned columns from database. in this case we didn't use RETURNING in postgres so it will return all columns
	err := s.DB.QueryRow(query , username).Scan(&user.ID , &user.Username , &user.Password , &user.CreatedAt)

	// this will occure when no result founded
	if err == sql.ErrNoRows {
		return nil, &Error{Message: fmt.Sprintf("user with username: %s not found", username) , Code: http.StatusNotFound}
	}

	return &user , nil
}


func (s *Storage) GetUserByIDDB(id int) (*User , *Error) {
	query := "SELECT * FROM users WHERE id=$1"
	
	user := User{}

	err := s.DB.QueryRow(query , id).Scan(&user.ID , &user.Username , &user.Password , &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, &Error{
			Message: "user not found",
			Code: http.StatusNotFound,
		}
	}

	return &user , nil
}

func (s *Storage) CreateUserDB(user *UserRequest) (*User , *Error) {
	// we use returning for insert because postgres by default will not return columns in insert command so we use it for fetching user and sending response to request source
	query := `INSERT INTO users (username , password, created_at) VALUES (
		$1,$2,$3) RETURNING id`
		
	// an empty instance for user id
	var id int
	// the only column it returns is id 
	err := s.DB.QueryRow(query , user.Username , user.Password , user.CreatedAt).Scan(&id)


	if err != nil && strings.Contains(err.Error(),"duplicate") {
		return nil , &Error{
			Message: fmt.Sprintf("user with username: %s already exists." , user.Username),
			Code: http.StatusConflict,
		}
	}

	if err != nil && !strings.Contains(err.Error(),"duplicate") {
		log.Fatal(err)
	}
	
	// fetching user by username
	foundedUser , findingErr := s.GetUserByIDDB(id)

	if findingErr != nil {
		return nil , findingErr
	}
	// we don't use pointer in this code and above error, because the GetUserByUsernameDB already returns pointer for user
	return foundedUser, nil
}