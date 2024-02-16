package nimble

import (
	"net/url"
	"strings"

	"github.com/44za12/mailsleuth/internal/requestor"
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	getUrl, _ := url.Parse("https://www.nimble.com/lib/register.php?email=" + email)
	err := requestor.GET(getUrl)
	if err != nil {
		return false, err
	}
	if requestor.Response.Body == "true" {
		return false, nil
	} else if strings.Contains(requestor.Response.Body, "email is already registered") {
		return true, nil
	}
	return false, nil
}
