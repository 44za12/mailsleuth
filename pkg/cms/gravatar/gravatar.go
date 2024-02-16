package gravatar

import (
	"net/url"

	"github.com/44za12/mailsleuth/internal/requestor"
	"github.com/44za12/mailsleuth/internal/utils"
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	hashedEmail := utils.HashString(email)
	getUrl, _ := url.Parse(`https://en.gravatar.com/` + hashedEmail + `.json`)
	err := requestor.GET(getUrl)
	if err != nil {
		return false, err
	}
	if requestor.Response.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}
