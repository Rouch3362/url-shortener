package main

import (
	"log"
	"time"

	"github.com/lithammer/shortuuid/v3"
	"golang.org/x/crypto/bcrypt"
)

// our custom error struct
type Error struct {
	Message string `json:"message"`
	Code  int    `json:"code"`
}

// wrote a method for our new custom error struct so we can use it for returning errors of functions
func (e *Error) Error() string {
	return e.Message
}

// a type for user response
type User struct {
	ID        int    	`json:"id" db:"id"`
	Username  string 	`json:"username" db:"username"`
	Password  string 	`json:"password" db:"password"`
	CreatedAt time.Time	`json:"createdAt" db:"created_at"`
}
type UserUpdateRequest struct {
	Username	string		`json:"username" db:"username"`
}
// a struct for creating new user
type UserRequest struct {
	Username  string 	`json:"username" db:"username"`
	Password  string 	`json:"password" db:"password"`
}
type UserResponse struct {
	ID        int    	`json:"id" db:"id"`
	Username  string 	`json:"username" db:"username"`
	CreatedAt time.Time	`json:"createdAt" db:"created_at"`
}

type UserUrlsResponse struct {
	ID        int    			`json:"id" db:"id"`
	Username  string 			`json:"username" db:"username"`
	Urls 	  []ShortedFormUrl  `json:"urls" db:"urls"`
	CreatedAt time.Time			`json:"createdAt" db:"created_at"`
}

type ShortedFormUrl struct {
	ID 			int 			`json:"id" db:"id"`
	OldUrl 		string			`json:"oldUrl" db:"old_url"`
	NewUrl 		string			`json:"newUrl" db:"new_url"`
	Clicks		int				`json:"clicks" db:"clicks"`
	CreatedAt	time.Time		`json:"createdAt" db:"created_at"`
}

// for showing users their tokens
type JwtToken struct {
	Access 	string 	`json:"access"`
	Refresh	string	`json:"refresh" db:"token"`
}

// type for refreshing expired tokens
type RefershTokenRequest struct {
	Refresh string	`json:"refresh" db:"token"`
}

// for filling when using verify token
type VerifyTokenResult struct {
	UserId 		int	
	Username 	string
	Type		string
}

// a struct for request to create new url
type CreateUrlRequest struct {
	Url string	`json:"url"`
}
// a struct for url response 
type Url struct {
	ID 			int 		`json:"id" db:"id"`
	User 		int			`json:"user" db:"user_id"`
	OldUrl 		string		`json:"oldUrl" db:"old_url"`
	NewUrl 		string		`json:"newUrl" db:"new_url"`
	Clicks		int				`json:"clicks" db:"clicks"`
	CreatedAt	time.Time	`json:"createdAt" db:"created_at"`
}
// a struct for nested data from database
type UrlReponse struct {
	ID 			int 			`json:"id" db:"id"`
	User 		UserResponse	`json:"user" db:"user_id"`
	OldUrl 		string			`json:"oldUrl" db:"old_url"`
	NewUrl 		string			`json:"newUrl" db:"new_url"`
	Clicks		int				`json:"clicks" db:"clicks"`
	CreatedAt	time.Time		`json:"createdAt" db:"created_at"`
}


// a method for creating an instancee of user before saving it to database
func (u *UserRequest) CreateUser() (*User, *Error) {
	validateErr := ValidateUserPayload(u.Username , u.Password)

	if validateErr != nil {
		return nil,validateErr
	}

	// hashing password
	hashPassword , hashErr := bcrypt.GenerateFromPassword([]byte(u.Password) , bcrypt.DefaultCost)
	
	if hashErr != nil {
		log.Fatal(hashErr)
	}
	user := &User{
		Username: u.Username, 
		Password: string(hashPassword),
		CreatedAt: time.Now().UTC(),
	}

	return user , nil
}


func (u *CreateUrlRequest) CreateUrl(userId int) *Url {
	// uses path prefix to add domain of api to uuid
	PATH_PREFIX := LoadEnvVariable("W_ADDR")

	// generating a uuid
	uuid := shortuuid.New()

	// create an instance for url
	return &Url{
		User: userId,
		OldUrl: u.Url,
		NewUrl: PATH_PREFIX+uuid,
		Clicks: 0,
		CreatedAt: time.Now().UTC(),
	}
}
