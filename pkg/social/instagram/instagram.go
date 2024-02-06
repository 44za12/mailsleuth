package instagram

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/44za12/mailsleuth/internal/requestor"
	"github.com/44za12/mailsleuth/internal/utils"
)

var (
	CSRF_URL, _      = url.Parse("https://www.instagram.com/accounts/emailsignup/")
	REGISTGER_URL, _ = url.Parse("https://www.instagram.com/accounts/web_create_ajax/attempt/")
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	err := requestor.GET(CSRF_URL)
	if err != nil {
		return false, err
	}
	csrfRe := regexp.MustCompile(`"csrf_token":"([^"]+)"`)
	match := csrfRe.FindStringSubmatch(string(requestor.Response.Body))
	if len(match) == 0 {
		return false, errors.New("CSRF Token not found")
	}
	token := match[1]
	if strings.EqualFold(token, "") {
		return false, errors.New("CSRF token not found")
	}
	data := make(map[string]string)
	data["email"] = email
	data["username"] = utils.RandomString(12)
	requestor.Request.Parameters = data
	requestor.Headers["Cookie"] = "csrftoken=" + token + ";"
	requestor.Headers["X-Csrftoken"] = token
	time.Sleep(time.Duration(time.Duration.Seconds(1)))
	err = requestor.POST(REGISTGER_URL)
	if err != nil {
		return false, err
	}
	if requestor.Response.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Status code error: %d", requestor.Response.StatusCode)
		return false, errors.New(msg)
	}
	isTaken, err := regexp.MatchString("email_is_taken", requestor.Response.Body)

	if err != nil {
		return false, err
	}

	if !isTaken {
		return false, nil
	}

	return true, nil
}
