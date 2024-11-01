package types


// the type for showing error messages using json
type ErrorMessage struct {
	Message   	string `json:"message"`
}


// a struct for showing user columns
type UserResponse struct {
	Id 			int	   `json:"id"`
	Username	string `json:"username"`
	CreatedAt 	string	`json:"created_at"`
}