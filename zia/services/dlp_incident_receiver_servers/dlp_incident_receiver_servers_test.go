package dlp_incident_receiver_servers

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestDLPIncidentReceiver_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	receivers, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting incident receivers : %v", err)
		return
	}
	if len(receivers) == 0 {
		t.Errorf("No incident receivers found")
		return
	}
	name := receivers[0].Name
	t.Log("Getting incident receiver by name:" + name)
	receiver, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting incident receiver by name: %v", err)
		return
	}
	if receiver.Name != name {
		t.Errorf("icap server name does not match: expected %s, got %s", name, receiver.Name)
		return
	}
	// Negative Test: Try to retrieve an incident receiver with a non-existent name
	nonExistentName := "ThisIncidentReceiverDoesNotExist"
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

	receivers, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting incident receiver: %v", err)
		return
	}
	if len(receivers) == 0 {
		t.Errorf("No incident receiver found")
		return
	}

	// Validate incident receiver
	for _, receiver := range receivers {
		// Checking if essential fields are not empty
		if receiver.ID == 0 {
			t.Errorf("incident receiver ID is empty")
		}
		if receiver.Name == "" {
			t.Errorf("incident receiver Name is empty")
		}
	}
}
