package spotify

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/44za12/mailsleuth/internal/requestor"
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

func Check(email string, requestor *requestor.Requestor) (bool, error) {
	url, _ := url.Parse(fmt.Sprintf("https://spclient.wg.spotify.com/signup/public/v1/account?validate=1&email=%s", email))
	err := requestor.GET(url)
	if err != nil {
		return false, err
	}
	var r response
	err = json.Unmarshal([]byte(requestor.Response.Body), &r)

	if err != nil {
		return false, err
	}

	if r.Status != 20 {
		return false, nil
	}
	return true, nil
}
