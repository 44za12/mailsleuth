package hubspot

import (
	"encoding/json"
	"net/url"

	"github.com/44za12/mailsleuth/internal/requestor"
	"github.com/44za12/mailsleuth/internal/utils"
)

var (
	POST_URL, _ = url.Parse("https://app.hubspot.com/api/login-api/v1/login")
)

type response struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	Correlationid string `json:"correlationId"`
}

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	rawData := `{"email":"` + email + `","password":"` + utils.RandomString(8) + `","rememberLogin":false}`
	requestor.Request.RawParameters = rawData
	err := requestor.POST(POST_URL)
	if err != nil {
		return false, err
	}
	var response response
	err = json.Unmarshal([]byte(requestor.Response.Body), &response)
	if err != nil {
		return false, err
	}
	if response.Status == "INVALID_USER" {
		return false, nil
	} else if response.Status == "INVALID_PASSWORD" {
		return true, nil
	}
	return false, nil
}
