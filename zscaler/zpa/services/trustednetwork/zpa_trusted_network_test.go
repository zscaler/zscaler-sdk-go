package trustednetwork

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestTrustedNetworks(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Test to retrieve all networks
	networks, _, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting trusted networks: %v", err)
		return
	}
	if len(networks) == 0 {
		t.Errorf("No trusted networks found")
		return
	}

	// Additional step: Use the ID of the first certificate to test the Get function
	firstNetworkID := networks[0].ID
	t.Run("Get by ID for first network", func(t *testing.T) {
		networkByID, _, err := Get(context.Background(), service, firstNetworkID)
		if err != nil {
			t.Fatalf("Error getting network by ID %s: %v", firstNetworkID, err)
		}
		if networkByID.ID != firstNetworkID {
			t.Errorf("Enrollment network ID does not match: expected %s, got %s", firstNetworkID, networkByID.ID)
		}
	})

	// Test to retrieve a profile by its name
	name := networks[0].Name
	adaptedName := common.RemoveCloudSuffix(name)
	t.Log("Getting trusted network by name:" + adaptedName)
	profile, _, err := GetByName(context.Background(), service, adaptedName)
	if err != nil {
		t.Errorf("Error getting trusted network by name: %v", err)
		return
	}
	if common.RemoveCloudSuffix(profile.Name) != adaptedName {
		t.Errorf("trusted network name does not match: expected %s, got %s", adaptedName, profile.Name)
		return
	}

	// Negative Test: Try to retrieve a profile with a non-existent name
	nonExistentName := "ThisTrustedNetworkNameDoesNotExist"
	_, _, err = GetByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}

	// Negative Test: Try to retrieve a network with a non-existent ID
	nonExistentID := "non_existent_id"
	t.Run("Get by non-existent ID", func(t *testing.T) {
		_, _, err := Get(context.Background(), service, nonExistentID)
		if err == nil {
			t.Errorf("Expected error when getting by non-existent ID, got nil")
		}
	})
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	networks, _, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting trusted networks: %v", err)
		return
	}
	if len(networks) == 0 {
		t.Errorf("No trusted network found")
		return
	}

	// Validate each network
	for _, network := range networks {
		// Checking if essential fields are not empty
		if network.ID == "" {
			t.Errorf("Trusted network ID is empty")
		}
		if network.Name == "" {
			t.Errorf("Trusted network Name is empty")
		}
		if network.NetworkID == "" {
			t.Errorf("Trusted network UDID is empty")
		}
	}
}

/*
	func TestCaseSensitivityOfGetByName(t *testing.T) {
		client, err := tests.NewOneAPIClient()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		service := services.New(client)

		// Assuming a network with the name "BD-TrustedNetwork01" exists
		knownName := "BD-TrustedNetwork01"

		// Case variations to test
		variations := []string{
			strings.ToUpper(knownName),
			strings.ToLower(knownName),
			cases.Title(language.English).String(knownName),
		}

		for _, variation := range variations {
			t.Logf("Attempting to retrieve trusted network with name variation: %s", variation)
			network, _, err := GetByName(context.Background(), service, variation)
			if err != nil {
				t.Errorf("Error getting trusted network with name variation '%s': %v", variation, err)
				continue
			}

			// Check if the profile's actual name matches the known name
			if common.RemoveCloudSuffix(network.Name) != knownName {
				t.Errorf("Expected trusted network name to be '%s' for variation '%s', but got '%s'", knownName, variation, network.Name)
			}
		}
	}
*/
func TestTrustedNetworkNamesWithSpaces(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Assuming that there are networks with the following name variations
	variations := []string{
		"BD Trusted Network 01", "BD  TrustedNetwork  01", "BD   TrustedNetwork   01",
		"BD    TrustedNetwork01", "BD  TrustedNetwork 01", "BD  Trusted Network   01",
		"BD   Trusted   Network 01",
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve network with name: %s", variation)
		network, _, err := GetByName(context.Background(), service, variation)
		if err != nil {
			t.Errorf("Error getting trusted network with name '%s': %v", variation, err)
			continue
		}

		// Verify if the network's actual name matches the expected variation
		if common.RemoveCloudSuffix(network.Name) != variation {
			t.Errorf("Expected trusted network name to be '%s' but got '%s'", variation, network.Name)
		}
	}
}

func TestTrustedNetworksByNetID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Fetch the list of all Trusted Networks
	networks, _, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting list of trusted networks: %v", err)
		return
	}
	if len(networks) == 0 {
		t.Errorf("No trusted networks found")
		return
	}

	// Assume the first network is the known network for this test
	knownNetwork := networks[0]
	t.Logf("Using known network with Name: %s and NetworkID: %s", knownNetwork.Name, knownNetwork.NetworkID)

	// Use the NetworkID from the known network to test GetByNetID
	networkByID, _, err := GetByNetID(context.Background(), service, knownNetwork.NetworkID)
	if err != nil {
		t.Errorf("Error getting trusted network with NetworkID '%s': %v", knownNetwork.NetworkID, err)
		return
	}

	// Check if the network's actual NetworkID matches the known NetworkID
	if networkByID.NetworkID != knownNetwork.NetworkID {
		t.Errorf("Expected trusted network NetworkID to be '%s', but got '%s'", knownNetwork.NetworkID, networkByID.NetworkID)
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
