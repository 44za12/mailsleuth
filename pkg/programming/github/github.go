package github

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/44za12/mailsleuth/internal/requestor"
)

var (
	GET_URL, _  = url.Parse("https://github.com/join")
	POST_URL, _ = url.Parse("https://github.com/signup_check/email")
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	err := requestor.GET(GET_URL)
	if err != nil {
		return false, err
	}
	tokenRegex := regexp.MustCompile(`<auto-check src="/signup_check/email[\s\S]*?value="([\S]+)"`)
	tokens := tokenRegex.FindStringSubmatch(requestor.Response.Body)
	if len(tokens) < 2 {
		return false, fmt.Errorf("failed to extract authenticity tokens")
	}
	authenticityToken := tokens[1]
	data := make(map[string]string)
	data["value"] = email
	data["authenticity_token"] = authenticityToken
	requestor.Request.Parameters = data
	requestor.Headers["Accept"] = "*/*"
	time.Sleep(time.Duration(time.Duration.Seconds(1)))
	err = requestor.POST(POST_URL)
	if err != nil {
		return false, err
	}
	if requestor.Response.StatusCode == http.StatusUnprocessableEntity {
		return true, nil
	} else if requestor.Response.StatusCode == http.StatusOK {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status code: %d", requestor.Response.StatusCode)
}
