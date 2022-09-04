package getdevices

const (
	getDevicesEndpoint = "getDevices"
)

type GetDevices struct {
	AgentVersion            string `json:"agentVersion"`
	CompanyName             string `json:"companyName"`
	ConfigDownloadTime      string `json:"config_download_time"`
	DeregistrationTimestamp string `json:"deregistrationTimestamp"`
	Detail                  string `json:"detail"`
	DownloadCount           int    `json:"download_count"`
	HardwareFingerprint     string `json:"hardwareFingerprint"`
	KeepAliveTime           string `json:"keepAliveTime"`
	LastSeenTime            string `json:"last_seen_time"`
	MacAddress              string `json:"macAddress"`
	MachineHostname         string `json:"machineHostname"`
	Manufacturer            string `json:"manufacturer"`
	OsVersion               string `json:"osVersion"`
	Owner                   string `json:"owner"`
	PolicyName              string `json:"policyName"`
	RegistrationState       string `json:"registrationState"`
	RegistrationTime        string `json:"registration_time"`
	State                   int    `json:"state"`
	TunnelVersion           string `json:"tunnelVersion"`
	Type                    int    `json:"type"`
	Udid                    string `json:"udid"`
	UpmVersion              string `json:"upmVersion"`
	User                    string `json:"user"`
	VpnState                int    `json:"vpnState"`
	ZappArch                string `json:"zappArch"`
}

func (service *Service) GetAll() ([]GetDevices, error) {
	var devices []GetDevices
	err := service.Client.Read(getDevicesEndpoint, &devices)
	return devices, err
}
