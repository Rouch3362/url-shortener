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



func GenerateJWTToken(payload *types.UserResponse, isAccessToken bool) string {
	// expires at 12 hours later of generation
	expHour := time.Hour * 12 

	// the default type of token is access token
	tokenType := "access"

	// changes based on argument entered
	if !isAccessToken {
		tokenType = "refresh"
		expHour *= 2
	}

	// payload of token
	claims := jwt.MapClaims{
		"id": payload.Id,
		"username": payload.Username,
		"exp": time.Now().Add(expHour).Unix(),
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


func DecodeJWTToken(token string) (*types.UserResponse, error) {
	jwtSecret := ReadEnvVar("JWT_SECRET")
	claims := jwt.MapClaims{}

	_ , err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	username := claims["username"].(string)
	userId :=	claims["id"].(int)
	user := types.UserResponse{Username: username, Id: userId, CreatedAt: ""}


	return &user,nil
}