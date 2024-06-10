package appconnectorschedule

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

const (
	mgmtConfig       = "/mgmtconfig/v1/admin/customers/"
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
func GetSchedule(service *services.Service) (*AssistantSchedule, *http.Response, error) {
	v := new(AssistantSchedule)
	path := fmt.Sprintf("%v", mgmtConfig+service.Client.Config.CustomerID+scheduleEndpoint)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Configure a App Connector schedule frequency to delete the in active connectors with configured frequency.
func CreateSchedule(service *services.Service, assistantSchedule AssistantSchedule) (*AssistantSchedule, *http.Response, error) {
	v := new(AssistantSchedule)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+scheduleEndpoint, nil, assistantSchedule, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func UpdateSchedule(service *services.Service, schedulerID string, assistantSchedule *AssistantSchedule) (*http.Response, error) {
	// Validate FrequencyInterval
	validIntervals := map[string]bool{"5": true, "7": true, "14": true, "30": true, "60": true, "90": true}
	if _, valid := validIntervals[assistantSchedule.FrequencyInterval]; !valid {
		return nil, fmt.Errorf("invalid FrequencyInterval: %s", assistantSchedule.FrequencyInterval)
	}

	// Check if the schedule is enabled
	if !assistantSchedule.Enabled {
		return nil, fmt.Errorf("cannot update a disabled schedule")
	}

	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+scheduleEndpoint, schedulerID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, assistantSchedule, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
