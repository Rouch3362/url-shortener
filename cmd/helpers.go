package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/Rouch3362/url-shortener/types"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func ReadEnvVar(varName string) string {
	env, err := godotenv.Read(".env")

	if err != nil{
		log.Fatal(err)
	}

	value := env[varName]

	return value 
}



func JsonGenerator(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Add("Content-Type" , "application/json")

	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(value)

	if err != nil {
		log.Fatal(err)
	}
}

func ExpireationTime(isAccessToken bool) int64 {
	// expires at 12 hours later of generation
	expHour := time.Hour * 24 * 7
	// changes based on argument entered
	if !isAccessToken {
		expHour *= 2
	}

	return time.Now().Add(expHour).Unix()
}

func GenerateJWTToken(payload *types.UserResponse, isAccessToken bool) string {
	 

	// the default type of token is access token
	tokenType := "access"

	// changes based on argument entered
	if !isAccessToken {
		tokenType = "refresh"
	}

	// payload of token
	claims := jwt.MapClaims{
		"userId": payload.Id,
		"username": payload.Username,
		"exp": ExpireationTime(isAccessToken),
		"type": tokenType,
	}

	// generating token 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// getting secret of tokens
	jwtSecret := ReadEnvVar("JWT_SECRET")


	// extracting token string from token
	tokenString, err := token.SignedString([]byte(jwtSecret))


	if err != nil {
		log.Fatal(err)
	}


	return tokenString
}

// generating refresh token and access token with each other and returning them at the same time
func GenerateAuthTokens(payload *types.UserResponse) *types.Token {
	access  := GenerateJWTToken(payload, true)
	refresh := GenerateJWTToken(payload, false)

	result := types.Token{
		AcccessToken: access,
		RefreshToken: refresh,
	}

	return &result
}

// extracting user's data from JWT tokens
func DecodeJWTToken(token string) (*types.UserResponse, error) {
	// getting jwt secret from env file
	jwtSecret := ReadEnvVar("JWT_SECRET")
	claims := jwt.MapClaims{}

	// parse jwt tokens
	_ , err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// extracting username and userID from jwt
	username := claims["username"].(string)
	userId 	 :=	int(claims["userId"].(float64))
	user := types.UserResponse{Username: username, Id: userId, CreatedAt: ""}


	return &user,nil
}