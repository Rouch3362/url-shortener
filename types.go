package main

import (
	"log"
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

// for showing users their tokens
type JwtToken struct {
	Access 	string 	`json:"access"`
	Refresh	string	`json:"refresh" db:"token"`
}

// type for refreshing expired tokens
type RefershTokenRequest struct {
	Refresh string	`json:"refresh" db:"token"`
}

type LoginRequest struct {
	Username string 	`json:"username" db:"username"`
	Passwrod string		`json:"password" db:"password"`
}

// a method for creating an instancee of user before saving it to database
func (u *UserRequest) CreateUser() (*UserRequest, *Error) {
	validateErr := ValidatePayload(u.Username , u.Password)

	if validateErr != nil {
		return nil,validateErr
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
