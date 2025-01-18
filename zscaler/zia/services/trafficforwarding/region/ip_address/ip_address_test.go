package ip_address

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
)

func TestByIPAddress(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.244.0/24")
	comment := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	staticIP, _, err := staticips.Create(context.Background(), service, &staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	if err != nil {
		t.Fatalf("Error creating static IP for testing: %v", err)
	}

	defer func() {
		_, err := staticips.Delete(context.Background(), service, staticIP.ID)
		if err != nil {
			t.Errorf("Error deleting static IP: %v", err)
		}
	}()

	result, err := GetByIPAddress(context.Background(), service, ipAddress)
	if err != nil {
		t.Fatalf("Error searching by IP address: %v", err)
	}

	if result == nil {
		t.Errorf("Expected results for IP address search, but got nil")
	}

	// Additional assertions can be added here as needed
}
