package outlook

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/44za12/mailsleuth/internal/requestor"
	"github.com/44za12/mailsleuth/internal/utils"
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	requestor.AddExtraHeaders = false
	URL, _ := url.Parse("https://login.microsoft.com/common/oauth2/token")
	requestor.Request.Parameters = map[string]string{
		"grant_type": "password",
		"password":   utils.RandomString(12),
		"client_id":  "4345a7b9-9a63-4910-a426-35363201d503",
		"username":   email,
		"resource":   "https://graph.windows.net",
		"scope":      "openid",
	}
	err := requestor.POST(URL)
	if err != nil {
		return false, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(requestor.Response.Body), &data); err != nil {
		panic(err)
	}
	jsonErrCode := data["error_codes"]
	x := fmt.Sprintf("%v", jsonErrCode)
	if strings.Contains(x, "50059") {
		//Domain not found in o365 directory. Exiting..."
		return false, nil
	} else if strings.Contains(x, "50034") {
		//User not found
		return false, nil
	} else if strings.Contains(x, "50126") {
		//Valid user, but invalid password
		return true, nil
	} else if strings.Contains(x, "50055") {
		//Valid user, expired password
		return true, nil
	} else if strings.Contains(x, "50056") {
		//User exists, but unable to determine if the password is correct
		return true, nil
	} else if strings.Contains(x, "50053") {
		//Account locked out
		return false, nil
	} else if strings.Contains(x, "50057") {
		//Account disabled
		return false, nil
	} else if strings.Contains(x, "50076") || strings.Contains(x, "50079") {
		//Possible valid login, MFA required.
		return true, nil
	} else if strings.Contains(x, "53004") {
		//Possible valid login, user must enroll in MFA
		return true, nil
	} else if strings.Contains(x, "") {
		//Possible valid login!
		return true, nil
	} else {
		return false, nil
		//Unknown response, run with -debug flag for more information.
	}
}
