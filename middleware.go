package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func CheckIfUserLoggedin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter , r *http.Request) {
		if mux.CurrentRoute(r).GetName() == "middleware:CheckIfUserLoggedin" {
			// get the header attribute with this format Bearer Token
			authorizationToken := r.Header.Get("Authorization")
			// this removes Bearer from the token
			authToken , err := ExtractRawToken(authorizationToken)

			// checks if token is valid and if not returns the error
			if err != nil {
				ErrorGenerator(w , err)
				return
			}

			// set Authorization header to raw token without Bearer keyword
			r.Header.Set("Authorization",authToken)
		}
		// if everthing is good goes to next handler
		next.ServeHTTP(w , r)
	})
}