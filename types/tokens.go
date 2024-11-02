package types


// a peroper format of token for showing to clients
type Token struct {
	AcccessToken 	string `json:"access_token"`
	RefreshToken 	string `json:"refresh_token"`	
}


type RefreshTokenRequest struct {
	RefreshToken 	string  `json:"refresh_token"`
}

type TokenDBRequest struct {
	AccessToken 	string
	RefreshToken 	string
	ExpiresAt		int64
	UserId 			int
}

func (r *RefreshTokenRequest) Validate() string{
	if r.RefreshToken == "" {
		return "refresh_token field is required."
	}

	return ""
}