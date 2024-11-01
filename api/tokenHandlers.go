package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rouch3362/url-shortener/cmd"
	"github.com/Rouch3362/url-shortener/types"
)

func (a *APIServer) LoginHandler(w http.ResponseWriter , r *http.Request) {
	// filling with client payload
	userPayload := &types.UserRequest{}

	err := json.NewDecoder(r.Body).Decode(userPayload)

	

	if err != nil {
		log.Fatal(err)
	}
	// checking for errors in payload
	validationErr := userPayload.Validator()

	if validationErr != "" {
		message := types.ErrorMessage{Message: validationErr}
		cmd.JsonGenerator(w, 400, message)
		return
	}

	// getting hashed password of user for comapring with entered password
	hashedPassword := a.DB.GetUserPassword(userPayload.Username)

	// checking if hash and password in same
	isPasswordValid := userPayload.ComparePassword(hashedPassword)

	if !isPasswordValid {
		message := types.ErrorMessage{Message: "username or password is not correct"}
		cmd.JsonGenerator(w , 401, message)
		return
	}

	// getting user from database
	userFromDB := a.DB.GetUserByUsername(userPayload.Username)
	// creating access and refresh token for sending to client
	accessToken := cmd.GenerateJWTToken(userFromDB, true)
	refreshToken := cmd.GenerateJWTToken(userFromDB, false)

	// formats it properly
	tokenResponse := types.Token{
		AcccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	// sending response
	cmd.JsonGenerator(w, 200, tokenResponse)
	
}