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
	tokenString , jwtErr := CreateJWT(user)

	if jwtErr != nil {
		log.Fatal(jwtErr)
	}

	token := &JwtToken{Access: tokenString}

	JsonGenerator(w , http.StatusOK , token)

}



