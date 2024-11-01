package db

import "log"



func (s *Storage) CreateTokenTable() {
	query := `CREATE TABLE IF NOT EXISTS tokens (
		id					SERIAL PRIMARY KEY NOT NULL,
		user_id				INT REFERENCES users ON DELETE CASCADE NOT NULL UNIQUE,
		access_token		TEXT NOT NULL,
		refersh_token		TEXT NOT NULL,
		expires_at			timestamp NOT NULL

	)`

	_, err := s.DB.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

