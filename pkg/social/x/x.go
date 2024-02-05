package x

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/44za12/mailsleuth/internal/utils"
)

type response struct {
	Valid bool   `json:"valid"`
	Msg   string `json:"msg"`
	Taken bool   `json:"taken"`
}

func Check(email string, client *http.Client) (bool, error) {
	twitterUrl := "https://api.x.com/i/users/email_available.json"
	data := url.Values{}
	data.Set("email", email)
	standardHeaders := utils.StandardHeaders()
	req, err := http.NewRequest("GET", twitterUrl+"?"+data.Encode(), nil)

	if err != nil {
		return false, err
	}
	for key, value := range standardHeaders {
		req.Header.Set(key, value)
	}
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Status code error: %d %s", resp.StatusCode, resp.Status)
		return false, errors.New(msg)
	}
	body, err := utils.DecodeResponseBody(resp)

	if err != nil {
		return false, err
	}

	utils.SaveResponse(string(body), "twitter.json")
	var response response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	if !response.Taken {
		return false, nil
	}
	return true, nil
}
