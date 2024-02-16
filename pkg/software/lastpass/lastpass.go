package lastpass

import (
	"net/url"

	"github.com/44za12/mailsleuth/internal/requestor"
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	getUrl, _ := url.Parse("https://lastpass.com/create_account.php?username=" + email + "&mistype=1&skipcontent=1&check=avail")
	err := requestor.GET(getUrl)
	if err != nil {
		return false, err
	}
	if requestor.Response.Body == "no" {
		return true, nil
	} else if requestor.Response.Body == "ok" {
		return false, nil
	}
	return false, nil
}
