package anydo

import (
	"encoding/json"
	"net/url"

	"github.com/44za12/mailsleuth/internal/requestor"
)

var (
	POST_URL, _ = url.Parse("https://sm-prod2.any.do/check_email")
)

type response struct {
	Exists bool `json:"user_exists"`
}

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	rawData := `{"email":"` + email + `"}`
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
	return response.Exists, nil
}
