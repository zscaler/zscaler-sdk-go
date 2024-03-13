package cloudconnectorgroup

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestIdPController(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	providers, _, err := service.GetAll()
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
	provider, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting identity provider by name: %v", err)
		return
	}
	if provider.Name != name {
		t.Errorf("identity provider name does not match: expected %s, got %s", name, provider.Name)
		return
	}
	// Negative Test: Try to retrieve a Idp with a non-existent name
	nonExistentName := "ThisIdpNameDoesNotExist"
	_, _, err = service.GetByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	providers, _, err := service.GetAll()
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.GetByName("non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
