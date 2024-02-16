package insightly

import (
	"net/url"
	"strings"

	"github.com/44za12/mailsleuth/internal/requestor"
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	postUrl, _ := url.Parse("https://accounts.insightly.com/signup/isemailvalid")
	data := make(map[string]string)
	data["emailaddress"] = email
	requestor.Request.Parameters = data
	err := requestor.POST(postUrl)
	if err != nil {
		return false, err
	}
	if requestor.Response.Body == "true" {
		return false, nil
	} else if strings.Contains(requestor.Response.Body, "An account exists for this address. Use another address or") {
		return true, nil
	}
	return false, nil
}
