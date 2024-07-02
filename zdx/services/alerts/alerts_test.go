package alerts

import (
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
)

func TestGetOngoingAlerts(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	// Call GetOngoingAlerts with the filters
	alerts, resp, err := GetOngoingAlerts(service, filters)
	if err != nil {
		t.Fatalf("Error getting ongoing alerts: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(alerts.Alerts) == 0 {
		t.Log("No ongoing alerts found.")
	} else {
		t.Logf("Retrieved %d ongoing alerts", len(alerts.Alerts))
		for _, alert := range alerts.Alerts {
			t.Logf("Alert ID: %d, Rule Name: %s", alert.ID, alert.RuleName)
		}
	}
}

func TestGetHistoricalAlerts(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Define a time filter for the last 14 days
	now := time.Now()
	from := now.Add(-10 * 24 * time.Hour).Unix() // 14 days ago
	to := now.Unix()                             // current time
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	alerts, resp, err := GetHistoricalAlerts(service, filters)
	if err != nil {
		t.Fatalf("Error getting historical alerts: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(alerts.Alerts) == 0 {
		t.Log("No historical alerts found.")
	} else {
		t.Logf("Retrieved %d historical alerts", len(alerts.Alerts))
		for _, alert := range alerts.Alerts {
			t.Logf("Alert ID: %d, Rule Name: %s", alert.ID, alert.RuleName)
		}
	}
}

// func TestGetAlert(t *testing.T) {
// 	client, err := tests.NewZdxClient()
// 	if err != nil {
// 		t.Fatalf("Error creating client: %v", err)
// 	}

// 	service := services.New(client)

// 	// //Get the first alert from historical alerts
// 	// alerts, resp, err := GetHistoricalAlerts(service)
// 	// if err != nil {
// 	// 	t.Fatalf("Error getting historical alerts: %v", err)
// 	// }

// 	// if resp.StatusCode != http.StatusOK {
// 	// 	t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
// 	// }

// 	// if len(alerts.Alerts) == 0 {
// 	// 	t.Log("No historical alerts found.")
// 	// 	return
// 	// }

// 	// firstAlertID := strconv.Itoa(alerts.Alerts[0].ID)

// 	// Get the specific alert by ID
// 	alert, resp, err := GetAlert(service, "7381380182807289758")
// 	if err != nil {
// 		t.Fatalf("Error getting alert: %v", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
// 	}

// 	if alert.ID == 0 {
// 		t.Log("No alert found with the specified ID.")
// 	} else {
// 		t.Logf("Retrieved alert ID: %d, Rule Name: %s", alert.ID, alert.RuleName)
// 	}
// }

func TestGetAffectedDevices(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Define a time filter for the last 10 days
	now := time.Now()
	from := now.Add(-10 * 24 * time.Hour).Unix() // 10 days ago
	to := now.Unix()                             // current time
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	// Get the first alert from historical alerts
	// alerts, resp, err := GetHistoricalAlerts(service, filters)
	// if err != nil {
	// 	t.Fatalf("Error getting historical alerts: %v", err)
	// }

	// if resp.StatusCode != http.StatusOK {
	// 	t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	// }

	// if len(alerts.Alerts) == 0 {
	// 	t.Log("No historical alerts found.")
	// 	return
	// }

	// // Log the retrieved alert IDs for debugging purposes
	// for _, alert := range alerts.Alerts {
	// 	t.Logf("Retrieved Alert ID: %d", alert.ID)
	// }

	// firstAlertID := strconv.Itoa(alerts.Alerts[0].ID)
	// t.Logf("Using First Alert ID: %s", firstAlertID)

	// Get the affected devices for the specific alert by ID using the same filters
	affectedDevices, resp, err := GetAffectedDevices(service, "7381380182807289758", filters)
	if err != nil {
		t.Fatalf("Error getting affected devices: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(affectedDevices.Devices) == 0 {
		t.Log("No affected devices found for the specified alert.")
	} else {
		t.Logf("Retrieved %d affected devices", len(affectedDevices.Devices))
		for _, device := range affectedDevices.Devices {
			t.Logf("Device ID: %d, Name: %s, User ID: %d, User Name: %s, User Email: %s", device.ID, device.Name, device.UserID, device.UserName, device.UserEmail)
		}
	}
}
