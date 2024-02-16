package kommo

import (
	"encoding/json"
	"net/url"

	"github.com/44za12/mailsleuth/internal/requestor"
)

type response struct {
	Status string `json:"status"`
	Mail   string `json:"mail"`
}

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	postUrl, _ := url.Parse("https://www.kommo.com/account/check_login.php")
	data := make(map[string]string)
	data["LOGIN"] = email
	requestor.Request.Parameters = data
	err := requestor.POST(postUrl)
	if err != nil {
		return false, err
	}
	var response response
	err = json.Unmarshal([]byte(requestor.Response.Body), &response)
	if err != nil {
		return false, err
	}
	if response.Status == "free" {
		return false, nil
	} else if response.Status == "used" {
		return true, nil
	}
	return false, nil
}
