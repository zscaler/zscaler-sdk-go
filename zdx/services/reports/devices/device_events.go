package devices

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
)

const (
	deviceEventsEndpoint = "v1/events"
)

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

// Gets the Events metrics trend for a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
// The event metrics include Zscaler, Hardware, Software and Network event changes.
func (service *Service) GetEvents(deviceID string, filters common.GetFromToFilters) (*DeviceEvents, *http.Response, error) {
	v := new(DeviceEvents)
	path := fmt.Sprintf("%v/%v/%v", devicesEndpoint, deviceID, deviceEventsEndpoint)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
