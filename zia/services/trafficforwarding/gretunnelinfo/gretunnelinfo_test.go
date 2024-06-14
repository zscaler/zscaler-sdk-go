package gretunnelinfo

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
	virtualipaddress "github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/virtualipaddress"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources()
}

func teardown() {
	cleanResources()
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true")) // simplified for clarity
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	client, err := tests.NewZiaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)
	resources, err := gretunnels.GetAll(service)
	if err != nil {
		log.Printf("Error retrieving resources during cleanup: %v", err)
		return
	}

	for _, r := range resources {
		if strings.HasPrefix(r.SourceIP, "tests-") {
			_, err := gretunnels.DeleteGreTunnels(service, r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func TestGRETunnelInfo(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.248.0/24")
	comment := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	// create static IP for testing
	service := services.New(client)
	staticIP, _, err := staticips.Create(service, &staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	// Check if the request was successful
	if err != nil {
		t.Fatalf("Error creating static IP for testing: %v", err)
	}
	defer deleteStaticIP(service, staticIP.ID, t)

	vipRecommendedList, err := virtualipaddress.GetAll(service, ipAddress)
	if err != nil {
		t.Errorf("Error getting recommended vip: %v", err)
		return
	}
	if len(vipRecommendedList) == 0 {
		t.Error("Expected retrieved recommended vip to be non-empty, but got empty slice")
	}

	withinCountry := true // Create a boolean variable

	greTunnel, _, err := gretunnels.CreateGreTunnels(service, &gretunnels.GreTunnels{
		SourceIP:      staticIP.IpAddress,
		Comment:       comment,
		WithinCountry: &withinCountry,
		IPUnnumbered:  true,
		PrimaryDestVip: &gretunnels.PrimaryDestVip{
			ID:        vipRecommendedList[0].ID,
			VirtualIP: vipRecommendedList[0].VirtualIp,
		},
		SecondaryDestVip: &gretunnels.SecondaryDestVip{
			ID:        vipRecommendedList[1].ID,
			VirtualIP: vipRecommendedList[1].VirtualIp,
		},
	})
	if err != nil {
		t.Fatalf("Error creating GRE tunnel: %v", err)
	}

	defer deleteGRETunnel(service, greTunnel.ID, t)

	// Get GRE tunnel information

	greTunnelInfo, err := GetGRETunnelInfo(service, ipAddress)
	if err != nil {
		t.Fatalf("Error retrieving GRE tunnel info: %v", err)
	}
	if greTunnelInfo.IPaddress != ipAddress {
		t.Errorf("Expected IP address %s, got %s", ipAddress, greTunnelInfo.IPaddress)
	}
	// Detailed Validation of GRE Tunnel Info
	t.Run("TestDetailedValidationOfGRETunnelInfo", func(t *testing.T) {
		if greTunnelInfo.GREEnabled == true {
			t.Errorf("Expected GREEnabled to be true")
		}
		if greTunnelInfo.PrimaryGW == "" {
			t.Errorf("Expected PrimaryGW to be non-empty")
		}
		if greTunnelInfo.SecondaryGW == "" {
			t.Errorf("Expected SecondaryGW to be non-empty")
		}
		// Add more checks as necessary for other fields like GRERangePrimary, GRERangeSecondary, etc.
	})

	t.Run("TestInvalidGRETunnelRetrieval", func(t *testing.T) {
		invalidIpAddress := "invalid-ip-address"
		_, err := GetGRETunnelInfo(service, invalidIpAddress)
		if err == nil {
			t.Errorf("Expected an error for invalid IP address, but got nil")
		} else {
			t.Logf("Received expected error for invalid IP address: %v", err)
		}
	})
}

// deleteStaticIP deletes a static IP resource
func deleteStaticIP(service *services.Service, id int, t *testing.T) {
	_, err := staticips.Delete(service, id)
	if err != nil {
		t.Errorf("Error deleting static IP: %v", err)
	}
}

// deleteGRETunnel deletes a GRE tunnel resource
func deleteGRETunnel(service *services.Service, id int, t *testing.T) {
	_, err := gretunnels.DeleteGreTunnels(service, id)
	if err != nil {
		t.Errorf("Error deleting GRE tunnel: %v", err)
	}
}
