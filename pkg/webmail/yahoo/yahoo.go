package yahoo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/44za12/mailsleuth/internal/requestor"
)

var (
	GET_URL, _  = url.Parse("https://login.yahoo.com/account/create?specId=yidregsimplified&lang=en-US&src=&done=https%3A%2F%2Fwww.yahoo.com&display=login")
	POST_URL, _ = url.Parse("https://login.yahoo.com/account/module/create?validateField=userId")
)

type yahooErrorResp struct {
	Errors []errItem `json:"errors"`
}

type errItem struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	var resp yahooErrorResp
	err := requestor.GET(GET_URL)
	if err != nil {
		return false, err
	}
	emailParts := strings.Split(email, "@")
	userName, domain := emailParts[0], emailParts[1]
	aCrumb := getACrumb(requestor.Response.Cookies)
	sessionIndex := getSessionIndex(requestor.Response.Body)
	data := `{"acrumb":"` + aCrumb + `","specId":"yidregsimplified","userId":"` + userName + `","sessionIndex":"` + sessionIndex + `","yidDomain":"` + domain + `"}`
	requestor.Request.RawParameters = data
	requestor.Headers["X-Requested-With"] = "XMLHttpRequest"
	requestor.Request.Cookies = requestor.Response.Cookies
	err = requestor.POST(POST_URL)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(requestor.Response.Body), &resp)
	if err != nil {
		return false, err
	}
	exists := checkUsernameExists(resp)
	return exists, nil
}

func getACrumb(cookies []*http.Cookie) string {
	for _, c := range cookies {
		re := regexp.MustCompile(`s=(?P<acrumb>[^;^&]*)`)
		match := re.FindStringSubmatch(c.Value)
		if len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func getSessionIndex(resp string) string {
	re := regexp.MustCompile(`value="([^"]+)" name="sessionIndex"`)
	match := re.FindStringSubmatch(resp)
	if len(match) > 1 {
		return string(match[1])
	}
	return ""
}

func checkUsernameExists(resp yahooErrorResp) bool {
	for _, item := range resp.Errors {
		if item.Name == "userId" && item.Error == "IDENTIFIER_EXISTS" {
			return true
		}
	}
	return false
}
