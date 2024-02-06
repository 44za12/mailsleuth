package amazon

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/44za12/mailsleuth/internal/requestor"
)

var (
	GET_URL, _  = url.Parse("https://www.amazon.com/ap/signin?openid.pape.max_auth_age=0&openid.return_to=https%3A%2F%2Fwww.amazon.com%2F%3F_encoding%3DUTF8%26ref_%3Dnav_ya_signin&openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.assoc_handle=usflex&openid.mode=checkid_setup&openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0&")
	POST_URL, _ = url.Parse("https://www.amazon.com/ap/signin")
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	err := requestor.GET(GET_URL)
	if err != nil {
		return false, err
	}
	formData, err := extractFormData(requestor.Response.Body)
	if err != nil {
		return false, err
	}
	formData["email"] = email
	requestor.Request.Parameters = formData
	err = requestor.POST(POST_URL)
	if err != nil {
		return false, err
	}
	return checkDivExists(requestor.Response.Body), nil
}

func extractFormData(htmlContent string) (map[string]string, error) {
	formData := make(map[string]string)
	formRegex := regexp.MustCompile(`(?s)<form[^>]*name="signIn"[^>]*>(.*?)</form>`)
	formMatch := formRegex.FindStringSubmatch(htmlContent)
	if formMatch == nil {
		return nil, fmt.Errorf("form not found")
	}
	inputRegex := regexp.MustCompile(`<input type="hidden" name="([^"]+)" value="([^"]*)"`)
	matches := inputRegex.FindAllStringSubmatch(formMatch[1], -1)
	for _, match := range matches {
		formData[match[1]] = match[2]
	}
	return formData, nil
}

func checkDivExists(htmlContent string) bool {
	divRegex := regexp.MustCompile(`<div[^>]*id="auth-password-missing-alert"[^>]*>`)
	return divRegex.MatchString(htmlContent)
}
