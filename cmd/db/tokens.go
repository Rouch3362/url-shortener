package db

import (
	"log"
	"github.com/Rouch3362/url-shortener/types"
)


// creating table for saving refresh tokens
func (s *Storage) createTokenTable() error {
	query := `CREATE TABLE IF NOT EXISTS tokens (
		id					SERIAL PRIMARY KEY NOT NULL,
		user_id				INT REFERENCES users ON DELETE CASCADE NOT NULL UNIQUE,
		refresh_token		TEXT NOT NULL,
		expires_at			INT NOT NULL

	)`

	_, err := s.DB.Exec(query)

	if err != nil {
		return err
	}
	// creating index for refresh_token field
	idxQuery := `CREATE INDEX IF NOT EXISTS token_index ON tokens(refresh_token)`

	_,err = s.DB.Exec(idxQuery)

	return err
}

// function for saving tokens in database
func (s *Storage) SaveToken(tokenInfo *types.TokenDBRequest) {
	// first removing any previous refresh tokens then saving new one
	s.RemovePreviousTokens(tokenInfo.UserId)
	query := `INSERT INTO tokens (user_id,refresh_token,expires_at) VALUES ($1, $2, $3)`
	

	_,err := s.DB.Exec(query, tokenInfo.UserId, tokenInfo.RefreshToken, tokenInfo.ExpiresAt)

	if err != nil {
		log.Fatal(err)
	}
}


// checking for token existance
func (s *Storage) DoesRefreshTokenExists(refreshToken string) bool {
	query := `SELECT EXISTS(SELECT id FROM tokens WHERE refresh_token = $1)`

	// checking if above query returend result
	var exists bool
	err := s.DB.QueryRow(query, refreshToken).Scan(&exists)

	if err != nil {
		log.Fatal(err)
	}

	return exists
}

// removing refresh token (this function is not useful for now but it will remain here for future usage)
func (s *Storage) RemoveRefreshToken(refreshToken string) error {
	query := `DELETE FROM tokens WHERE refresh_token = $1`

	_, err := s.DB.Exec(query, refreshToken)


	if err != nil {
		return err
	}

	return nil
}

// removing user's all refresh tokens saved into database 
func (s *Storage) RemovePreviousTokens(userId int) {
	query := `DELETE FROM tokens WHERE user_id = $1`
	
	_, err := s.DB.Exec(query, userId)

	if err != nil {
		log.Fatal(err)
	}
}