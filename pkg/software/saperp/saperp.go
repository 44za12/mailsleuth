package saperp

import (
	"encoding/json"
	"net/url"

	"github.com/44za12/mailsleuth/internal/requestor"
)

type response struct {
	Employee        bool `json:"employee"`
	HasUID          bool `json:"hasUid"`
	LinkingRequired bool `json:"linkingRequired"`
	Shared          bool `json:"shared"`
}

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	URL, _ := url.Parse("https://core-api.account.sap.com/uid-core/employee/" + email + "/verify")
	err := requestor.GET(URL)
	if err != nil {
		return false, err
	}
	var response response
	err = json.Unmarshal([]byte(requestor.Response.Body), &response)
	if err != nil {
		return false, err
	}
	if !response.HasUID {
		return false, nil
	}
	return true, nil
}
