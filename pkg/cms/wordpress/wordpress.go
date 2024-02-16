package wordpress

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/44za12/mailsleuth/internal/requestor"
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	getUrl, _ := url.Parse("https://public-api.wordpress.com/rest/v1.1/users/" + email + "/auth-options?http_envelope=1&locale=fr")
	requestor.Headers["TE"] = "Trailers"
	cookies := []*http.Cookie{
		{Name: "G_ENABLED_IDPS", Value: "google"},
		{Name: "ccpa_applies", Value: "true"},
		{Name: "usprivacy", Value: "1YNN"},
		{Name: "landingpage_currency", Value: "EUR"},
		{Name: "wordpress_test_cookie", Value: "WP Cookie check"},
	}
	requestor.Request.Client.Jar.SetCookies(getUrl, cookies)
	err := requestor.GET(getUrl)
	if err != nil {
		return false, err
	}
	if strings.Contains(requestor.Response.Body, "email_verified") {
		return true, nil
	}
	return false, nil
}
