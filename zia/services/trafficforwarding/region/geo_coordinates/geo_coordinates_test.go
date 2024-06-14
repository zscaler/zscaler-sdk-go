package geo_coordinates

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
)

func TestGeoCoordinates(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.243.0/24")
	comment := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Create static IP for testing
	staticIP, _, err := staticips.Create(service, &staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	if err != nil {
		t.Fatalf("Error creating static IP for testing: %v", err)
	}

	// Clean up: delete the static IP after test
	defer func() {
		_, err := staticips.Delete(service, staticIP.ID)
		if err != nil {
			t.Errorf("Error deleting static IP: %v", err)
		}
	}()

	// Retrieve the GeoCoordinates using the latitude and longitude from the staticIP
	coordinate, err := GetByGeoCoordinates(service,
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
