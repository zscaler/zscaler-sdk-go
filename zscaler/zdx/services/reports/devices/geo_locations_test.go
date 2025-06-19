package devices

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

func TestGetGeoLocations(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := GeoLocationFilter{
		GetFromToFilters: common.GetFromToFilters{
			From: int(from),
			To:   int(to),
		},
	}

	// Call GetGeoLocations with the filters
	geoLocations, resp, err := GetGeoLocations(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting geo locations: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(geoLocations) == 0 {
		t.Log("No geolocations found.")
	} else {
		for _, geoLocation := range geoLocations {
			t.Logf("Retrieved geolocation: ID: %s, Name: %s, GeoType: %s, Description: %s", geoLocation.ID, geoLocation.Name, geoLocation.GeoType, geoLocation.Description)
			for _, child := range geoLocation.Children {
				t.Logf("Child GeoLocation - ID: %s, GeoType: %s, Description: %s", child.ID, child.GeoType, child.Description)
			}
		}
	}
}
