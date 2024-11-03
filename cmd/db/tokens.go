package db

import (
	"log"

	"github.com/Rouch3362/url-shortener/types"
)



func (s *Storage) createTokenTable() error {
	query := `CREATE TABLE IF NOT EXISTS tokens (
		id					SERIAL PRIMARY KEY NOT NULL,
		user_id				INT REFERENCES users ON DELETE CASCADE NOT NULL UNIQUE,
		access_token		TEXT NOT NULL,
		refresh_token		TEXT NOT NULL,
		expires_at			INT NOT NULL

	)`

	_, err := s.DB.Exec(query)

	if err != nil {
		return err
	}

	idxQuery := `CREATE INDEX IF NOT EXISTS token_index ON tokens(refresh_token)`

	_,err = s.DB.Exec(idxQuery)

	return err
}


func (s *Storage) SaveToken(tokenInfo *types.TokenDBRequest) {
	query := `INSERT INTO tokens (user_id,access_token,refresh_token,expires_at) VALUES ($1, $2, $3, $4)`
	

	_,err := s.DB.Exec(query, tokenInfo.UserId, tokenInfo.AccessToken, tokenInfo.RefreshToken, tokenInfo.ExpiresAt)

	if err != nil {
		log.Fatal(err)
	}
}


func (s *Storage) DoesRefreshTokenExists(refreshToken string) bool {
	query := `SELECT EXISTS(SELECT id FROM tokens WHERE refresh_token = $1)`

	var exists bool
	err := s.DB.QueryRow(query, refreshToken).Scan(&exists)

	if exists {
		err := s.RemoveRefreshToken(refreshToken)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	return exists
}


func (s *Storage) RemoveRefreshToken(refreshToken string) error {
	query := `DELETE FROM tokens WHERE refresh_token = $1`

	_, err := s.DB.Exec(query, refreshToken)


	if err != nil {
		return err
	}

	return nil
}