package getotp

import (
	"fmt"
	"net/url"
)

const (
	getOtpEndpoint = "/public/v1/getOtp"
)

type OtpResponse struct {
	AntiTemperingDisableOtp string `json:"antiTemperingDisableOtp"`
	DeceptionSettingsOtp    string `json:"deceptionSettingsOtp"`
	ExitOtp                 string `json:"exitOtp"`
	LogoutOtp               string `json:"logoutOtp"`
	Otp                     string `json:"otp"`
	RevertOtp               string `json:"revertOtp"`
	UninstallOtp            string `json:"uninstallOtp"`
	ZdpDisableOtp           string `json:"zdpDisableOtp"`
	ZdxDisableOtp           string `json:"zdxDisableOtp"`
	ZiaDisableOtp           string `json:"ziaDisableOtp"`
	ZpaDisableOtp           string `json:"zpaDisableOtp"`
}

type GetOtpQuery struct {
	Udid string `json:"udid,omitempty" url:"udid,omitempty"`
}

func (service *Service) GetOtp(udid string) (*OtpResponse, error) {
	queryParams := url.Values{}
	if udid != "" {
		queryParams.Set("udid", udid)
	}

	// Build the full URL with query parameters
	fullURL := getOtpEndpoint
	if len(queryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s", getOtpEndpoint, queryParams.Encode())
	}

	var otpResponse OtpResponse
	_, err := service.Client.NewRequestDo("GET", fullURL, nil, nil, &otpResponse)
	if err != nil {
		return nil, err
	}
	return &otpResponse, err
}
