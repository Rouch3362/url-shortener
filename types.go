package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// our custom error struct
type Error struct {
	Message string `json:"message"`
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
	} else if len(u.Username) < 8 || len(u.Password) < 8 {
		return nil , &Error{
			Message: "username and password must be longer than 8 characters.",
			Code: http.StatusBadRequest,
		}
	}
	// hashing password
	hashPassword , hashErr := bcrypt.GenerateFromPassword([]byte(u.Password) , bcrypt.DefaultCost)
	
	if hashErr != nil {
		log.Fatal(hashErr)
	}
	user := &UserRequest{
		Username: u.Username, 
		Password: string(hashPassword),
		CreatedAt: time.Now().UTC(),
	}

	return user , nil
}
