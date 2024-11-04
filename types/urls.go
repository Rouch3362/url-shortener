package types

import (
	"regexp"

	"github.com/joho/godotenv"
	"github.com/lithammer/shortuuid/v4"
)

type CreateUrlRequest struct {
	Url string `json:"url"`
}

// validating the URL user enters in the CreateUrlRequest struct
func (c *CreateUrlRequest) Validator() string {
	if len(c.Url) == 0 {
		return "url field can not be empty"
	}

	urlRegex := `^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`

	isValid , _ := regexp.MatchString(urlRegex, c.Url)

	if !isValid {
		return "the url is not valid"
	}

	return ""
}



type DBCreateUrlRequest struct {
	UserId  	int  
	LongUrl 	string
	ShortUrl	string
}


func (d *DBCreateUrlRequest) CreateUrl() {
	// generating short uuid
	uuid := shortuuid.New()

	// getting the base URL for adding short uuid to it
	
	env,_ := godotenv.Read(".env")

	d.ShortUrl = env["W_ADDR"]+uuid
}
