package getotp

import (
	"context"
	"fmt"
	"net/url"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	getOtpEndpoint = "/zcc/papi/public/v1/getOtp"
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

func GetOtp(ctx context.Context, service *zscaler.Service, udid string) (*OtpResponse, error) {
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
	_, err := service.Client.NewZccRequestDo(ctx, "GET", fullURL, nil, nil, &otpResponse)
	if err != nil {
		return nil, err
	}
	return &otpResponse, err
}
