package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type APIServer struct {
	Addr string
	DB   Storage
}

func NewApiServer(addr string, db Storage) *APIServer {
	apiInstance := &APIServer{
		Addr: addr,
		DB:   db,
	}

	return apiInstance
}

func (a *APIServer) Run() error {
	router := mux.NewRouter()
	// router.Use(CheckIfUserLoggedin)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.Use(CheckIfUserLoggedin)
	
	// test route 
	subRouter.HandleFunc("/hello" , a.SayHello).Methods("GET")
	
	// name specified here can determined in middleware for running it if this route is protected
	subRouter.HandleFunc("/url", a.CreateUrlHandler).Methods("POST").Name("middleware:CheckIfUserLoggedin")
	
	// a diffrent route for redirecting users to the actual url
	router.HandleFunc("/urls/{uuid}" , a.GetUrlHandler).Methods("GET")
	
	// users routes
	subRouter.HandleFunc("/user", a.CreateUserHandler).Methods("POST")
	subRouter.HandleFunc("/users/{username}", a.GetUserByUsernameHandler).Methods("GET")
	subRouter.HandleFunc("/users/{username}", a.DeleteUserHandler).Methods("DELETE").Name("middleware:CheckIfUserLoggedin")
	subRouter.HandleFunc("/users/{username}/urls", a.GetUsersUrlHandler).Methods("GET")
	subRouter.HandleFunc("/user/login" , a.LoginHandler).Methods("POST")
	subRouter.HandleFunc("/user/login/refresh" , a.RefreshTokenHandler).Methods("POST")
	
	
	err := http.ListenAndServe(a.Addr , router)


	return err

}

// testing api is alive
func (a *APIServer) SayHello(w http.ResponseWriter , r *http.Request) {
	helloStruct := struct{Message string}{"hello world"}

	JsonGenerator(w , http.StatusOK , helloStruct)

} 


func (a *APIServer) GetUserByUsernameHandler(w http.ResponseWriter , r *http.Request) {
	username := mux.Vars(r)["username"]
	user , _ , err := a.DB.GetUserByUsernameDB(username)

	if err != nil {
		ErrorGenerator(w, err)
		return
	}

	JsonGenerator(w, http.StatusOK , user)

}




// creating new user by POST method 
func (a *APIServer) CreateUserHandler(w http.ResponseWriter , r *http.Request) {
	// creting an empty instance for accesing its method
	user := &UserRequest{}

	// decode user payload to User struct
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		log.Fatal(err)
	} 

	// create a user and default values of createdAt field
	u , userErr  := user.CreateUser()

	// sends error if payload is not valid
	if userErr != nil {
		ErrorGenerator(w, userErr)
		return
	}

	// save created user to database
	createdUser , DBerr := a.DB.CreateUserDB(u)

	// sends error if user is already exist
	if DBerr != nil {
		ErrorGenerator(w , DBerr)
		return
	}
	
	// if everything is okay returns created user
	JsonGenerator(w , http.StatusCreated , createdUser)
}	

func (a *APIServer) DeleteUserHandler(w http.ResponseWriter , r *http.Request) {
	username := mux.Vars(r)["username"]


	authToken := r.Header.Get("Authorization")

	userCreden , tokenErr := VerifyToken(authToken , true)

	if tokenErr != nil {
		ErrorGenerator(w , tokenErr)
		return
	}

	// if the token provided was not access token
	if userCreden.Username == "" {
		ErrorGenerator(w , &Error{"for deleting user you must enter access token not refresh token", http.StatusBadRequest})
		return
	}
	// checks if user that request this operation is the same user of deleting user
	if userCreden.Username != username && userCreden.UserId == 0 {
		ErrorGenerator(w , &Error{"Access Denied." , http.StatusForbidden})
		return
	}

	_, err := a.DB.DelteUserDB(username)


	if err != nil {
		ErrorGenerator(w , err)
		return
	}

	JsonGenerator(w , http.StatusOK, struct{Message string}{fmt.Sprintf("%s user deleted successfully.",username)})
}

func (a *APIServer) LoginHandler(w http.ResponseWriter , r *http.Request) {
	user := &UserRequest{}

	json.NewDecoder(r.Body).Decode(user)

	// validate fields user enters
	validateErr := ValidateUserPayload(user.Username, user.Password)

	if validateErr != nil {
		ErrorGenerator(w , validateErr)
		return
	}

	// check if user with entered username exists
	userExists, userPassword ,notFoundErr := a.DB.GetUserByUsernameDB(user.Username)

	if notFoundErr != nil {
		ErrorGenerator(w, notFoundErr)
		return
	}

	// check if entered password is valid for the user
	passErr := IsPasswordValid(userPassword,user.Password)

	if passErr != nil {
		ErrorGenerator(w , passErr)
		return
	}

	// creating jwt token
	tokenString , jwtErr := CreateJWT(userExists)
	refreshString, refreshErr := CreateRefreshToken(userExists.ID)
	// create one row in db for refresh token
	if err := a.DB.CreateRefreshTokenDB(userExists.ID , refreshString); err != nil {
		log.Fatal(err)
	}
	if jwtErr != nil || refreshErr != nil{
		log.Fatal(jwtErr , refreshErr)
	}
	// creates an instance of jwt's results
	token := &JwtToken{Access: tokenString, Refresh: refreshString}

	JsonGenerator(w , http.StatusOK , token)

}


func (a *APIServer) RefreshTokenHandler(w http.ResponseWriter , r *http.Request) {
	// an instance for RefreshRequest for decoding the values of request
	refreshRequest := &RefershTokenRequest{}
	

	json.NewDecoder(r.Body).Decode(refreshRequest)

	// if refresh field was empty
	if refreshRequest.Refresh == "" {
		ErrorGenerator(w , &Error{"refresh field not provided." , http.StatusBadRequest})
		return
	}
	// verfy provided token
	userCreden , err := VerifyToken(refreshRequest.Refresh,true)

	if err != nil {
		ErrorGenerator(w , err)
		return
	}
	// get user by its id
	user , usrEr := a.DB.GetUserByIDDB(userCreden.UserId)

	// if user not found
	if usrEr != nil {
		/* checks if user credential returned from verify token is
		a refresh token credential or a access token if the credentials is 
		for access token this returns an err */

		if userCreden.Username != "" {
			ErrorGenerator(w , &Error{
				"for refreshing token you must enter the refresh token not access token",
				http.StatusBadRequest })
			return
		}
		ErrorGenerator(w , usrEr)
		return
	}
	// generating new tokens
	accessTokenString , accErr := CreateJWT(user)
	refreshTokenString, refErr := CreateRefreshToken(user.ID)

	if accErr != nil || refErr != nil {
		log.Fatal(accErr , refErr)
	}
	// deleted the used refresh token
	refDelErr := a.DB.DeleteRefreshTokenDB(refreshRequest.Refresh)
	
	// if row is black listed or once used return an error
	if refDelErr != nil {
		ErrorGenerator(w , refDelErr)
		return
	}
	// creates an instance for returning new tokens to user
	token := &JwtToken{Access: accessTokenString , Refresh: refreshTokenString}
	// creates new row to database for new token
	if err := a.DB.CreateRefreshTokenDB(user.ID,refreshTokenString); err != nil {
		log.Fatal(err)
	}

	JsonGenerator(w,http.StatusOK , token)
	
}



func (a *APIServer) CreateUrlHandler(w http.ResponseWriter , r *http.Request) {
	UrlRequest :=  &CreateUrlRequest{}
	// cleaned token from  middleware
	authToken := r.Header.Get("Authorization")


	userCreden , tokenErr := VerifyToken(authToken , true)

	// check if token is set
	if tokenErr != nil {
		ErrorGenerator(w , tokenErr)
		return
	}

	userExist , _ , usrErr := a.DB.GetUserByUsernameDB(userCreden.Username)


	if usrErr != nil {
		ErrorGenerator(w , usrErr)
		return
	}


	json.NewDecoder(r.Body).Decode(UrlRequest)


	validateErr := ValidateUrlPayload(UrlRequest)

	if validateErr != nil {
		ErrorGenerator(w , validateErr)
		return
	}


	urlInstance := UrlRequest.CreateUrl(userExist.ID)

	createdUrl, urlErr := a.DB.CreateUrlDB(urlInstance)

	if urlErr != nil {
		ErrorGenerator(w , urlErr)
		return
	}

	JsonGenerator(w , http.StatusCreated , createdUrl)

}


func (a *APIServer) GetUsersUrlHandler(w http.ResponseWriter , r *http.Request) {
	
	username := mux.Vars(r)["username"]

	response , err := a.DB.GetUserUrlsDB(username)

	if err != nil {
		ErrorGenerator(w , err)
		return
	}

	JsonGenerator(w,http.StatusOK,response)

}


func (a *APIServer) GetUrlHandler(w http.ResponseWriter , r *http.Request) {
	PATH_PREFIX := LoadEnvVariable("W_ADDR")
	uuid := mux.Vars(r)["uuid"]

	// gets the url by its uuid from database
	url , err := a.DB.GetUrlByShortUrl(PATH_PREFIX+uuid)

	if err != nil {
		ErrorGenerator(w , err)
		return
	}

	// increase clicks by 1
	a.DB.IncreaseClicksUrlDB(url.NewUrl)

	// it will save on database but not shows the returned value after update in database because of less database calls its just a value showen to user 
	url.Clicks += 1

	JsonGenerator(w , http.StatusFound , url)
}






