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
	GetUserByUsernameDB(string) (*UserResponse , *Error)
	GetUserByIDDB(int) (*UserResponse , *Error)
	GetUserUrlsDB(string) (*UserUrlsResponse , *Error)
	DeleteUserDB(string) (string , *Error)
	UpdateUserDB(*User) (*UserResponse , *Error) 

	// refresh token commands
	CreateRefreshTokenDB() (string , error)
	DeleteRefreshTokenDB() error
	GetRefreshTokenDB(string) *Error
	
	// url commands
	CreateUrlDB() (string , *Error)
	UpdateUrlDB(int) (*UrlReponse,*Error)
	DeleteUrlDB(string) *Error
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
		user_id		INT REFERENCES users ON DELETE CASCADE NOT NULL,
		old_url		TEXT		 NOT NULL,
		new_url 	VARCHAR(200) NOT NULL,
		clicks		INT			 NOT NULL,
		created_at 	timestamp	 NOT NULL
	)`

	_ , err := s.DB.Exec(query)


	return err
}

func (s *Storage) CreateJwtTable() error {
	query := `CREATE TABLE IF NOT EXISTS token (
		id 			SERIAL PRIMARY KEY UNIQUE,
		user_id 	INT REFERENCES users ON DELETE CASCADE NOT NULL UNIQUE,
		refresh		TEXT NOT NULL,
		created_at	timestamp NOT NULL
	)`

	_ ,err := s.DB.Query(query)

	if err != nil {
		return err
	}

	return nil
}

// deletes token if exists by a user id
func (s *Storage) DeletePreviousToken(userId int)  {
	query := `DELETE FROM token WHERE user_id=$1 RETURNING refresh`

	var refresh string

	s.DB.QueryRow(query,userId).Scan(&refresh)
}

// simply just creates a row in database for tokens that provided
func (s *Storage) CreateRefreshTokenDB(userId int , token string) error {
	// first deletes token if a token already exist in database and not expired
	s.DeletePreviousToken(userId)
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
		return BlackListedTokenError()
	}
	return nil
}


func (s *Storage) GetUserByUsernameDB(username string) (*UserResponse, string , *Error) {
	query := "SELECT * FROM users WHERE username=$1"

	// created an instance for filling it with result from database
	user := UserResponse{}
	var userPass string 
	// QueryRow returns only one row and if we use scan after it it will return an error or nil
	// scan accepts destination for returned columns from database. in this case we didn't use RETURNING in postgres so it will return all columns
	err := s.DB.QueryRow(query , username).Scan(&user.ID , &user.Username, &userPass , &user.CreatedAt)
	// this will occure when no result founded
	if err == sql.ErrNoRows {
		return nil,"",NotFoundError(username)
	} else if err != nil {
		log.Fatal(err)
	}

	return &user, userPass , nil
}


func (s *Storage) GetUserByIDDB(id int) (*UserResponse , *Error) {
	query := "SELECT id,username,created_at FROM users WHERE id=$1"
	
	user := UserResponse{}
	
	// insert each of columns in the userInstance
	err := s.DB.QueryRow(query , id).Scan(&user.ID , &user.Username , &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, NotFoundError("user")
	}

	return &user , nil
}


func (s *Storage) GetUserUrlsDB(username string) (*UserUrlsResponse , *Error) {
	// first check if a user with that username exists
	_ , _ , userErr := s.GetUserByUsernameDB(username)

	if userErr != nil {
		return nil , userErr
	}

	// ==============================================

	query := `SELECT users.id, users.username, urls.id, urls.old_url, urls.new_url,urls.clicks, urls.created_at, users.created_at 
			  FROM users 
			  JOIN urls ON urls.user_id = users.id WHERE users.username=$1`

	result , err := s.DB.Query(query , username)
	
	if err != nil {
		log.Fatal(err)
	}
	// set an instance for filling results in it
	response := UserUrlsResponse{}

	// iterate over results and add it to response 
	for result.Next() {
		// an instance for url so in end of every iterate we append it to urls field response
		urls := ShortedFormUrl{}
		// filling infos of every field
		scanErr := result.Scan(
			// user info field
			&response.ID,
			&response.Username,
			// urls field info
			&urls.ID,
			&urls.OldUrl,
			&urls.NewUrl,
			&urls.Clicks,
			&urls.CreatedAt,
			// user info field
			&response.CreatedAt,
		)
		// checks if there are no rows
		if scanErr == sql.ErrNoRows {
			return nil , NotFoundError("urls of "+username)
		}
		// append url instance to urls field in response instance
		response.Urls = append(response.Urls, urls)
	
	}

	// checks if there are any urls for user with that given username
	if len(response.Urls) < 1 {
		return nil , &Error{
			Message: "user has no urls.",
			Code: http.StatusNotFound,
		}
	}
	
	return &response,nil
}


func (s *Storage) DeleteUserDB(username string) (*UserResponse, *Error) {
	query := `DELETE FROM users WHERE username=$1 RETURNING id,username,created_at`
	
	user := UserResponse{}

	err := s.DB.QueryRow(query , username).Scan(&user.ID,&user.Username,&user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, NotFoundError(username)
	} else if err != nil {
		log.Fatal(err)
	}

	return &user , nil
	
}

func (s *Storage) CreateUserDB(user *User) (*UserResponse , *Error) {
	// we use returning for insert because postgres by default will not return columns in insert command so we use it for fetching user and sending response to request source
	query := `INSERT INTO users (username , password, created_at) VALUES (
		$1,$2,$3) RETURNING id,username,created_at`
		
	// an empty instance for user
	foundedUser := UserResponse{}

	// the only column it returns is id 
	err := s.DB.QueryRow(query , user.Username , user.Password , user.CreatedAt).Scan(
		&foundedUser.ID,
		&foundedUser.Username,
		&foundedUser.CreatedAt,
	)


	if err != nil && strings.Contains(err.Error(),"duplicate") {
		return nil, &Error{
			Message: fmt.Sprintf("user with username: %s already exists." , user.Username),
			Code: http.StatusConflict,
		}
	}

	if err != nil && !strings.Contains(err.Error(),"duplicate") {
		log.Fatal(err)
	}
	

	return &foundedUser, nil
}


func (s *Storage) UpdateUserDB(newUsername string , oldUsername string) (*UserResponse , *Error) {
	query := `UPDATE users SET username=$1 WHERE username=$2 RETURNING id,username,created_at`

	var response UserResponse

	err := s.DB.QueryRow(query , newUsername,oldUsername).Scan(&response.ID,&response.Username,&response.CreatedAt)

	if err == sql.ErrNoRows {
		return nil , NotFoundError(oldUsername)
	} else if err != nil {
		log.Fatal(err)
	}
	
	return &response, nil

}


func (s *Storage) CreateUrlDB(url *Url) (*Url,*Error) {
	query := `INSERT INTO urls (user_id , old_url , new_url, clicks , created_at) VALUES (
		$1 , $2 , $3, $4 , $5) RETURNING *`

	// creating instance for filling the returen results from insert command
		createdUrl := Url{}
	err := s.DB.QueryRow(query , url.User , url.OldUrl, url.NewUrl, url.Clicks ,url.CreatedAt).Scan(
		&createdUrl.ID,
		&createdUrl.User,
		&createdUrl.OldUrl,
		&createdUrl.NewUrl,
		&createdUrl.Clicks,
		&createdUrl.CreatedAt,
	)

	if err != nil {
		log.Fatal(err)
	}

	return &createdUrl, nil
}


func (s *Storage) UpdateUrlDB(updatedUrl string,id int) (*UrlReponse,*Error) {
	query := `UPDATE urls SET old_url=$1 WHERE id=$2 RETURNING new_url`

	var short_url string
	err := s.DB.QueryRow(query,updatedUrl,id).Scan(&short_url)

	if err == sql.ErrNoRows {
		return nil, NotFoundError("url")
	} else if err != nil {
		log.Fatal(err)
	}


	response , _ := s.GetUrlByShortUrl(short_url) 

	return response , nil
}


func  (s *Storage) DeleteUrlDB(id int) *Error {
	query := `DELETE FROM urls WHERE id=$1 RETURNING new_url`

	var new_url_returned string
	err := s.DB.QueryRow(query,id).Scan(&new_url_returned)
	// checks if url existed or not
	if err == sql.ErrNoRows {
		return NotFoundError("url")
	} else if err != nil {
		log.Fatal(err)
	}
	return nil
}


func (s *Storage) IncreaseClicksUrlDB(newUrl string) {
	query := `UPDATE urls SET clicks=clicks+1 WHERE new_url=$1`

	_,  err := s.DB.Exec(query , newUrl)

	if err != nil {
		log.Fatal(err)
	}

}


func (s *Storage) GetUrlByShortUrl(shortUrl string) (*UrlReponse , *Error) {
	query := `SELECT urls.id, urls.user_id, urls.old_url, urls.new_url, urls.clicks,urls.created_at, users.id, users.username, users.created_at 
			  FROM urls  
			  JOIN users ON urls.user_id = users.id WHERE new_url=$1`

	url := UrlReponse{}
	
	err := s.DB.QueryRow(query , shortUrl).Scan(
		&url.ID,
		&url.User.ID,
		&url.OldUrl,
		&url.NewUrl,
		&url.Clicks,
		&url.CreatedAt,
		&url.User.ID,
		&url.User.Username,
		&url.User.CreatedAt,
	)



	if err == sql.ErrNoRows {
		return nil, NotFoundError(shortUrl)
	} else if err != nil{
		log.Fatal(err)
	}


	return &url,nil
}
