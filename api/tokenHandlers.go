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


	// saving generated tokens into database
	refreshTokenData := &types.TokenDBRequest{
		RefreshToken: tokenResponse.RefreshToken,
		UserId: userFromDB.Id,
		ExpiresAt: cmd.ExpireationTime(false),
	} 

	a.DB.SaveToken(refreshTokenData)

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

	// checking if refresh token that user sent with request is valid and exists on database or not
	isRefreshTokenValid := a.DB.DoesRefreshTokenExists(refreshInstance.RefreshToken)


	if !isRefreshTokenValid {
		message := types.ErrorMessage{Message: "refresh token is not valid"}
		cmd.JsonGenerator(w, 401, message)
		return
	}

	// decoding jwt token to get user info's and not to reach for database calls
	userCredentials , _, err := cmd.VerifyJWTToken(refreshInstance.RefreshToken, false)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, 401, message)
		return
	}

	// generating new tokens for user
	tokenResponse := cmd.GenerateAuthTokens(userCredentials)


	// saving new generated tokens into database
	refreshTokenData := &types.TokenDBRequest{
		RefreshToken: tokenResponse.RefreshToken,
		UserId: userCredentials.Id,
		ExpiresAt: cmd.ExpireationTime(false),
	} 

	a.DB.SaveToken(refreshTokenData)


	cmd.JsonGenerator(w, 201, tokenResponse)
}