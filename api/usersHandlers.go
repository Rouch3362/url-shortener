package api

import (
	"encoding/json"
	"net/http"
	"github.com/Rouch3362/url-shortener/cmd"
	"github.com/Rouch3362/url-shortener/types"
	"github.com/gorilla/mux"
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
	result, err := a.DB.CreateNewUser(userPayload)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, http.StatusConflict, message)
		return
	}

	// returning result to the client
	cmd.JsonGenerator(w, 201, result)
}



func (a *APIServer) GetUser(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	username := pathVars["username"]

	result, err := a.DB.GetUserURLs(username)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, http.StatusNotFound, message)
		return
	}

	cmd.JsonGenerator(w, http.StatusOK, result)

}


func (a *APIServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(types.CtxKey).(string)

	if !ok {
		authToken := r.Header.Get("Authorization")
		userInfo, _, _ := cmd.VerifyJWTToken(authToken, false)
		username = userInfo.Username
	}

	err := a.DB.DeleteUserDB(username)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, http.StatusOK, message)
		return
	}

	message := types.ErrorMessage{Message: "user deleted successfully"}
	cmd.JsonGenerator(w, http.StatusOK, message)
}