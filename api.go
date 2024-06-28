package main

import (
	"encoding/json"
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

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	subRouter.HandleFunc("/hello" , a.SayHello).Methods("GET")
	subRouter.HandleFunc("/user", a.CreateUserHandler).Methods("POST")
	subRouter.HandleFunc("/user/{username}", a.GetUserByUsernameHandler).Methods("GET")
	subRouter.HandleFunc("/user/login" , a.LoginHandler).Methods("POST")
	subRouter.HandleFunc("/user/login/refresh" , a.RefershTokenHandler).Methods("POST")
	
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
	user , err := a.DB.GetUserByUsernameDB(username)

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


func (a *APIServer) LoginHandler(w http.ResponseWriter , r *http.Request) {
	user := &LoginRequest{}

	json.NewDecoder(r.Body).Decode(user)

	// validate fields user enters
	validateErr := ValidatePayload(user.Username, user.Passwrod)

	if validateErr != nil {
		ErrorGenerator(w , validateErr)
		return
	}

	// check if user with entered username exists
	userExists, notFoundErr := a.DB.GetUserByUsernameDB(user.Username)

	if notFoundErr != nil {
		ErrorGenerator(w, notFoundErr)
		return
	}

	// check if entered password is valid for the user
	passErr := IsPasswordValid(userExists.Password,user.Passwrod)

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


func (a *APIServer) RefershTokenHandler(w http.ResponseWriter , r *http.Request) {
	// an instance for RefreshRequest for decoding the values of request
	refreshRequest := &RefershTokenRequest{}
	

	json.NewDecoder(r.Body).Decode(refreshRequest)

	// if refresh field was empty
	if refreshRequest.Refresh == "" {
		ErrorGenerator(w , &Error{"refresh field not provided." , http.StatusBadRequest})
		return
	}
	// verfy provided token
	userId , err := VerifyToken(refreshRequest.Refresh,true)

	if err != nil {
		ErrorGenerator(w , err)
		return
	}
	// get user by its id
	user , usrEr := a.DB.GetUserByIDDB(userId)

	// if user not found
	if usrEr != nil {
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


