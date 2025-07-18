package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

const (
	deviceTopProcessEndpoint = "top-processes"
)

type DeviceTopProcesses struct {
	TimeStamp    int            `json:"timestamp,omitempty"`
	TopProcesses []TopProcesses `json:"top_processes,omitempty"`
}

type TopProcesses struct {
	Category  string      `json:"category,omitempty"`
	Processes []Processes `json:"processes,omitempty"`
}

type Processes struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetDeviceTopProcesses gets the top processes for a device's deep trace session
func GetDeviceTopProcesses(ctx context.Context, service *zscaler.Service, deviceID int, traceID string, filters common.GetFromToFilters) ([]DeviceTopProcesses, *http.Response, error) {
	var v []DeviceTopProcesses
	path := fmt.Sprintf("%v/%v/deeptraces/%v/%v", devicesEndpoint, deviceID, traceID, deviceTopProcessEndpoint)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
