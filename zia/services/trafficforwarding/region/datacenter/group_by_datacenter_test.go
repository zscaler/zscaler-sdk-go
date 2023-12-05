package datacenter

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
)

func TestGroupByDatacenter(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.241.0/24")
	comment := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	staticipsService := staticips.New(client)
	staticIP, _, err := staticipsService.Create(&staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	if err != nil {
		t.Fatalf("Creating static ip failed: %v", err)
	}

	defer func() {
		_, err := staticipsService.Delete(staticIP.ID)
		if err != nil {
			t.Errorf("Deleting static ip failed: %v", err)
		}
	}()

	// Create common parameters with the source IP address
	commonParams := common.DatacenterSearchParameters{SourceIp: ipAddress}

	// Test for each individual search parameter
	t.Run("TestRoutableIP", func(t *testing.T) {
		commonParams.RoutableIP = true
		results, err := service.SearchByDatacenters(commonParams)
		if err != nil {
			t.Errorf("Error searching datacenters with RoutableIP: %v", err)
		}
		if len(results) == 0 {
			t.Errorf("Expected non-zero results for RoutableIP search")
		}
	})

	t.Run("TestWithinCountryOnly", func(t *testing.T) {
		commonParams.WithinCountryOnly = true
		results, err := service.SearchByDatacenters(commonParams)
		if err != nil {
			t.Errorf("Error searching datacenters with RoutableIP: %v", err)
		}
		if len(results) == 0 {
			t.Errorf("Expected non-zero results for RoutableIP search")
		}
	})

	t.Run("TestIncludePrivateServiceEdge", func(t *testing.T) {
		commonParams.IncludePrivateServiceEdge = true
		results, err := service.SearchByDatacenters(commonParams)
		if err != nil {
			t.Errorf("Error searching datacenters with RoutableIP: %v", err)
		}
		if len(results) == 0 {
			t.Errorf("Expected non-zero results for RoutableIP search")
		}
	})

	t.Run("TestIncludeCurrentVips", func(t *testing.T) {
		commonParams.IncludeCurrentVips = true
		results, err := service.SearchByDatacenters(commonParams)
		if err != nil {
			t.Errorf("Error searching datacenters with RoutableIP: %v", err)
		}
		if len(results) == 0 {
			t.Errorf("Expected non-zero results for RoutableIP search")
		}
	})

	// Test for each individual search parameter
	t.Run("TestSourceIp", func(t *testing.T) {
		results, err := service.SearchByDatacenters(common.DatacenterSearchParameters{SourceIp: ipAddress})
		if err != nil {
			t.Errorf("Error searching datacenters with SourceIp: %v", err)
		}
		if len(results) == 0 {
			t.Errorf("Expected non-zero results for SourceIp search")
		}
	})

	t.Run("TestLatitudeLongitude", func(t *testing.T) {
		// Adjust to include source IP, latitude, and longitude from staticIP
		results, err := service.SearchByDatacenters(common.DatacenterSearchParameters{
			SourceIp:  ipAddress,
			Latitude:  float64(staticIP.Latitude),
			Longitude: float64(staticIP.Longitude),
		})

		if err != nil {
			t.Errorf("Error searching datacenters with Latitude/Longitude: %v", err)
		}
		if len(results) == 0 {
			t.Errorf("Expected non-zero results for Latitude/Longitude search")
		}
	})

	// Test with all parameters combined
	t.Run("TestAllParameters", func(t *testing.T) {
		results, err := service.SearchByDatacenters(common.DatacenterSearchParameters{
			RoutableIP:                true,
			WithinCountryOnly:         true,
			IncludePrivateServiceEdge: true,
			IncludeCurrentVips:        true,
			SourceIp:                  ipAddress,
			Latitude:                  float64(staticIP.Latitude),
			Longitude:                 float64(staticIP.Longitude),
		})

		if err != nil {
			t.Errorf("Error searching datacenters with all parameters: %v", err)
		}
		if len(results) == 0 {
			t.Errorf("Expected non-zero results for combined search")
		}
		// Additional assertions can be added here as per requirement
	})
}
