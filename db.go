package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"log"
	"net/http"
)


type DBCommands interface {
	// user commands
	CreateUserDB(*User) (*User , error)
	GetUserByUsernameDB(string) (*User , *Error)
	GetUserByIDDB(int) (*User , *Error)
	// delete user
	// update user

	// refresh token commands
	CreateRefreshTokenDB() (string , error)
	DeleteRefreshTokenDB() error
	GetRefreshTokenDB(string) *Error
	
	// url commands
	CreateUrlDB() (string , *Error)
	GetUrlByShortUrl(string) (*Url , *Error)
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
	if err != nil {
		return err
	}

	err = s.CreateJwtTable()

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

func (s *Storage) CreateJwtTable() error {
	query := `CREATE TABLE IF NOT EXISTS token (
		id 			SERIAL PRIMARY KEY UNIQUE,
		user_id 	INT REFERENCES users NOT NULL,
		refresh		TEXT NOT NULL,
		created_at	timestamp NOT NULL
	)`

	_ ,err := s.DB.Query(query)

	if err != nil {
		return err
	}

	return nil
}

// simply just creates a row in database for tokens that provided
func (s *Storage) CreateRefreshTokenDB(userId int , token string) error {
	query := `INSERT INTO token (user_id , refresh , created_at) VALUES (
		$1 , $2 , $3
	)` 


	_,err := s.DB.Exec(query , userId, token ,time.Now().UTC())

	if err != nil {
		return err
	}

	return nil
}


func (s *Storage) DeleteRefreshTokenDB(token string) *Error {
	query := `DELETE FROM token WHERE refresh=$1`

	
	res ,err := s.DB.Exec(query , token)
	
	if err != nil {
		log.Fatal(err)
	} 
	// checks if a row deleted or not and if not returns and error that means token already used
	if count , err := res.RowsAffected(); err == nil && count < 1 {
		if err != nil {
			log.Fatal(err)
		}
		return &Error{
			Message: "token is not valid any more(black listed).",
			Code: http.StatusUnauthorized,
		}
	}
	return nil
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

	// insert each of columns in the userInstance
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
		$1,$2,$3) RETURNING *`
		
	// an empty instance for user
	foundedUser := User{}
	// the only column it returns is id 
	err := s.DB.QueryRow(query , user.Username , user.Password , user.CreatedAt).Scan(
		&foundedUser.ID,
		&foundedUser.Username,
		&foundedUser.Password,
		&foundedUser.CreatedAt,
	)


	if err != nil && strings.Contains(err.Error(),"duplicate") {
		return nil , &Error{
			Message: fmt.Sprintf("user with username: %s already exists." , user.Username),
			Code: http.StatusConflict,
		}
	}

	if err != nil && !strings.Contains(err.Error(),"duplicate") {
		log.Fatal(err)
	}
	

	return &foundedUser, nil
}


func (s *Storage) CreateUrlDB(url *Url) (*Url,*Error) {
	query := `INSERT INTO urls (user_id , old_url , new_url , created_at) VALUES (
		$1 , $2 , $3 , $4) RETURNING *`

	// creating instance for filling the returen results from insert command
	createdUrl := Url{}
	err := s.DB.QueryRow(query , url.User , url.OldUrl, url.NewUrl, url.CreatedAt).Scan(
		&createdUrl.ID,
		&createdUrl.User,
		&createdUrl.OldUrl,
		&createdUrl.NewUrl,
		&createdUrl.CreatedAt,
	)

	if err != nil {
		log.Fatal(err)
	}

	return &createdUrl, nil
}


func (s *Storage) GetUrlByShortUrl(shortUrl string) (*Url , *Error) {
	query := `SELECT * FROM urls WHERE new_url=$1`

	url := Url{}
	
	err := s.DB.QueryRow(query , shortUrl).Scan(&url.ID , &url.User , &url.OldUrl , &url.NewUrl , &url.CreatedAt)

	fmt.Println(shortUrl)

	if err == sql.ErrNoRows {
		return nil, &Error{
			Message: "url not found.",
			Code: http.StatusNotFound,
		}
	}

	return &url,nil
}
