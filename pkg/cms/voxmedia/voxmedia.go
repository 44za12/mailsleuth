package voxmedia

import (
	"net/url"
	"strings"

	"github.com/44za12/mailsleuth/internal/requestor"
)

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	getUrl, _ := url.Parse("https://auth.voxmedia.com/chorus_auth/email_valid.json")
	requestor.Headers["TE"] = "Trailers"
	requestor.Request.Parameters = map[string]string{"email": email}
	err := requestor.POST(getUrl)
	if err != nil {
		return false, err
	}
	if strings.Contains(requestor.Response.Body, `"available":true`) {
		return false, nil
	} else if strings.Contains(requestor.Response.Body, "You cannot use this email address") {
		return true, nil
	}
	return false, nil
}
