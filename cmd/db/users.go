package db

import (
	"log"

	"github.com/Rouch3362/url-shortener/types"
)


// creating table for storing urls
func (s *Storage) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id 			SERIAL PRIMARY KEY UNIQUE,
		username	VARCHAR(100) NOT NULL UNIQUE,
		password 	VARCHAR(100) NOT NULL UNIQUE,
		created_at 	timestamp	NOT NULL DEFAULT now()
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

// creating new users
func (s *Storage)CreateNewUser(user *types.UserRequest) (*types.UserResponse, error){
	// needed query for inserting new users with returning columns
	query := `INSERT INTO users (username,password) VALUES ($1, $2) RETURNING id,username,created_at`

	// a empty instance to fill later with values returend from sql query
	userResponse := types.UserResponse{}  

	// executing query for creating user and saving result of the query
	err := s.DB.QueryRow(query, user.Username, user.Password).Scan(
		&userResponse.Id,
		&userResponse.Username,
		&userResponse.CreatedAt,
	)

	if err != nil {
		log.Fatal(err)
	}

	return &userResponse, nil
}

// getting user hashed password from database 
func (s *Storage) GetUserPassword(username string) string {
	query := `SELECT password FROM users WHERE username = $1`

	// an empty string for filling after query executing
	var hashedPassword string

	
	err := s.DB.QueryRow(query, username).Scan(&hashedPassword)


	if err != nil {
		log.Fatal(err)
	}

	return hashedPassword
}


// finding user by username from database
func (s *Storage) GetUserByUsername(username string) *types.UserResponse {
	query := `SELECT id,username,created_at FROM users WHERE username = $1`
	
	
	// creating empty instance of user response for filling after query executed
	result := &types.UserResponse{}

	err := s.DB.QueryRow(query, username).Scan(
		&result.Id,
		&result.Username,
		&result.CreatedAt,
	)
	if err != nil {
		log.Fatal(err)
	}

	return result
}