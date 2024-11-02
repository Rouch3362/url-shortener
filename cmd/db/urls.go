package db

import (
	"log"

	"github.com/Rouch3362/url-shortener/types"
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



func (s *Storage) createUrls(urlPayload *types.DBCreateUrlRequest) error {
	query := `INSERT INTO urls(user_id,long_url,short_url) VALUES ($1 , $2 , $3)`

	_, err := s.DB.Exec(query)


	if err != nil {
		log.Fatal(err)
	}


	return nil
}