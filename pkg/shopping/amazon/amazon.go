package amazon

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/44za12/mailsleuth/internal/utils"
)

func Check(email string, client *http.Client) (bool, error) {
	userAgent := utils.RandomUserAgent()
	req, err := http.NewRequest("GET", "https://www.amazon.com/ap/signin?openid.pape.max_auth_age=0&openid.return_to=https%3A%2F%2Fwww.amazon.com%2F%3F_encoding%3DUTF8%26ref_%3Dnav_ya_signin&openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.assoc_handle=usflex&openid.mode=checkid_setup&openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0&", nil)
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	cookies := resp.Cookies()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	htmlContent := string(body)
	formData, err := extractFormData(htmlContent)
	if err != nil {
		return false, err
	}
	formData["email"] = email
	formValues := url.Values{}
	for key, value := range formData {
		formValues.Set(key, value)
	}

	req, err = http.NewRequest("POST", "https://www.amazon.com/ap/signin", bytes.NewBufferString(formValues.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err = client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	return checkDivExists(string(responseBody)), nil
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
