package zoho

import (
	"net/url"
	"strings"

	"github.com/44za12/mailsleuth/internal/requestor"
	"github.com/44za12/mailsleuth/internal/utils"
)

var (
	CSRF_URL, _ = url.Parse("https://accounts.zoho.in/signin?servicename=VirtualOffice&signupurl=https://www.zoho.com/mail/zohomail-pricing.html&serviceurl=https://mail.zoho.in")
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	err := requestor.GET(CSRF_URL)
	if err != nil {
		return false, err
	}
	token := requestor.GetCookie("iamcsr")
	data := map[string]string{
		"mode":        "primary",
		"cli_time":    utils.GetCurrentTimeStampAsStr(),
		"servicename": "VirtualOffice",
		"serviceurl":  "https://mail.zoho.in",
		"signupurl":   "https://www.zoho.com/mail/zohomail-pricing.html",
	}
	requestor.Request.Parameters = data
	requestor.Headers["X-ZCSRF-TOKEN"] = "iamcsrcoo=" + token
	REGISTGER_URL, _ := url.Parse("https://accounts.zoho.in/signin/v2/lookup/" + email)
	err = requestor.POST(REGISTGER_URL)
	if err != nil {
		return false, err
	}
	if strings.Contains(requestor.Response.Body, "User exists") {
		return true, nil
	} else if strings.Contains(requestor.Response.Body, "User does not exist") {
		return false, nil
	}
	return false, nil
}
