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
	hashedPassword, err := a.DB.GetUserPassword(userPayload.Username)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, http.StatusNotFound, message)
		return
	}

	// checking if hash and password in same
	isPasswordValid := userPayload.ComparePassword(hashedPassword)

	if !isPasswordValid {
		message := types.ErrorMessage{Message: "username or password is not correct"}
		cmd.JsonGenerator(w , 401, message)
		return
	}

	// getting user from database
	userFromDB, _ := a.DB.GetUserByUsername(userPayload.Username)
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

	// -------------------------------------
	// this block of code makes sure that the user who's refreshing the token is the same user who's logged in and checks it with Authorization header
	authToken := r.Header.Get("Authorization")

	// decoding jwt token to get user info's and not to reach for database calls
	userCredentialsFromRefersh , _, err := cmd.VerifyJWTToken(refreshInstance.RefreshToken, false)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, 401, message)
		return
	}
	// extracting user info for checking if the user who is requesting refresh token is the same use who is logged in
	parsedAuthToken , _ , _ := cmd.VerifyJWTToken(authToken, false)

	if parsedAuthToken.Id != userCredentialsFromRefersh.Id {
		message := types.ErrorMessage{Message: "you are not allowed to do this (the user of access token is not the same user of refresh token)"}
		cmd.JsonGenerator(w, 403, message)
		return
	}
	// ----------------------------

	// checking if refresh token that user sent with request is valid and exists on database or not
	isRefreshTokenValid := a.DB.DoesRefreshTokenExists(refreshInstance.RefreshToken)


	if !isRefreshTokenValid {
		message := types.ErrorMessage{Message: "refresh token is not valid"}
		cmd.JsonGenerator(w, 401, message)
		return
	}

	

	// generating new tokens for user
	tokenResponse := cmd.GenerateAuthTokens(userCredentialsFromRefersh)


	// saving new generated tokens into database
	refreshTokenData := &types.TokenDBRequest{
		RefreshToken: tokenResponse.RefreshToken,
		UserId: userCredentialsFromRefersh.Id,
		ExpiresAt: cmd.ExpireationTime(false),
	} 

	a.DB.SaveToken(refreshTokenData)


	cmd.JsonGenerator(w, 201, tokenResponse)
}