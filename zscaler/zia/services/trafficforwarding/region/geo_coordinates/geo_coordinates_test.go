package geo_coordinates

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
)

func TestGeoCoordinates(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.243.0/24")
	comment := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Create static IP for testing
	staticIP, _, err := staticips.Create(context.Background(), service, &staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	if err != nil {
		t.Fatalf("Error creating static IP for testing: %v", err)
	}

	// Clean up: delete the static IP after test
	defer func() {
		_, err := staticips.Delete(context.Background(), service, staticIP.ID)
		if err != nil {
			t.Errorf("Error deleting static IP: %v", err)
		}
	}()

	// Retrieve the GeoCoordinates using the latitude and longitude from the staticIP
	coordinate, err := GetByGeoCoordinates(context.Background(), service,
		float64(staticIP.Latitude),
		float64(staticIP.Longitude),
	)
	if err != nil {
		t.Errorf("Error retrieving GeoCoordinates: %v", err)
		return
	}

	if coordinate == nil {
		t.Errorf("No region found for coordinates: Latitude %v, Longitude %v", staticIP.Latitude, staticIP.Longitude)
		return
	}
}
