package x

import (
	"encoding/json"
	"net/url"

	"github.com/44za12/mailsleuth/internal/requestor"
)

type response struct {
	Valid bool   `json:"valid"`
	Msg   string `json:"msg"`
	Taken bool   `json:"taken"`
}

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	twitterUrl := "https://api.x.com/i/users/email_available.json"
	data := url.Values{}
	data.Set("email", email)
	URL, _ := url.Parse(twitterUrl + "?" + data.Encode())
	err := requestor.GET(URL)
	if err != nil {
		return false, err
	}
	var response response
	err = json.Unmarshal([]byte(requestor.Response.Body), &response)
	if err != nil {
		return false, err
	}
	if !response.Taken {
		return false, nil
	}
	return true, nil
}
