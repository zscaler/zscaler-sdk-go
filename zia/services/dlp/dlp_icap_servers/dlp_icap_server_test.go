package dlp_icap_servers

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDLPICAPServer_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	servers, err := GetAll(service)
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
	server, err := GetByName(service, name)
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

	for _, server := range servers {
		if server.URL == "" {
			t.Errorf("URL field is empty for server ID %d", server.ID)
		} else if !strings.HasPrefix(server.URL, "icaps://") {
			t.Errorf("Invalid URL format for server ID %d: %s", server.ID, server.URL)
		}

		if server.Status == "" {
			t.Errorf("Status field is empty for server ID %d", server.ID)
		} else if server.Status != "ENABLED" && server.Status != "DISABLED" {
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

	servers, err := GetAll(service)
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

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	// Assuming a group with the name "ZS_BD_ICAP_01" exists
	knownName := "ZS_BD_ICAP_01"

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
