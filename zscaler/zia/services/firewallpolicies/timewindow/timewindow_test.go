package timewindow

import (
	"context"
	"strings"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestTimeWindow_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	tWindows, err := GetAll(context.Background(), service)
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
	tWindow, err := GetTimeWindowByName(context.Background(), service, tWindowName)
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
	_, err = GetTimeWindowByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	timeWindows, err := GetAll(context.Background(), service)
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

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Assuming a service with the name "Weekends" exists
	knownName := "Weekends"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve service with name variation: %s", variation)
		timeWindows, err := GetTimeWindowByName(context.Background(), service, variation)
		if err != nil {
			t.Errorf("Error getting service with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if timeWindows.Name != knownName {
			t.Errorf("Expected role name to be '%s' for variation '%s', but got '%s'", knownName, variation, timeWindows.Name)
		}
	}
}
