package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rouch3362/url-shortener/cmd"
	"github.com/Rouch3362/url-shortener/types"
)



// handling POST requests for shorting an URL 
func createUrlsHandler(w http.ResponseWriter, r *http.Request) {
	UrlRequest := &types.CreateUrlRequest{}

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
	}
}