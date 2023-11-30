package dlp_icap_servers

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestDLPICAPServer_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	servers, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting icap servers: %v", err)
		return
	}
	if len(servers) == 0 {
		t.Errorf("No icap server found")
		return
	}
	name := servers[0].Name
	t.Log("Getting icap server by name:" + name)
	server, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting icap server by name: %v", err)
		return
	}
	if server.Name != name {
		t.Errorf("icap server name does not match: expected %s, got %s", name, server.Name)
		return
	}
	// Negative Test: Try to retrieve an icap server with a non-existent name
	nonExistentName := "ThisIcapServerDoesNotExist"
	_, err = service.GetByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	servers, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting icap server: %v", err)
		return
	}
	if len(servers) == 0 {
		t.Errorf("No icap server found")
		return
	}

	// Validate icap server
	for _, server := range servers {
		// Checking if essential fields are not empty
		if server.ID == 0 {
			t.Errorf("icap server ID is empty")
		}
		if server.Name == "" {
			t.Errorf("icap server Name is empty")
		}
	}
}
