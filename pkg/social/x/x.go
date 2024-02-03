package x

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	r, err := http.Get(twitterUrl + "?" + data.Encode())

	if err != nil {
		return false, err
	}

	r.Header.Add("User-Agent", utils.RandomUserAgent())

	if err != nil {
		return false, err
	}

	if r.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Status code error: %d %s", r.StatusCode, r.Status)
		return false, errors.New(msg)
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return false, err
	}

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
