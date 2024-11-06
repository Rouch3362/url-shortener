package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"github.com/Rouch3362/url-shortener/types"
	"github.com/lib/pq"
)

// creating table for storing urls
func (s *Storage) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id 			SERIAL PRIMARY KEY UNIQUE,
		username	VARCHAR(100) NOT NULL UNIQUE,
		password 	VARCHAR(100) NOT NULL UNIQUE,
		created_at 	timestamp	NOT NULL DEFAULT now()
	)`

	_, err := s.DB.Exec(query)

	if err != nil {
		return err
	}

	// making the most requested column and index for accessing it faster
	idxQuery := "CREATE INDEX IF NOT EXISTS users_index ON users(username)"

	_, err = s.DB.Query(idxQuery)

	return err
}

// creating new users
func (s *Storage) CreateNewUser(user *types.UserRequest) (*types.UserResponse, error) {
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

	if err, ok := err.(*pq.Error); ok && err.Code.Name() == types.UNIQUE_VIOLATION {
		return nil, fmt.Errorf("user with username %v already exists", user.Username)
	}

	return &userResponse, nil
}

// getting user hashed password from database
func (s *Storage) GetUserPassword(username string) (string, error) {
	query := `SELECT password FROM users WHERE username = $1`

	// an empty string for filling after query executing
	var hashedPassword string

	err := s.DB.QueryRow(query, username).Scan(&hashedPassword)

	if err == sql.ErrNoRows {
		return "",errors.New("user not found")
	}

	return hashedPassword,nil
}

// finding user by username from database
func (s *Storage) GetUserByUsername(username string) (*types.UserResponse,error) {
	query := `SELECT id,username,created_at FROM users WHERE username = $1`

	// creating empty instance of user response for filling after query executed
	result := &types.UserResponse{}

	err := s.DB.QueryRow(query, username).Scan(
		&result.Id,
		&result.Username,
		&result.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}


	return result, err
}

func (s *Storage) GetUserURLs(username string) (*types.UserURLsResponse,error) {
	query := `SELECT 
		users.id,
		urls.id,
		urls.long_url, 
		urls.short_url, 
		users.username, 
		users.created_at, 
		urls.clicks, 
		urls.created_at FROM users 
		JOIN urls ON urls.user_id = users.id 
		WHERE users.username = $1`

	
	// an instance for saving user info and their urls
	response := types.UserURLsResponse{} 

	// executing query
	result, _ := s.DB.Query(query, username)
	

	// iterating over each returend values
	for result.Next() {
		// creating an url instance for each result
		url := types.URL{}
		// filling every field of response
		err := result.Scan(
			&response.Id,
			&url.Id,
			&url.OriginalURL,
			&url.ShortURL,
			&response.Username,
			&response.CreatedAt,
			&url.Clicks,
			&url.CreatedAt,
		)

		if err != nil {
			log.Fatal(err)
		}
		// appending user's urls scanned urls
		response.Urls = append(response.Urls, url)
	}

	// if user didn't have any urls in database
	if len(response.Urls) < 1 {
		return nil,errors.New("user has no urls or user does not exists")
	}	

	return &response,nil
}


func (s *Storage) DeleteUserDB(username string) error {
	query := `DELETE FROM users WHERE username = $1`

	result, err := s.DB.Exec(query, username)

	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}