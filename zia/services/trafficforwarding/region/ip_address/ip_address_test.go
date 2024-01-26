package ip_address

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
)

func TestByIPAddress(t *testing.T) {
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

	result, err := service.GetByIPAddress(ipAddress)
	if err != nil {
		t.Fatalf("Error searching by IP address: %v", err)
	}

	if result == nil {
		t.Errorf("Expected results for IP address search, but got nil")
	}

	// Additional assertions can be added here as needed
}
