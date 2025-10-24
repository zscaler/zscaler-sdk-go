package idpcontroller

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestIdPController(t *testing.T) {
	// service, err := tests.NewOneAPIClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	service, err := tests.NewZPAClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	providers, _, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting identity providers: %v", err)
		return
	}
	if len(providers) == 0 {
		t.Errorf("No identity provider found")
		return
	}

	name := providers[0].Name
	t.Log("Getting identity provider by name:" + name)
	provider, _, err := GetByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting identity provider by name: %v", err)
		return
	}
	if provider.Name != name {
		t.Errorf("identity provider name does not match: expected %s, got %s", name, provider.Name)
		return
	}

	// Additional step: Use the ID of the first provider to test the Get function
	firstProviderID := providers[0].ID
	t.Log("Getting identity provider by ID:" + firstProviderID)
	providerByID, _, err := Get(context.Background(), service, firstProviderID)
	if err != nil {
		t.Errorf("Error getting identity provider by ID: %v", err)
		return
	}
	if providerByID.ID != firstProviderID {
		t.Errorf("identity provider ID does not match: expected %s, got %s", firstProviderID, providerByID.ID)
		return
	}

	// Negative Test: Try to retrieve a Idp with a non-existent name
	nonExistentName := "ThisIdpNameDoesNotExist"
	_, _, err = GetByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}

	// Negative Test: Try to retrieve an Idp with a non-existent ID
	nonExistentID := "non_existent_id"
	_, _, err = Get(context.Background(), service, nonExistentID)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent ID, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	// service, err := tests.NewOneAPIClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	service, err := tests.NewZPAClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	providers, _, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting identity provider: %v", err)
		return
	}
	if len(providers) == 0 {
		t.Errorf("No identity provider found")
		return
	}

	// Validate each group
	for _, provider := range providers {
		// Checking if essential fields are not empty
		if provider.ID == "" {
			t.Errorf("Identity provider ID is empty")
		}
		if provider.Name == "" {
			t.Errorf("Identity provider Name is empty")
		}
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	// service, err := tests.NewOneAPIClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	service, err := tests.NewZPAClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
