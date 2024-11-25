package appconnectorschedule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	mgmtConfig       = "/zpa/mgmtconfig/v1/admin/customers/"
	scheduleEndpoint = "/connectorSchedule"
)

type AssistantSchedule struct {
	// The unique identifier for the App Connector auto deletion configuration for a customer. This field is only required for the PUT request to update the frequency of the App Connector Settings.
	ID string `json:"id,omitempty"`

	// The unique identifier of the ZPA tenant.
	CustomerID string `json:"customerId"`

	// Indicates if the App Connectors are included for deletion if they are in a disconnected state based on frequencyInterval and frequency values.
	DeleteDisabled bool `json:"deleteDisabled"`

	// Indicates if the setting for deleting App Connectors is enabled or disabled.
	Enabled bool `json:"enabled"`

	// The scheduled frequency at which the disconnected App Connectors are deleted.
	Frequency string `json:"frequency"`

	// The interval for the configured frequency value. The minimum supported value is 5.
	FrequencyInterval string `json:"frequencyInterval"`
}

// Get a Configured App Connector schedule frequency.
func GetSchedule(ctx context.Context, service *zscaler.Service) (*AssistantSchedule, *http.Response, error) {
	v := new(AssistantSchedule)
	path := fmt.Sprintf("%v", mgmtConfig+service.Client.GetCustomerID()+scheduleEndpoint)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Configure a App Connector schedule frequency to delete the in active connectors with configured frequency.
func CreateSchedule(ctx context.Context, service *zscaler.Service, assistantSchedule AssistantSchedule) (*AssistantSchedule, *http.Response, error) {
	v := new(AssistantSchedule)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+scheduleEndpoint, nil, assistantSchedule, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func UpdateSchedule(ctx context.Context, service *zscaler.Service, schedulerID string, assistantSchedule *AssistantSchedule) (*http.Response, error) {
	// Validate FrequencyInterval
	validIntervals := map[string]bool{"5": true, "7": true, "14": true, "30": true, "60": true, "90": true}
	if _, valid := validIntervals[assistantSchedule.FrequencyInterval]; !valid {
		return nil, fmt.Errorf("invalid FrequencyInterval: %s", assistantSchedule.FrequencyInterval)
	}

	// Check if the schedule is enabled
	if !assistantSchedule.Enabled {
		return nil, fmt.Errorf("cannot update a disabled schedule")
	}

	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+scheduleEndpoint, schedulerID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, nil, assistantSchedule, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
