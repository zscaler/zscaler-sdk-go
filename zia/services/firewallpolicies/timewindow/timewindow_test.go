package timewindow

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestTimeWindow_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	tWindows, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting time windows: %v", err)
		return
	}
	if len(tWindows) == 0 {
		t.Errorf("No time windows found")
		return
	}
	tWindowName := tWindows[0].Name
	t.Log("Getting time window by name:" + tWindowName)
	tWindow, err := service.GetTimeWindowByName(tWindowName)
	if err != nil {
		t.Errorf("Error getting time windows by name: %v", err)
		return
	}
	if tWindow.Name != tWindowName {
		t.Errorf("time window name does not match: expected %s, got %s", tWindowName, tWindow.Name)
		return
	}
	// Negative Test: Try to retrieve a time window with a non-existent name
	nonExistentName := "ThisTimeWindowDoesNotExist"
	_, err = service.GetTimeWindowByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	timeWindows, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting time window : %v", err)
		return
	}
	if len(timeWindows) == 0 {
		t.Errorf("No time window  found")
		return
	}

	// Validate time window
	for _, timeWindow := range timeWindows {
		// Checking if essential fields are not empty
		if timeWindow.ID == 0 {
			t.Errorf("time window  ID is empty")
		}
		if timeWindow.Name == "" {
			t.Errorf("time window Name is empty")
		}
	}
}
