package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

const (
	deviceEventsEndpoint = "events"
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
// The event metrics include Zscaler, Hardware, Software, and Network event changes.
func GetEvents(ctx context.Context, service *zscaler.Service, deviceID int, filters common.GetFromToFilters) ([]DeviceEvents, *http.Response, error) {
	var v []DeviceEvents
	path := fmt.Sprintf("%v/%v/%v", devicesEndpoint, deviceID, deviceEventsEndpoint)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
