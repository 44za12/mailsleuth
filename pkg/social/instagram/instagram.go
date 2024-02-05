package instagram

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/44za12/mailsleuth/internal/utils"
)

func Check(email string, client *http.Client) (bool, error) {
	standardHeaders := utils.StandardHeaders()
	token, err := getCSRFToken(client, standardHeaders)

	if err != nil {
		return false, err
	}

	if strings.EqualFold(token, "") {
		return false, errors.New("CSRF token not found")
	}

	attempUrl := "https://www.instagram.com/accounts/web_create_ajax/attempt/"
	data := url.Values{}
	data.Set("email", email)
	data.Set("username", utils.RandomString(12))
	req, err := http.NewRequest("POST", attempUrl, strings.NewReader(data.Encode()))

	if err != nil {
		return false, err
	}

	for key, value := range standardHeaders {
		req.Header.Set(key, value)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "csrftoken="+token+";")
	req.Header.Add("X-Csrftoken", token)

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Status code error: %d %s", resp.StatusCode, resp.Status)
		return false, errors.New(msg)
	}

	body, err := utils.DecodeResponseBody(resp)

	if err != nil {
		return false, err
	}
	match, err := regexp.MatchString("email_is_taken", string(body))

	if err != nil {
		return false, err
	}

	if !match {
		return false, nil
	}

	return true, nil
}

func getCSRFToken(client *http.Client, standardHeaders map[string]string) (string, error) {
	url := "https://www.instagram.com/accounts/emailsignup/"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}
	for key, value := range standardHeaders {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Status code error: %d %s", resp.StatusCode, resp.Status)
		return "", errors.New(msg)
	}
	body, err := utils.DecodeResponseBody(resp)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`"csrf_token":"([^"]+)"`)
	match := re.FindStringSubmatch(string(body))
	if len(match) == 0 {
		return "", errors.New("CSRF Token not found")
	}

	return match[1], nil
}
