package getpasswords

import (
	"fmt"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v2/zcc/services"
)

const (
	getPasswordsEndpoint = "/public/v1/getPasswords"
)

type Passwords struct {
	ExitPass             string `json:"exitPass"`
	LogoutPass           string `json:"logoutPass"`
	UninstallPass        string `json:"uninstallPass"`
	ZdSettingsAccessPass string `json:"zdSettingsAccessPass"`
	ZdxDisablePass       string `json:"zdxDisablePass"`
	ZiaDisablePass       string `json:"ziaDisablePass"`
	ZpaDisablePass       string `json:"zpaDisablePass"`
}

type GetPasswordsQueryParams struct {
	Username string `url:"username"`
	OsType   string `url:"osType"`
}

func GetPasswords(service *services.Service, username, osType string) (*Passwords, error) {
	queryParams := url.Values{}
	if username != "" {
		queryParams.Set("username", username)
	}
	if osType != "" {
		queryParams.Set("osType", osType)
	}

	// Build the full URL with query parameters
	fullURL := getPasswordsEndpoint
	if len(queryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s", getPasswordsEndpoint, queryParams.Encode())
	}

	var passwords Passwords
	_, err := service.Client.NewRequestDo("GET", fullURL, nil, nil, &passwords)
	if err != nil {
		return nil, err
	}
	return &passwords, err
}
