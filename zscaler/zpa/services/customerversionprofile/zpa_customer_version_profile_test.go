package customerversionprofile

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestCustomerVersionProfile(t *testing.T) {
	// service, err := tests.NewOneAPIClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	service, err := tests.NewZPAClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Normal case for GetAll
	profiles, _, err := GetAll(context.Background(), service)
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
	profile, _, err := GetByName(context.Background(), service, name)
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
		_, _, err := GetByName(context.Background(), service, nonExistentName)
		if err == nil {
			t.Errorf("Expected error when no version profile is found, but got none")
		}
	})
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	// service, err := tests.NewOneAPIClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	service, err := tests.NewZPAClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

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
				version, _, err := GetByName(context.Background(), service, variation)
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
