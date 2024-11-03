package api

import (
	"net/http"
	"github.com/Rouch3362/url-shortener/cmd"
	"github.com/Rouch3362/url-shortener/types"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extracting Authorization header value
		authToken := r.Header.Get("Authorization")

		// if auth token does not exists
		if authToken == "" {
			message := types.ErrorMessage{Message: "Authorization token is required for this action."}
			cmd.JsonGenerator(w, 401, message)
			return
		}

		// validating and checking if user used access token not refresh token
		_,isAccessToken, err := cmd.VerifyJWTToken(authToken, true)

		if err != nil {
			message := types.ErrorMessage{Message: err.Error()}
			cmd.JsonGenerator(w, 401, message)
			return
		}

		// if user used refersh token instead of acccess token
		if !isAccessToken {
			message := types.ErrorMessage{Message: "use access token instead of refresh token."}
			cmd.JsonGenerator(w, 401, message)
			return
		}

		// if everything is fine the app goes forward
		next.ServeHTTP(w , r)
	})
}