package api

import (
	"encoding/json"
	"net/http"
	"github.com/Rouch3362/url-shortener/cmd"
	"github.com/Rouch3362/url-shortener/types"
)


func (a *APIServer) createUserHandler(w http.ResponseWriter , r *http.Request) {
	// creating empty instance for filling with user payload	
	userPayload := &types.UserRequest{}
	
	// decoding user payload to userPayload
	json.NewDecoder(r.Body).Decode(userPayload)

	// checking for any errors in payload
	payloadError := userPayload.Validator()


	if payloadError != "" {
		errorMessage := types.ErrorMessage{Message: payloadError}
		cmd.JsonGenerator(w , 400 , errorMessage)
		return
	}

	// hashing user's password for security :)
	userPayload.HashPassword()

	// saving user into database
	result, _ := a.DB.CreateNewUser(userPayload)

	// returning result to the client
	cmd.JsonGenerator(w, 201, result)
}