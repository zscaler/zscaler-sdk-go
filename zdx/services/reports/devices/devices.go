package devices

const (
	devicesEndpoint = "/devices"
)

/*
https://help.zscaler.com/zdx/reports#/devices-get
Gets the list of all active devices and its basic details.
The JSON must contain the user’s ID and email address to associate the device to the user.
If the time range is not specified, the endpoint defaults to the last 2 hours.
*/
type DeviceDetail struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type Hardware struct {
	HWModel     string `json:"hw_model,omitempty"`
	HWMFG       string `json:"hw_mfg,omitempty"`
	HWType      string `json:"hw_type,omitempty"`
	HWSerial    string `json:"hw_serial,omitempty"`
	TotMem      string `json:"tot_mem,omitempty"`
	GPU         string `json:"gpu,omitempty"`
	DiskSize    string `json:"disk_size,omitempty"`
	DiskModel   string `json:"disk_model,omitempty"`
	DiskType    string `json:"disk_type,omitempty"`
	CPUMFG      string `json:"cpu_mfg,omitempty"`
	CPUModel    string `json:"cpu_model,omitempty"`
	SpeedGHZ    int    `json:"speed_ghz,omitempty"`
	LogicalProc int    `json:"logical_proc,omitempty"`
	NumCores    int    `json:"num_cores,omitempty"`
}

type Network struct {
	NetType     string `json:"net_type,omitempty"`
	Status      string `json:"status,omitempty"`
	IPv4        string `json:"ipv4,omitempty"`
	IPv6        string `json:"ipv6,omitempty"`
	DNSSRVS     string `json:"dns_srvs,omitempty"`
	DNSSuffix   string `json:"dns_suffix,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
	MAC         string `json:"mac,omitempty"`
	GUID        string `json:"guid,omitempty"`
	WiFiAdapter string `json:"wifi_adapter,omitempty"`
	WiFiType    string `json:"wifi_type,omitempty"`
	SSID        string `json:"ssid,omitempty"`
	Channel     string `json:"channel,omitempty"`
	BSSID       string `json:"bssid,omitempty"`
}

type Software struct {
	OSName        string `json:"os_name,omitempty"`
	OSVer         string `json:"os_ver,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	NetBios       string `json:"netbios,omitempty"`
	User          string `json:"user,omitempty"`
	ClientConnVer string `json:"client_conn_ver,omitempty"`
	ZDXVer        string `json:"zdx_ver,omitempty"`
}
