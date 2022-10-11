package devices

const (
	deviceCloudPathEndpoint = "/cloudpath"
)

/*
https://help.zscaler.com/zdx/reports#/devices/{deviceid}/apps/{appid}/cloudpath-probes/{probeid}/cloudpath-get
Gets the Cloud Path hop data for an application on a specific device.
Includes the summary data for the entire path like the total number of hops, packet loss, latency, and tunnel type (if available).
It also includes a similar summary of data for each individual hop. If the time range is not specified, the endpoint defaults to the last 2 hours.
*/

type DeviceCloudPath struct {
	TimeStamp int         `json:"timestamp,omitempty"`
	CloudPath []CloudPath `json:"cloudpath,omitempty"`
}

type CloudPath struct {
	SRC           string `json:"src,omitempty"`
	DST           string `json:"dst,omitempty"`
	NumHops       int    `json:"num_hops,omitempty"`
	Latency       int    `json:"latency,omitempty"`
	Loss          int    `json:"loss,omitempty"`
	NumUnrespHops int    `json:"num_unresp_hops,omitempty"`
	TunnelType    int    `json:"tunnel_type,omitempty"`
	Hops          []Hops `json:"hops,omitempty"`
}

type Hops struct {
	IP          string `json:"ip,omitempty"`
	GWMac       string `json:"gw_mac,omitempty"`
	GWMacVendor string `json:"gw_mac_vendor,omitempty"`
	PkgSent     int    `json:"pkt_sent,omitempty"`
	PkgRcvd     int    `json:"pkt_rcvd,omitempty"`
	LatencyMin  int    `json:"latency_min,omitempty"`
	LatencyMax  int    `json:"latency_max,omitempty"`
	LatencyAvg  int    `json:"latency_avg,omitempty"`
	LatencyDiff int    `json:"latency_diff,omitempty"`
}
