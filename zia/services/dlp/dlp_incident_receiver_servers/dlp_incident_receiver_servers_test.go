package dlp_incident_receiver_servers

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDLPIncidentReceiver_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	receivers, err := GetAll(service)
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
	receiver, err := GetByName(service, name)
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
	_, err = GetByName(service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestGetById(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Get all servers to find a valid ID
	servers, err := GetAll(service)
	if err != nil {
		t.Fatalf("Error getting all icap servers: %v", err)
	}
	if len(servers) == 0 {
		t.Fatalf("No icap servers found for testing")
	}

	// Choose the first server's ID for testing
	testID := servers[0].ID

	// Retrieve the server by ID
	server, err := Get(service, testID)
	if err != nil {
		t.Errorf("Error retrieving icap server with ID %d: %v", testID, err)
		return
	}

	// Verify the retrieved server
	if server == nil {
		t.Errorf("No server returned for ID %d", testID)
		return
	}

	if server.ID != testID {
		t.Errorf("Retrieved server ID mismatch: expected %d, got %d", testID, server.ID)
	}
}

func TestURLAndStatusFields(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Retrieve all servers
	servers, err := GetAll(service)
	if err != nil {
		t.Fatalf("Error getting all icap servers: %v", err)
	}
	if len(servers) == 0 {
		t.Fatalf("No icap servers found for testing")
	}

	// Iterate through each server and check URL and Status fields
	for _, server := range servers {
		// Check if URL field is populated and valid
		if server.URL == "" {
			t.Errorf("URL field is empty for server ID %d", server.ID)
		} else if !strings.HasPrefix(server.URL, "icaps://") { // Adjust this condition based on your URL format
			t.Errorf("Invalid URL format for server ID %d: %s", server.ID, server.URL)
		}

		// Check if Status field is populated and valid
		if server.Status == "" {
			t.Errorf("Status field is empty for server ID %d", server.ID)
		} else if server.Status != "ENABLED" && server.Status != "DISABLED" { // Assuming possible statuses
			t.Errorf("Invalid status for server ID %d: %s", server.ID, server.Status)
		}
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	receivers, err := GetAll(service)
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

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	// Assuming a group with the name "ZS_BD_INC_RECEIVER_01" exists
	knownName := "ZS_BD_INC_RECEIVER_01"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		server, err := GetByName(service, variation)
		if err != nil {
			t.Errorf("Error getting icap server with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if server.Name != knownName {
			t.Errorf("Expected group name to be '%s' for variation '%s', but got '%s'", knownName, variation, server.Name)
		}
	}
}
