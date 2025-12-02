package timewindow

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestTimeWindow(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "timewindow", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Test GetAll and GetByName
	t.Run("GetAll and GetByName", func(t *testing.T) {
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
	})

	// Test response format validation
	t.Run("ResponseFormatValidation", func(t *testing.T) {
		timeWindows, err := GetAll(context.Background(), service)
		if err != nil {
			t.Errorf("Error getting time window: %v", err)
			return
		}
		if len(timeWindows) == 0 {
			t.Errorf("No time window found")
			return
		}

		for _, timeWindow := range timeWindows {
			if timeWindow.ID == 0 {
				t.Errorf("time window ID is empty")
			}
			if timeWindow.Name == "" {
				t.Errorf("time window Name is empty")
			}
		}
	})

	// Test case sensitivity
	t.Run("CaseSensitivityOfGetByName", func(t *testing.T) {
		knownName := "Weekends"
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
			if timeWindows.Name != knownName {
				t.Errorf("Expected role name to be '%s' for variation '%s', but got '%s'", knownName, variation, timeWindows.Name)
			}
		}
	})
}
