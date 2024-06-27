package main


type Error struct {
	Error 	string	`json:"error"`
	Code	int		`json:"code"`
}


type User struct {
	ID			int		`json:"id"`
	Username 	string	`json:"username"`
	Password	string	`json:"password"`
}


