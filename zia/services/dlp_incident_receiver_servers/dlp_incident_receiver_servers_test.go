package dlp_incident_receiver_servers

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestDLPIncidentReceiver_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	servers, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting incident receiver servers: %v", err)
		return
	}
	if len(servers) == 0 {
		t.Errorf("No incident receiver servers found")
		return
	}
	name := servers[0].Name
	t.Log("Getting incident receiver servers by name:" + name)
	server, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting incident receiver servers by name: %v", err)
		return
	}
	if server.Name != name {
		t.Errorf("incident receiver servers name does not match: expected %s, got %s", name, server.Name)
		return
	}
}
