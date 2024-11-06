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

type URLResponse struct{
	OriginalURL string `json:"original_url"`	
}

type URL struct {
	Id				int 	`json:"id"`
	OriginalURL 	string 	`json:"original_url"`
	ShortURL 		string	`json:"short_url"`
	Clicks			int		`json:"clicks"`
	CreatedAt 		string	`json:"created_at"`
}

type URLObject struct {
	Id				int 	`json:"id"`
	OriginalURL 	string 	`json:"original_url"`
	ShortURL 		string	`json:"short_url"`
	User			string	`json:"user"`
	Clicks			int		`json:"clicks"`
	CreatedAt 		string	`json:"created_at"`
}

type UserURLsResponse struct {
	Id			int    `json:"id"`
	Username 	string `json:"username"`
	CreatedAt	string `json:"created_at"`
	Urls		[]URL  `json:"urls"`
}