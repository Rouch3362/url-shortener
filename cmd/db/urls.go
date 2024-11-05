package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Rouch3362/url-shortener/cmd"
	"github.com/Rouch3362/url-shortener/types"
	"github.com/lib/pq"
)

// creating table for storing users
func (s *Storage) createUrlsTable() error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		id 			SERIAL PRIMARY KEY UNIQUE,
		user_id	    INT REFERENCES users ON DELETE CASCADE NOT NULL,
		short_url 	VARCHAR(100) NOT NULL,
		long_url 	TEXT NOT NULL,
		clicks 		INT NOT NULL DEFAULT 0,	
		created_at 	timestamp	NOT NULL DEFAULT now()
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



func (s *Storage) CreateUrlDB(urlPayload *types.DBCreateUrlRequest) error {
	query := `INSERT INTO urls(user_id,long_url,short_url) VALUES ($1 , $2 , $3)`

	_, err := s.DB.Exec(query, urlPayload.UserId, urlPayload.LongUrl, urlPayload.ShortUrl)


	if err, ok := err.(*pq.Error); ok && err.Code.Name() == types.FOREIGN_KEY_ERROR {
		return errors.New("user does not exists (JWT token is not valid)")
	}


	return nil
}


func (s *Storage) IncreaseURLClicks(urlId string) {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short_url = $1`

	_, err := s.DB.Exec(query, urlId)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Storage) GetURL(urlId string) (string, error){
	query := `SELECT long_url FROM urls WHERE short_url = $1`

	W_ADDR := cmd.ReadEnvVar("W_ADDR") 
	// combining url id with domain name
	shortURL := W_ADDR+urlId

	var originalUrl string
	
	err := s.DB.QueryRow(query, shortURL).Scan(&originalUrl)

	if err == sql.ErrNoRows {
		return "", errors.New("URL not found")
	}
	// increasing click field whenever we pull the url from database
	s.IncreaseURLClicks(shortURL)

	return originalUrl,nil
}