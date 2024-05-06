package get_password

const (
	getPasswordsEndpoint = "/public/v1/getPasswords"
)

type GetPassword struct {
	ExitPass             string `json:"exitPass"`
	LogoutPass           string `json:"logoutPass"`
	UninstallPass        string `json:"uninstallPass"`
	ZdSettingsAccessPass string `json:"zdSettingsAccessPass"`
	ZdxDisablePass       string `json:"zdxDisablePass"`
	ZiaDisablePass       string `json:"ziaDisablePass"`
	ZpaDisablePass       string `json:"zpaDisablePass"`
}

type GetPasswordsQuery struct {
	OsType   int    `json:"osType,omitempty" url:"osType,omitempty"`
	Username string `json:"username,omitempty" url:"username,omitempty"`
}

func (service *Service) GetPasswords(osType int, username string) (*GetPassword, error) {
	var passwords GetPassword
	_, err := service.Client.NewRequestDo("GET", getPasswordsEndpoint, GetPasswordsQuery{OsType: osType, Username: username}, nil, &passwords)
	if err != nil {
		return nil, err
	}
	return &passwords, err
}
