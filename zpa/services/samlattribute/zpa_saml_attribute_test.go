package samlattribute

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestSAMLAttribute(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	attributes, _, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(attributes) == 0 {
		t.Errorf("No saml attribute found")
		return
	}
	name := attributes[0].Name
	t.Log("Getting saml attribute by name:" + name)
	attribute, _, err := GetByName(service, name)
	if err != nil {
		t.Errorf("Error getting saml attribute by name: %v", err)
		return
	}
	if attribute.Name != name {
		t.Errorf("identity provider name does not match: expected %s, got %s", name, attribute.Name)
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	attributes, _, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(attributes) == 0 {
		t.Errorf("No saml attributes found")
		return
	}

	// Validate each group
	for _, attribute := range attributes {
		// Checking if essential fields are not empty
		if attribute.ID == "" {
			t.Errorf("saml attributes ID is empty")
		}
		if attribute.Name == "" {
			t.Errorf("saml attributes Name is empty")
		}
	}
}

func TestNonExistentSAMLAttributeName(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)
	_, _, err = GetByName(service, "NonExistentName")
	if err == nil {
		t.Errorf("Expected error when getting non-existent SAML attribute by name, got none")
	}
}

func TestEmptyResponse(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)
	attributes, _, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting SAML attributes: %v", err)
		return
	}
	if attributes == nil {
		t.Errorf("Received nil response for SAML attributes")
		return
	}
}

func TestGetSAMLAttributeByID(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)
	attributes, _, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting all SAML attributes: %v", err)
		return
	}

	if len(attributes) == 0 {
		t.Errorf("No SAML attributes found")
		return
	}

	specificID := attributes[0].ID
	attribute, _, err := Get(service, specificID)
	if err != nil {
		t.Errorf("Error getting SAML attribute by ID: %v", err)
		return
	}
	if attribute.ID != specificID {
		t.Errorf("Mismatch in attribute ID: expected '%s', got %s", specificID, attribute.ID)
		return
	}
}

func TestAllFieldsOfSAMLAttribute(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)
	attributes, _, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting all SAML attributes: %v", err)
		return
	}

	if len(attributes) == 0 {
		t.Errorf("No SAML attributes found")
		return
	}

	specificID := attributes[0].ID
	attribute, _, err := Get(service, specificID)
	if err != nil {
		t.Errorf("Error getting SAML attribute by ID: %v", err)
		return
	}

	// Now check each field
	if attribute.ID == "" {
		t.Errorf("ID is empty")
	}
	if attribute.IdpID == "" {
		t.Errorf("IdpID is empty")
	}
	if attribute.IdpName == "" {
		t.Errorf("IdpName is empty")
	}
	if attribute.Name == "" {
		t.Errorf("Name is empty")
	}
	if attribute.SamlName == "" {
		t.Errorf("SamlName is empty")
	}
}

func TestResponseHeadersAndFormat(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)
	_, resp, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting SAML attributes: %v", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		t.Errorf("Expected content type to start with 'application/json', got %s", contentType)
	}
}
