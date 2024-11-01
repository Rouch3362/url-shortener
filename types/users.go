package types

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// struct of user payload for creating user
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// hashing user's password
func (u *UserRequest) HashPassword() {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(u.Password) , 10)

	if err != nil {
		log.Fatal(err)
	}

	u.Password = string(hashByte)
}

// comparing plain texted password and its hash
func (u *UserRequest) ComparePassword(hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password))
	
	return err != nil
}

// validating payloads of user
func (c *UserRequest) Validator() string {

	if len(c.Username) < 3 {
		return "username is a required field and must be greater than or equal to 3 characters"
	}


	if len(c.Password) < 4 || len(c.Password) > 72 {
		return "password is a required field and must be at least 4 characters and maximum 72 characters"
	}

	return ""
}