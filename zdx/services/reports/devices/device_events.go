package devices

const (
	deviceEventsEndpoint = "/events"
)

/*
https://help.zscaler.com/zdx/reports#/devices/{deviceid}/events-get
Gets the Events metrics trend for a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
The event metrics include Zscaler, Hardware, Software and Network event changes.
*/

type DeviceEvents struct {
	TimeStamp int      `json:"timestamp,omitempty"`
	Events    []Events `json:"instances,omitempty"`
}

type Events struct {
	Category    string `json:"category,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Prev        string `json:"prev,omitempty"`
	Curr        string `json:"curr,omitempty"`
}
