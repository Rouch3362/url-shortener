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
	tokenResponse := cmd.GenerateAuthTokens(userFromDB)

	tokenDB := types.TokenDBRequest{
		UserId: userFromDB.Id,
		AccessToken: tokenResponse.AcccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt: cmd.ExpireationTime(false),
	}

	a.DB.SaveToken(&tokenDB)

	// sending response
	cmd.JsonGenerator(w, 200, tokenResponse)
	
}



func (a *APIServer) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshInstance := &types.RefreshTokenRequest{}

	err := json.NewDecoder(r.Body).Decode(refreshInstance)

	if err != nil {
		log.Fatal(err)
	}

	validationErr := refreshInstance.Validate()

	if validationErr != "" {
		message := types.ErrorMessage{Message: validationErr}
		cmd.JsonGenerator(w, 400, message)
		return
	}

	isRefreshTokenValid := a.DB.DoesRefreshTokenExists(refreshInstance.RefreshToken)


	if !isRefreshTokenValid {
		message := types.ErrorMessage{Message: "refresh token is not valid"}
		cmd.JsonGenerator(w, 401, message)
		return
	}

	userCredentials , err := cmd.DecodeJWTToken(refreshInstance.RefreshToken)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, 401, message)
		return
	}

	tokenResponse := cmd.GenerateAuthTokens(userCredentials)

	cmd.JsonGenerator(w, 201, tokenResponse)
}