package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)



func LoadEnvVariable(varName string) string {
	// loads .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
	// gets env variable value based on its name
	value := os.Getenv(varName)

	return value
}


func ValidatePayload(username , password string) *Error{
	// checks if the request for creating user has requried fields
	if username == "" || password == "" {
		return &Error{
			Message: "username and password fields are required.", 
			Code: http.StatusBadRequest,
		}
	// check if length of values is longer than 8 characters
	} else if len(username) < 8 || len(password) < 8 {
		return &Error{
			Message: "username and password must be longer than 8 characters.",
			Code: http.StatusBadRequest,
		}
	}

	return nil
}

// a helper function for producing responses
func JsonGenerator(w http.ResponseWriter , statusCode int , value any) {
	w.Header().Add("Content-Type" , "application/json")
	
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(value)

	if err != nil {
		log.Fatal(err)
	}
}


func CreateJWT(user *LoginRequest) (string , error) {
	// claims (fields) we want in our jwt token
	claims := jwt.MapClaims{
		"username": user.Username,
		// expires after 1 day of creation
		"exp": time.Now().Add(time.Hour*24).UTC(),
	}

	// creating token based on ecryption algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256 , claims)

	// loads jwt secret from env file	
	JWT_SECRET := LoadEnvVariable("JWT_SECRET")
	// getting generated access token
	tokenString , err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "" , err
	}


	return tokenString,nil
}



// checks if a password and its hash is the same
func IsPasswordValid(hash , password string) *Error{
	// if not the same returns error
	err := bcrypt.CompareHashAndPassword([]byte(hash) , []byte(password))

	if err != nil {
		return &Error{
			Message: "username or password is invalid",
			Code: http.StatusUnauthorized,
		}
	}

	return nil
}

// generates error based on Error struct
func ErrorGenerator(w http.ResponseWriter, error *Error) {
	JsonGenerator(w,error.Code , error)
}

