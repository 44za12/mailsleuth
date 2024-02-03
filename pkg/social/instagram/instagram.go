package instagram

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/44za12/mailsleuth/internal/utils"
)

func Check(email string, client *http.Client) (bool, error) {
	userAgent := utils.RandomUserAgent()
	token, err := getCSRFToken(client, userAgent)

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
	r, err := http.NewRequest("POST", attempUrl, strings.NewReader(data.Encode()))

	if err != nil {
		return false, err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("User-Agent", userAgent)
	r.Header.Add("Cookie", "csrftoken="+token+";")
	r.Header.Add("X-Csrftoken", token)

	res, err := client.Do(r)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Status code error: %d %s", res.StatusCode, res.Status)
		return false, errors.New(msg)
	}

	body, err := io.ReadAll(res.Body)

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

func getCSRFToken(client *http.Client, userAgent string) (string, error) {
	url := "https://www.instagram.com/accounts/emailsignup/"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Status code error: %d %s", res.StatusCode, res.Status)
		return "", errors.New(msg)
	}
	re := regexp.MustCompile(`"csrf_token":"([^"]+)"`)
	match := re.FindStringSubmatch(string(body))

	if len(match) == 0 {
		return "", errors.New("CSRF Token not found")
	}

	return match[1], nil
}
