package virtualipaddress

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/region/datacenter"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
)

func TestVIPs(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.244.0/24")
	comment := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	staticIPService := staticips.New(client)
	staticIP, _, err := staticIPService.Create(&staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})

	if err != nil {
		t.Fatalf("Error creating static IP for testing: %v", err)
	}

	defer func() {
		_, err := staticIPService.Delete(staticIP.ID)
		if err != nil {
			t.Errorf("Error deleting static IP: %v", err)
		}
	}()

	// Search for datacenters by source IP
	dataCenterService := datacenter.New(client)
	dataCenterList, err := dataCenterService.SearchByDatacenters(common.DatacenterSearchParameters{
		SourceIp: ipAddress,
	})
	if err != nil {
		t.Errorf("Error searching datacenter by SourceIp: %v", err)
		return
	}
	if len(dataCenterList) == 0 {
		t.Error("Expected retrieved datacenter to be non-empty, but got empty slice")
		return
	}

	// Test for GetZscalerVIPs
	t.Run("TestGetZscalerVIPs", func(t *testing.T) {
		// Check if there are any datacenters available
		if len(dataCenterList) == 0 {
			t.Fatal("No datacenters available for testing")
		}

		// Use the name of the first datacenter from the list
		datacenterName := dataCenterList[0].Datacenter.Name
		vips, err := service.GetZscalerVIPs(datacenterName)
		if err != nil {
			t.Fatalf("Error fetching VIPs for datacenter %s: %v", datacenterName, err)
		}

		// Assuming vips contains a slice of VIP objects directly
		if vips == nil || len(vips.DataCenter) == 0 {
			t.Errorf("Expected VIPs for datacenter %s, got none", datacenterName)
		}

	})

	// Test for GetZSGREVirtualIPList
	t.Run("TestGetZSGREVirtualIPList", func(t *testing.T) {
		vips, err := service.GetZSGREVirtualIPList(staticIP.IpAddress, 3)
		if err != nil {
			t.Fatalf("Error fetching GRE VIP list: %v", err)
		}
		if len(*vips) < 3 {
			t.Errorf("Expected at least 3 VIPs, got %d", len(*vips))
		}
	})

	// Test for GetPairZSGREVirtualIPsWithinCountry
	t.Run("TestGetPairZSGREVirtualIPsWithinCountry", func(t *testing.T) {
		sourceIP := ipAddress // Assuming ipAddress from the staticIP
		countryCode := "US"   // Replace with the appropriate country code

		pairVips, err := service.GetPairZSGREVirtualIPsWithinCountry(sourceIP, countryCode)
		if err != nil {
			t.Fatalf("Error fetching pair of VIPs within country: %v", err)
		}

		if len(*pairVips) < 2 {
			t.Errorf("Expected at least a pair of VIPs, got %d", len(*pairVips))
		}

		// Additional checks can be added based on the structure of GREVirtualIPList
	})

}
