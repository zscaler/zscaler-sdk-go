package applications

import (
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
)

func TestGetAppScores(t *testing.T) {
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

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := GetAllApps(service, filters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetAppScores test.")
		return
	}

	firstAppID := apps[0].ID

	// Call GetAppScores with the first app's ID
	scores, resp, err := GetAppScores(service, firstAppID, filters)
	if err != nil {
		t.Fatalf("Error getting app scores: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(scores) == 0 {
		t.Log("No app scores found.")
	} else {
		t.Logf("Retrieved %d app scores", len(scores))
		for _, score := range scores {
			t.Logf("Metric: %s, Unit: %s", score.Metric, score.Unit)
			for _, dp := range score.DataPoints {
				t.Logf("Timestamp: %d, Value: %f", dp.TimeStamp, dp.Value)
			}
		}
	}
}

func TestGetAppMetrics(t *testing.T) {
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

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := GetAllApps(service, filters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetAppMetrics test.")
		return
	}

	firstAppID := apps[0].ID

	// Call GetAppMetrics with the first app's ID
	metrics, resp, err := GetAppMetrics(service, firstAppID, filters)
	if err != nil {
		t.Fatalf("Error getting app metrics: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(metrics) == 0 {
		t.Log("No app metrics found.")
	} else {
		t.Logf("Retrieved %d app metrics", len(metrics))
		for _, metric := range metrics {
			t.Logf("Metric: %s, Unit: %s", metric.Metric, metric.Unit)
			for _, dp := range metric.DataPoints {
				t.Logf("Timestamp: %d, Value: %f", dp.TimeStamp, dp.Value)
			}
		}
	}
}
