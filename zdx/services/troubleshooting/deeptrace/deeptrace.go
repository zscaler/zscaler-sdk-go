package deeptrace

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
)

const (
	deepTracesEndpoint = "v1/devices"
)

type DeepTraceSession struct {
	TraceID      string       `json:"trace_id"`
	TraceDetails TraceDetails `json:"trace_details,omitempty"`
	Status       string       `json:"status,omitempty"`
	CreatedAt    int          `json:"created_at,omitempty"`
	StartedAt    int          `json:"started_at,omitempty"`
	EndedAt      int          `json:"ended_at,omitempty"`
}

type TraceDetails struct {
	SessionName        string `json:"session_name"`
	UserID             string `json:"user_id,omitempty"`
	Username           string `json:"username,omitempty"`
	DeviceID           string `json:"device_id,omitempty"`
	DeviceName         string `json:"device_name,omitempty"`
	WebProbeID         string `json:"web_probe_id,omitempty"`
	WebProbeName       string `json:"web_probe_name,omitempty"`
	CloudPathProbeID   string `json:"cloudpath_probe_id,omitempty"`
	CloudPathProbeName string `json:"cloud_path_name,omitempty"`
	SessionLength      string `json:"session_length,omitempty"`
	ProbeDevice        bool   `json:"probe_device,omitempty"`
}

type DeepTraceSessionPayload struct {
	SessionName          string `json:"session_name"`
	AppID                int    `json:"app_id"`
	WebProbeID           int    `json:"web_probe_id"`
	CloudPathProbeID     int    `json:"cloud_path_probe_id"`
	SessionLengthMinutes int    `json:"session_length_minutes"`
	ProbeDevice          bool   `json:"probe_device"`
}

func GetDeepTraces(service *services.Service, deviceID int) ([]DeepTraceSession, *http.Response, error) {
	var response []DeepTraceSession
	path := fmt.Sprintf("%s/%d/deeptraces", deepTracesEndpoint, deviceID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return response, resp, nil
}

func GetDeepTraceSession(service *services.Service, deviceID int, traceID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%d/deeptraces/%s", deepTracesEndpoint, deviceID, traceID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CreateDeepTraceSession(service *services.Service, deviceID int, payload DeepTraceSessionPayload) (*DeepTraceSession, *http.Response, error) {
	var response DeepTraceSession
	path := fmt.Sprintf("%s/%d/deeptraces", deepTracesEndpoint, deviceID)
	resp, err := service.Client.NewRequestDo("POST", path, nil, payload, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}

func DeleteDeepTraceSession(service *services.Service, deviceID int, traceID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%d/deeptraces/%s", deepTracesEndpoint, deviceID, traceID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
