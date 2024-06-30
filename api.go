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
	subRouter.HandleFunc("/url/{uuid}", a.DeleteUrl).Methods("DELETE").Name("middleware:CheckIfUserLoggedin")
	subRouter.HandleFunc("/url/{uuid}", a.UpdateUrl).Methods("PATCH","PUT").Name("middleware:CheckIfUserLoggedin")
	// a diffrent route for redirecting users to the actual url
	router.HandleFunc("/urls/{uuid}" , a.GetUrlHandler).Methods("GET")
	
	// users routes
	subRouter.HandleFunc("/user/register", a.CreateUserHandler).Methods("POST")
	subRouter.HandleFunc("/users/{username}", a.GetUserByUsernameHandler).Methods("GET")
	subRouter.HandleFunc("/users/{username}", a.UpdateUserHandler).Methods("PUT" , "PATCH").Name("middleware:CheckIfUserLoggedin")
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
	json.NewDecoder(r.Body).Decode(user)

	err := ValidateUserPayload(user.Username,user.Password)

	if err != nil {
		ErrorGenerator(w , err)
		return
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

	userCreden , tokenErr := VerifyToken(authToken)

	if tokenErr != nil {
		ErrorGenerator(w , tokenErr)
		return
	}

	// if the token provided was not access token
	if userCreden.Username == "" {
		ErrorGenerator(w , AccessTokenNeededError())
		return
	}
	// checks if user that request this operation is the same user of deleting user
	if accErr := CheckIfIsAccessOrRefresh(userCreden.Type , true); accErr != nil {
		ErrorGenerator(w , accErr)
		return
	}

	_, err := a.DB.DeleteUserDB(username)


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
	tokenString , jwtErr := CreateJWT(userExists,false)
	refreshString, refreshErr := CreateJWT(userExists,true)
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
		ErrorGenerator(w , RequiredFieldsError([]string{"refresh"}))
		return
	}
	// verfy if provided token in payload
	userCreden , err := VerifyToken(refreshRequest.Refresh)

	if err != nil {
		ErrorGenerator(w , err)
		return
	}


	if refErrr := CheckIfIsAccessOrRefresh(userCreden.Type , false); refErrr != nil {
		ErrorGenerator(w , refErrr)
		return
	}
	

	// get user by its id
	user , usrEr := a.DB.GetUserByIDDB(userCreden.UserId)

	// if user not found
	if usrEr != nil {
		ErrorGenerator(w , usrEr)
		return
	}
	// generating new tokens
	accessTokenString , accErr := CreateJWT(user,false)
	refreshTokenString, refErr := CreateJWT(user,true)

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


	userCreden , tokenErr := VerifyToken(authToken)

	// check if token is set
	if tokenErr != nil {
		ErrorGenerator(w , tokenErr)
		return
	}

	if accErr := CheckIfIsAccessOrRefresh(userCreden.Type , true); accErr != nil {
		ErrorGenerator(w , accErr)
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


func (a *APIServer) UpdateUserHandler(w http.ResponseWriter , r *http.Request) {
	username := mux.Vars(r)["username"]

	authToken := r.Header.Get("Authorization")

	userCreden , tokenErr := VerifyToken(authToken)

	if tokenErr != nil {
		ErrorGenerator(w , tokenErr)
		return
	}

	if accErr := CheckIfIsAccessOrRefresh(userCreden.Type , true); accErr != nil {
		ErrorGenerator(w, accErr)
		return
	}

	// check the actual user credential not those saved in token because after updating username the username in token wont't update since new token request
	user, _, err := a.DB.GetUserByUsernameDB(username)

	if err != nil {
		ErrorGenerator(w , err)
		return
	}

	if userCreden.UserId !=  user.ID{
		ErrorGenerator(w , AccessDeniedError())
		return
	}

	userPayload := &UserUpdateRequest{}

	json.NewDecoder(r.Body).Decode(userPayload)


	if userPayload.Username == "" {
		ErrorGenerator(w, RequiredFieldsError([]string{"username"}))
		return
	}

	// we don't want error value because we handled the not found error earlier
	response, _ := a.DB.UpdateUserDB(userPayload.Username,username)
	
	if err != nil {
		ErrorGenerator(w , err)
		return
	}


	JsonGenerator(w,http.StatusOK,response)

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

func (a *APIServer) UpdateUrl(w http.ResponseWriter , r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	authToken := r.Header.Get("Authorization")

	userCreden , tokenErr := VerifyToken(authToken)

	if tokenErr != nil {
		ErrorGenerator(w , tokenErr)
		return
	}


	// checks if access token provided not refresh token
	if accErr := CheckIfIsAccessOrRefresh(userCreden.Type , true); accErr != nil {
		ErrorGenerator(w , accErr)
		return
	}

	// decode user payload to a go struct
	urlRequest := &CreateUrlRequest{}

	json.NewDecoder(r.Body).Decode(urlRequest)
	// checks if url field is in payload
	validateErr := ValidateUrlPayload(urlRequest)

	if validateErr != nil {
		ErrorGenerator(w , validateErr)
		return
	}

	// load path prefix variable from env
	PATH_PREFIX := LoadEnvVariable("W_ADDR")
	// make url in accepted format
	short_url 	:= PATH_PREFIX+uuid
	// get users of that url and checks if url is exists
	url,urlErr := a.DB.GetUrlByShortUrl(short_url)
	// if url does not exist
	if urlErr != nil {
		ErrorGenerator(w , urlErr)
		return
	}
	// checks if the user requested to update is the same user that created url
	if url.User.ID != userCreden.UserId {
		ErrorGenerator(w , &Error{"Access Denied." , http.StatusForbidden})
		return
	}

	
	// updating url. at this situation we don't want the Error because we handled not found error in above codes.
	response , _ := a.DB.UpdateUrlDB(urlRequest.Url,url.ID)
	
	JsonGenerator(w , http.StatusOK , response)
}

func (a *APIServer) DeleteUrl(w http.ResponseWriter , r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	// token verification
	authToken := r.Header.Get("Authorization")

	userCreden , tokenErr := VerifyToken(authToken)

	if tokenErr != nil {
		ErrorGenerator(w , tokenErr)
		return
	}
	// ================================================


	// checks if access token provided not refresh token
	if accErr := CheckIfIsAccessOrRefresh(userCreden.Type , true); accErr != nil {
		ErrorGenerator(w , accErr)
		return
	}

	// load path prefix variable from env
	PATH_PREFIX := LoadEnvVariable("W_ADDR")
	// male url in accepted format
	short_url 	:= PATH_PREFIX+uuid
	// get users of that url and checks if url is exists
	url,urlErr := a.DB.GetUrlByShortUrl(short_url)
	// if url does not exist
	if urlErr != nil {
		ErrorGenerator(w , urlErr)
		return
	}
	// checks if the user requested to delete is the same user that created url
	if url.User.ID != userCreden.UserId {
		ErrorGenerator(w , AccessDeniedError())
		return
	}

	
	// deleting url
	a.DB.DeleteUrlDB(url.ID)
	
	JsonGenerator(w , http.StatusOK , struct{Message string}{fmt.Sprintf("%s Deleted successfully.", short_url)})
}




