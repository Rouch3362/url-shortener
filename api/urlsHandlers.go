package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Rouch3362/url-shortener/cmd"
	"github.com/Rouch3362/url-shortener/types"
	"github.com/gorilla/mux"
)

// handling POST requests for shorting an URL
func (a *APIServer) createUrlsHandler(w http.ResponseWriter, r *http.Request) {
	UrlRequest := &types.CreateUrlRequest{}

	authToken := r.Header.Get("Authorization")


	err := json.NewDecoder(r.Body).Decode(UrlRequest)

	if err != nil {
		log.Fatal(err)
	}
	// validate the URL field that user entered
	validationError := UrlRequest.Validator()

	// shows the proper error message to user if the URL is not valid 
	if validationError != "" {
		message := types.ErrorMessage{Message: validationError}
		cmd.JsonGenerator(w, 400, message)
		return
	}

	userCredentials , _, _ := cmd.VerifyJWTToken(authToken,false)

	// creating an instance for url
	urlInstance := types.DBCreateUrlRequest{
		UserId: userCredentials.Id,
		LongUrl: UrlRequest.Url,
	}
	// makes an uuid for saved long URL and saves that to a field called short URL
	urlInstance.CreateUrl()


	err = a.DB.CreateUrlDB(&urlInstance)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, http.StatusUnauthorized, message)
		return
	}

	cmd.JsonGenerator(w , 200 , urlInstance)
}


func (a *APIServer) getUrlHandler(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)

	urlId := pathVars["id"]

	originlaURL, err := a.DB.GetURL(urlId)

	if err != nil {
		message := types.ErrorMessage{Message: err.Error()}
		cmd.JsonGenerator(w, http.StatusNotFound, message)
		return
	}

	urlResponse := types.URLResponse{OriginalURL: originlaURL}

	cmd.JsonGenerator(w, http.StatusOK, urlResponse)
}




func (a *APIServer) deleteUrlHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	urlId	:= routeVars["id"]

	getUrlUser, err := a.DB.GetURLObject(urlId)

	if err != nil {
		message := types.ErrorMessage{Message:  err.Error()}
		cmd.JsonGenerator(w, http.StatusNotFound, message)
		return
	}

	usernameFromJWT := r.Context().Value(types.CtxKey)
	fmt.Println(usernameFromJWT, getUrlUser.User)
	if getUrlUser.User != usernameFromJWT {
		message := types.ErrorMessage{Message: "you are not allowed to do this"}
		cmd.JsonGenerator(w, http.StatusMethodNotAllowed, message)
		return
	}


	a.DB.DeleteURL(urlId)

	message := types.ErrorMessage{Message: "URL deleted successfully"}
	cmd.JsonGenerator(w, http.StatusOK, message)
}