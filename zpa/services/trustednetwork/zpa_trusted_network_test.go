package trustednetwork

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestTrustedNetworks(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Assuming a network with the name "BD-TrustedNetwork03" exists
	knownName := "BDTrustedNetwork"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve network with name variation: %s", variation)
		network, _, err := service.GetByName(variation)
		if err != nil {
			t.Errorf("Error getting trusted network with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the network's actual name matches the known name
		if common.RemoveCloudSuffix(network.Name) != knownName {
			t.Errorf("Expected trusted network name to be '%s' for variation '%s', but got '%s'", knownName, variation, network.Name)
		}
	}
}

func TestTrustedNetworkNamesWithSpaces(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Assuming that there are networks with the following name variations
	variations := []string{
		"BD Trusted Network 01", "BD  TrustedNetwork  01", "BD   TrustedNetwork   01",
		"BD    TrustedNetwork01", "BD  TrustedNetwork 01", "BD  Trusted Network   01",
		"BD   Trusted   Network 01",
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve network with name: %s", variation)
		network, _, err := service.GetByName(variation)
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Fetch the list of all Trusted Networks
	networks, _, err := service.GetAll()
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
	networkByID, _, err := service.GetByNetID(knownNetwork.NetworkID)
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.GetByName("non-existent-name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
