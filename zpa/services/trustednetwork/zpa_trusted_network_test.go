package trustednetwork

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestTrustedNetworks(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	nets, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting trusted networks: %v", err)
		return
	}
	if len(nets) == 0 {
		t.Errorf("No trusted networks found")
		return
	}
	name := nets[0].Name
	t.Log("Getting trusted network by name:" + name)
	net, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting trusted network by name: %v", err)
		return
	}
	if net.Name != name {
		t.Errorf("Trusted network name does not match: expected %s, got %s", name, net.Name)
		return
	}
}
