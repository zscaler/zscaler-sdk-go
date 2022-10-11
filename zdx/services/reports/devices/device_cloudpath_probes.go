package devices

const (
	deviceCloudPathProbesEndpoint = "/cloudpath-probes"
)

/*
https://help.zscaler.com/zdx/reports#/devices/{deviceid}/apps/{appid}/cloudpath-probes-get
Gets the list of all active Cloud Path probes on a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
/devices/{deviceid}/apps/{appid}/cloudpath-probes
*/
type DeviceCloudPathProbes struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	NumProbes int    `json:"num_probes,omitempty"`
}

type AverageLatency struct {
	LegSRC  string `json:"leg_src,omitempty"`
	LegDst  string `json:"leg_dst,omitempty"`
	Latency int    `json:"latency,omitempty"`
}
