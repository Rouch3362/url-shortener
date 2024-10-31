package db



// creating table for storing urls
func (s *Storage) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id 			SERIAL PRIMARY KEY UNIQUE,
		username	VARCHAR(100) NOT NULL UNIQUE,
		password 	VARCHAR(100) NOT NULL UNIQUE,
		created_at 	timestamp	NOT NULL
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