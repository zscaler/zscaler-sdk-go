package virtualipaddress

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/region/datacenter"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
)

func TestVIPs(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.244.0/24")
	comment := "tests-" + acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	staticIP, _, err := staticips.Create(service, &staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	if err != nil {
		t.Fatalf("Error creating static IP for testing: %v", err)
	}

	defer func() {
		_, err := staticips.Delete(service, staticIP.ID)
		if err != nil {
			t.Errorf("Error deleting static IP: %v", err)
		}
	}()

	// Search for datacenters by source IP
	dataCenterList, err := datacenter.SearchByDatacenters(service, common.DatacenterSearchParameters{
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
		vips, err := GetZscalerVIPs(service, datacenterName)
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
		vips, err := GetZSGREVirtualIPList(service, staticIP.IpAddress, 3)
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

		pairVips, err := GetPairZSGREVirtualIPsWithinCountry(service, sourceIP, countryCode)
		if err != nil {
			t.Fatalf("Error fetching pair of VIPs within country: %v", err)
		}

		if len(*pairVips) < 2 {
			t.Errorf("Expected at least a pair of VIPs, got %d", len(*pairVips))
		}

		// Additional checks can be added based on the structure of GREVirtualIPList
	})

	// Test for GetAll
	t.Run("TestGetAll", func(t *testing.T) {
		allVips, err := GetAll(service, ipAddress)
		if err != nil {
			t.Fatalf("Error fetching all VIPs for source IP: %v", err)
		}
		if len(allVips) == 0 {
			t.Errorf("Expected VIPs for source IP %s, got none", ipAddress)
		}
	})

	// Test for getAllStaticIPs
	t.Run("TestGetAllStaticIPs", func(t *testing.T) {
		staticIPs, err := getAllStaticIPs(service)
		if err != nil {
			t.Fatalf("Error fetching all static IPs: %v", err)
		}
		if len(staticIPs) == 0 {
			t.Errorf("Expected static IPs, got none")
		}
	})

	// Test for GetAllSourceIPs
	t.Run("TestGetAllSourceIPs", func(t *testing.T) {
		allSourceIPs, err := GetAllSourceIPs(service)
		if err != nil {
			t.Fatalf("Error fetching all source IPs: %v", err)
		}
		if len(allSourceIPs) == 0 {
			t.Errorf("Expected VIPs for all source IPs, got none")
		}
	})
}
