package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)



func LoadEnvVariable(varName string) string {
	// loads .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
	// gets env variable value based on its name
	value := os.Getenv(varName)

	return value
}


func ValidateUserPayload(username , password string) *Error{
	// checks if the request for creating user has requried fields
	if username == "" || password == "" {
		return RequiredFieldsError([]string{"username","password"})
	// check if length of values is longer than 8 characters
	} else if len(username) < 8 || len(password) < 8 {
		return &Error{
			Message: "username and password must be longer than 8 characters.",
			Code: http.StatusBadRequest,
		}
	}

	return nil
}

func ValidateUrlPayload(url *CreateUrlRequest) *Error {
	
	if url.Url == "" {
		return RequiredFieldsError([]string{"url"})
	}
	
	UrlRegex := `^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`
	matched , _  := regexp.MatchString(UrlRegex,url.Url)

	if !matched {
		return &Error{
			Message: "the value of url field you entered is not an URL.",
			Code: http.StatusBadRequest,
		}
	}


	return nil
}

// a helper function for producing responses
func JsonGenerator(w http.ResponseWriter , statusCode int , value any) {
	w.Header().Add("Content-Type" , "application/json")
	
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(value)

	if err != nil {
		log.Fatal(err)
	}
}

// isRefreshToken argument determines if this generation is for a refresh token or not and if it is the exp data is more than access
func CreateJWT(user *UserResponse , isRefreshToken bool) (string , error) {
	expHour := time.Hour * 24
	tokenType := "access"
	if isRefreshToken {
		tokenType = "refresh"
		expHour = time.Hour * 48
	}

	// claims (fields) we want in our jwt token
	claims := jwt.MapClaims{
		"username": user.Username,
		"id": user.ID,
		"type": tokenType,
		// expires after 1 day of creation
		"exp": time.Now().Add(expHour).Unix(),
	}

	// creating token based on ecryption algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256 , claims)

	// loads jwt secret from env file	
	JWT_SECRET := LoadEnvVariable("JWT_SECRET")
	// getting generated access token
	tokenString , err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "" , err
	}


	return tokenString,nil
}

// verifies if token is valid and not expired
func VerifyToken(token string) (*VerifyTokenResult,*Error) {
	// loading jwt secret
	JWT_SECRET := LoadEnvVariable("JWT_SECRET")

	// parse token 
	parsedClaims, err := jwt.Parse(token , func(t *jwt.Token) (interface{}, error) {return []byte(JWT_SECRET) , nil})


	if err != nil {
		return nil, &Error{
			Message: "Token is not valid or expired.",
			Code: http.StatusUnauthorized,
		}
	}

	

	// get claims of the token
	claims := parsedClaims.Claims.(jwt.MapClaims)

	verifyTokenRes := VerifyTokenResult{}

	// convert token fields to a go struct
	verifyTokenRes.UserId = int(claims["id"].(float64))
	verifyTokenRes.Username = string(claims["username"].(string))
	verifyTokenRes.Type = string(claims["type"].(string))


	return &verifyTokenRes,nil
}


// checks if a password and its hash is the same
func IsPasswordValid(hash , password string) *Error{
	// if not the same returns error
	err := bcrypt.CompareHashAndPassword([]byte(hash) , []byte(password))

	if err != nil {
		return &Error{
			Message: "username or password is invalid",
			Code: http.StatusUnauthorized,
		}
	}

	return nil
}

// generates error based on Error struct
func ErrorGenerator(w http.ResponseWriter, error *Error) {
	JsonGenerator(w,error.Code , error)
}

func ExtractRawToken(token string) (string, *Error){
	// checks if token is entered
	if token != "" && (strings.Contains(token , "Bearer ") || strings.Contains(token,"bearer")) {
		token = token[len("Bearer "):]
	}

	// checks if tokne is not empty
	if token == "" {
		return "" , NotAuthorizedError()
	}

	return token,nil
}

/* checks if user credential returned from verify token is
a refresh token credential or a access token if the credentials is 
for access token this returns an err or opsite */
// the checkAccess argument is for knowing if the called is for a refresh token or an access token
func CheckIfIsAccessOrRefresh(tokenType string , checkAccess bool) *Error {
	if checkAccess && tokenType == "refresh" {
		return AccessTokenNeededError()
	} else if !checkAccess && tokenType == "access" {
		return RefreshTokenNeededError()
	}
	return nil
}



// error's
func AccessDeniedError() *Error {
	return &Error{"Access Denied." , http.StatusForbidden}
}

func RequiredFieldsError(fields []string) *Error {
	
	textFields := ""

	for _,field := range fields {
		textFields += fmt.Sprintf("%s, " , field)
	}

	return &Error{fmt.Sprintf("%s field(s) are required." , textFields) , http.StatusBadRequest}
}

func NotFoundError(entity string) *Error {
	return &Error{fmt.Sprintf("%s not found." , entity) , http.StatusNotFound}
}

func BlackListedTokenError() *Error {
	return &Error{"token is not valid any more (black listed)." , http.StatusUnauthorized}
}

func AccessTokenNeededError() *Error {
	return &Error{"for this action you must use access token not refresh token." , http.StatusUnauthorized}
}

func RefreshTokenNeededError() *Error {
	return &Error{"for this action you must use refresh token not access token.", http.StatusUnauthorized}
}

func NotAuthorizedError() *Error {
	return &Error{"authorization token not provided." , http.StatusUnauthorized}
}
 
