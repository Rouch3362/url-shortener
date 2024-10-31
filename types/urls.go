package types

import "regexp"

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