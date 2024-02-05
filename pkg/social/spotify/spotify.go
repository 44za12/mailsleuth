package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/44za12/mailsleuth/internal/utils"
)

type Spotify struct {
	Exists bool
}

type response struct {
	Status int `json:"status"`
	Errors struct {
		Email string `json:"email"`
	} `json:"errors"`
	Country                         string `json:"country"`
	CanAcceptLicensesInOneStep      bool   `json:"can_accept_licenses_in_one_step"`
	RequiresMarketingOptIn          bool   `json:"requires_marketing_opt_in"`
	RequiresMarketingOptInText      bool   `json:"requires_marketing_opt_in_text"`
	MinimumAge                      int    `json:"minimum_age"`
	CountryGroup                    string `json:"country_group"`
	SpecificLicenses                bool   `json:"specific_licenses"`
	TermsConditionsAcceptance       string `json:"terms_conditions_acceptance"`
	PrivacyPolicyAcceptance         string `json:"privacy_policy_acceptance"`
	SpotifyMarketingMessagesOption  string `json:"spotify_marketing_messages_option"`
	PretickEula                     bool   `json:"pretick_eula"`
	ShowCollectPersonalInfo         bool   `json:"show_collect_personal_info"`
	UseAllGenders                   bool   `json:"use_all_genders"`
	UseOtherGender                  bool   `json:"use_other_gender"`
	UsePreferNotToSayGender         bool   `json:"use_prefer_not_to_say_gender"`
	ShowNonRequiredFieldsAsOptional bool   `json:"show_non_required_fields_as_optional"`
	DateEndianness                  int    `json:"date_endianness"`
	IsCountryLaunched               bool   `json:"is_country_launched"`
	AllowedCallingCodes             []struct {
		CountryCode string `json:"country_code"`
		CallingCode string `json:"calling_code"`
	} `json:"allowed_calling_codes"`
	PushNotifications bool `json:"push-notifications"`
}

func Check(email string, client *http.Client) (bool, error) {
	url := fmt.Sprintf("https://spclient.wg.spotify.com/signup/public/v1/account?validate=1&email=%s", email)
	req, err := http.NewRequest("GET", url, nil)
	standardHeaders := utils.StandardHeaders()
	if err != nil {
		return false, err
	}
	for key, value := range standardHeaders {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := utils.DecodeResponseBody(resp)

	if err != nil {
		return false, err
	}
	utils.SaveResponse(string(body), "spotify.json")
	var r response
	err = json.Unmarshal(body, &r)

	if err != nil {
		return false, err
	}

	if r.Status != 20 {
		return false, nil
	}
	return true, nil
}
