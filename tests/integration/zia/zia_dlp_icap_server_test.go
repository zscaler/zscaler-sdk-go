package integration

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_icap_servers"
)

func TestDLPICAPServer_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := dlp_icap_servers.New(client)

	servers, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting icap servers: %v", err)
		return
	}
	if len(servers) == 0 {
		t.Errorf("No icap servers found")
		return
	}
	name := servers[0].Name
	t.Log("Getting icap servers by name:" + name)
	server, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting icap servers by name: %v", err)
		return
	}
	if server.Name != name {
		t.Errorf("icap server name does not match: expected %s, got %s", name, server.Name)
		return
	}
}
