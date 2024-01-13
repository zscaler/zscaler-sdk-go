package enrollmentcert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestEnrollmentCert(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	service := New(client)

	certificates, _, err := service.GetAll()
	if err != nil {
		t.Fatalf("Error getting enrollment certificates: %v", err)
		return
	}
	if len(certificates) == 0 {
		t.Fatalf("No enrollment certificate found")
		return
	}

	// Check if GetAll returns specific certificate names
	requiredNames := []string{"Root", "Client", "Connector", "Service Edge", "Isolation Client"}
	for _, reqName := range requiredNames {
		found := false
		for _, cert := range certificates {
			if cert.Name == reqName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected certificate with name %s not found in GetAll response", reqName)
		}
	}

	// Test GetByName for each specific certificate name
	for _, reqName := range requiredNames {
		t.Run("GetByName for "+reqName, func(t *testing.T) {
			certificate, _, err := service.GetByName(reqName)
			if err != nil {
				t.Fatalf("Error getting enrollment certificate by name %s: %v", reqName, err)
			}
			if certificate.Name != reqName {
				t.Errorf("Enrollment certificate name does not match: expected %s, got %s", reqName, certificate.Name)
			}
		})
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

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	service := New(client)

	requiredNames := []string{"Root", "Client", "Connector", "Service Edge", "Isolation Client"}

	for _, knownName := range requiredNames {
		// Case variations to test for each knownName
		variations := []string{
			strings.ToUpper(knownName),
			strings.ToLower(knownName),
			cases.Title(language.English).String(knownName),
		}

		for _, variation := range variations {
			t.Run(fmt.Sprintf("GetByName case sensitivity test for %s", variation), func(t *testing.T) {
				t.Logf("Attempting to retrieve certificate with name variation: %s", variation)
				certificate, _, err := service.GetByName(variation)
				if err != nil {
					t.Errorf("Error getting certificate with name variation '%s': %v", variation, err)
					return
				}

				// Check if the certificate's actual name matches the known name
				if certificate.Name != knownName {
					t.Errorf("Expected certificate name to be '%s' for variation '%s', but got '%s'", knownName, variation, certificate.Name)
				}
			})
		}
	}
}
