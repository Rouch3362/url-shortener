package main

import (
	"net/http"
	"time"
)

// our custom error struct
type Error struct {
	Message string `json:"error"`
	Code  int    `json:"code"`
}

// wrote a method for our new custom error struct so we can use it for returning errors of functions
func (e *Error) Error() string {
	return e.Message
}

type User struct {
	ID        int    	`json:"id" db:"id"`
	Username  string 	`json:"username" db:"username"`
	Password  string 	`json:"password" db:"password"`
	CreatedAt time.Time	`json:"createdAt" db:"created_at"`
}

type UserRequest struct {
	Username  string 	`json:"username" db:"username"`
	Password  string 	`json:"password" db:"password"`
	CreatedAt time.Time	`json:"createdAt" db:"created_at"`
}

// a method for creating an instancee of user before saving it to database
func (u *UserRequest) CreateUser() (*UserRequest, *Error) {
	// checks if the request for creating user has requried fields
	if u.Username == "" || u.Password == "" {
		return nil , &Error{
			Message: "username and password fields are required.", 
			Code: http.StatusBadRequest,
		}
	}

	if len(u.Username) < 8 || len(u.Password) < 8 {
		return nil , &Error{
			Message: "username and password should be bigger than 8 characters.",
			Code: http.StatusBadRequest,
		}
	}

	user := &UserRequest{
		Username: u.Username, 
		Password: u.Password,
		CreatedAt: time.Now().UTC(),
	}

	return user , nil
}
