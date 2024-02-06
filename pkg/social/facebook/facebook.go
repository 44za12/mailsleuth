package facebook

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/44za12/mailsleuth/internal/requestor"
)

var (
	LOGIN_URL, _           = url.Parse("https://mbasic.facebook.com/login/identify/?ctx=recover&search_attempts=0&alternate_search=0&toggle_search_mode=0")
	FORGOT_PASSWORD_URL, _ = url.Parse("https://mbasic.facebook.com/recover/initiate/?c=%2Flogin%2F&fl=initiate_view&ctx=msite_initiate_view")
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	requestor.AddExtraHeaders = false
	err := requestor.GET(LOGIN_URL)
	if err != nil {
		return false, err
	}
	formData := addParameters(requestor.Response.Body, email)
	if formData == nil {
		return false, fmt.Errorf("error finding parameters")
	}
	requestor.Request.Parameters = formData
	err = requestor.POST(LOGIN_URL)
	if err != nil {
		return false, err
	}
	err = requestor.GET(FORGOT_PASSWORD_URL)
	if err != nil {
		return false, err
	}
	return findDivWithClass(requestor.Response.Body), nil
}

func addParameters(htmlContent, email string) map[string]string {
	formData := make(map[string]string)
	inputRegex := regexp.MustCompile(`<input type="hidden" name="([^"]+)" value="([^"]*)"`)
	matches := inputRegex.FindAllStringSubmatch(htmlContent, -1)
	if matches == nil {
		return nil
	}
	for _, match := range matches {
		switch match[1] {
		case "lsd":
			{
				formData["lsd"] = match[2]
			}
		case "jazoest":
			{
				formData["jazoest"] = match[2]
			}
		}
	}
	formData["email"] = email
	formData["did_submit"] = "Search"
	return formData
}

func findDivWithClass(htmlContent string) bool {
	pattern := `<div\s+(?:class=(?:"bb\s+bc"|'bb\s+bc'|'bb\s+bc'|"bb\s+bc"))[^>]*>`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(htmlContent)
}
