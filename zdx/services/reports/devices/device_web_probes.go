package devices

const (
	deviceWebProbesEndpoint = "/web-probes"
)

/*
https://help.zscaler.com/zdx/reports#/devices/{deviceid}/apps/{appid}/web-probes-get
Gets the list of all active web probes on a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
/devices/{deviceid}/apps/{appid}/web-probes
*/

type DeviceWebProbes struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	NumProbes int    `json:"num_probes,omitempty"`
	AvgScore  int    `json:"avg_score,omitempty"`
	AvgPFT    int    `json:"avg_pft,omitempty"`
}
