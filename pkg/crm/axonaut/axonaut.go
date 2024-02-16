package axonaut

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/44za12/mailsleuth/internal/requestor"
	"github.com/44za12/mailsleuth/internal/utils"
)

var (
	CSRF_URL, _ = url.Parse("https://axonaut.com/onboarding/")
	POST_URL, _ = url.Parse("https://axonaut.com/onboarding/subscription")
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	err := requestor.GET(CSRF_URL)
	if err != nil {
		return false, err
	}
	csrfRegex := regexp.MustCompile(`<input type="hidden" name="_csrf" value="([^"]*)"`)
	match := csrfRegex.FindStringSubmatch(string(requestor.Response.Body))
	if len(match) == 0 {
		return false, errors.New("axonaut: csrf token not found")
	}
	token := match[1]
	cookieRegex := regexp.MustCompile(`setCookie\('landing', ([^"]*)\);`)
	match = cookieRegex.FindStringSubmatch(string(requestor.Response.Body))
	if len(match) == 0 {
		return false, errors.New("axonaut: cookie not found")
	}
	cookie := match[1]
	formData := make(map[string]string)
	formData["_csrf"] = token
	formData["cookie"] = cookie
	formData["emailAddress"] = email
	formData["password"] = utils.RandomString(8)
	requestor.Request.Parameters = formData
	err = requestor.POST(POST_URL)
	if err != nil {
		return false, err
	}
	switch requestor.Response.Body {
	case "1":
		return false, nil
	case "0":
		return true, nil
	}
	return false, errors.New("axonaut: Some error in the formdata")
}
