package customerversionprofile

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestCustomerVersionProfile(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	// Normal case for GetAll
	profiles, _, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting customer version profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No customer version profile found")
		return
	}

	name := profiles[0].Name
	t.Log("Getting customer version profile by name:" + name)

	// Normal case for GetByName
	profile, _, err := GetByName(service, name)
	if err != nil {
		t.Errorf("Error getting customer version profile by name: %v", err)
		return
	}
	if profile.Name != name {
		t.Errorf("customer version profile name does not match: expected %s, got %s", name, profile.Name)
		return
	}

	// Test no version profile found case for GetByName
	t.Run("No version profile found case for GetByName", func(t *testing.T) {
		// Use a name that does not exist
		nonExistentName := "NonExistentVersionProfileName"
		_, _, err := GetByName(service, nonExistentName)
		if err == nil {
			t.Errorf("Expected error when no version profile is found, but got none")
		}
	})

	// Simulate network error by using an invalid URL
	t.Run("Network error case for GetAll", func(t *testing.T) {
		service.Client.Config.CustomerID = "invalid-customer-id"
		_, _, err := GetAll(service)
		if err == nil {
			t.Errorf("Expected network error when calling GetAll with invalid customer ID, but got none")
		}
	})

	// Reset the customer ID after the test
	service.Client.Config.CustomerID = client.Config.CustomerID
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	requiredNames := []string{"New Release", "Default", "Previous Default", "Default - el8"}
	anyVariationSucceeded := false
	var errorMsgs []string

	for _, knownName := range requiredNames {
		// Case variations to test for each knownName
		variations := []string{
			strings.ToUpper(knownName),
			strings.ToLower(knownName),
			cases.Title(language.English).String(knownName),
		}

		for _, variation := range variations {
			t.Run(fmt.Sprintf("GetByName case sensitivity test for %s", variation), func(t *testing.T) {
				t.Logf("Attempting to retrieve customer version profile with name variation: %s", variation)
				version, _, err := GetByName(service, variation)
				if err != nil {
					errorMsg := fmt.Sprintf("Error getting customer version profile with name variation '%s': %v", variation, err)
					errorMsgs = append(errorMsgs, errorMsg)
					return
				}

				// Check if the customer version profile's actual name matches the known name
				if version.Name != knownName {
					errorMsg := fmt.Sprintf("Expected customer version profile name to be '%s' for variation '%s', but got '%s'", knownName, variation, version.Name)
					errorMsgs = append(errorMsgs, errorMsg)
					return
				}

				anyVariationSucceeded = true
			})
		}
	}

	if !anyVariationSucceeded {
		for _, msg := range errorMsgs {
			t.Error(msg)
		}
	}
}
