package alerts

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
)

func TestGetOngoingAlerts(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	alerts, resp, err := GetOngoingAlerts(service)
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

	alerts, resp, err := GetHistoricalAlerts(service)
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

/*
	func TestGetAlert(t *testing.T) {
		client, err := tests.NewZdxClient()
		if err != nil {
			t.Fatalf("Error creating client: %v", err)
		}

		service := services.New(client)

		//Get the first alert from historical alerts
		alerts, resp, err := GetHistoricalAlerts(service)
		if err != nil {
			t.Fatalf("Error getting historical alerts: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if len(alerts.Alerts) == 0 {
			t.Log("No historical alerts found.")
			return
		}

		firstAlertID := strconv.Itoa(alerts.Alerts[0].ID)

		// Get the specific alert by ID
		alert, resp, err := GetAlert(service, firstAlertID)
		if err != nil {
			t.Fatalf("Error getting alert: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if alert.ID == 0 {
			t.Log("No alert found with the specified ID.")
		} else {
			t.Logf("Retrieved alert ID: %d, Rule Name: %s", alert.ID, alert.RuleName)
		}
	}
*/
func TestGetAffectedDevices(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Get the first alert from ongoing alerts
	// alerts, resp, err := GetOngoingAlerts(service)
	// if err != nil {
	// 	t.Fatalf("Error getting ongoing alerts: %v", err)
	// }

	// if resp.StatusCode != http.StatusOK {
	// 	t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	// }

	// if len(alerts.Alerts) == 0 {
	// 	t.Log("No ongoing alerts found.")
	// 	return
	// }

	// firstAlertID := strconv.Itoa(alerts.Alerts[0].ID)

	// Get the affected devices for the specific alert by ID
	affectedDevices, resp, err := GetAffectedDevices(service, "7379740076528711991")
	if err != nil {
		t.Fatalf("Error getting affected devices: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(affectedDevices.Alerts) == 0 {
		t.Log("No affected devices found for the specified alert.")
	} else {
		t.Logf("Retrieved %d affected devices", len(affectedDevices.Alerts))
		for _, alert := range affectedDevices.Alerts {
			t.Logf("Alert ID: %d, Rule Name: %s", alert.ID, alert.RuleName)
		}
	}
}
