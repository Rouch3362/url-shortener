package types


// a peroper format of token for showing to clients
type Token struct {
	AcccessToken 	string `json:"access_token"`
	RefreshToken 	string `json:"refresh_token"`	
}