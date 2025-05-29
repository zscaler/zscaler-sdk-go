package enrollmentcert

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestEnrollmentCert(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	certificates, _, err := GetAll(context.Background(), service)
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
			certificate, _, err := GetByName(context.Background(), service, reqName)
			if err != nil {
				t.Fatalf("Error getting enrollment certificate by name %s: %v", reqName, err)
			}
			if certificate.Name != reqName {
				t.Errorf("Enrollment certificate name does not match: expected %s, got %s", reqName, certificate.Name)
			}
		})
	}

	// Additional step: Use the ID of the first certificate to test the Get function
	firstCertID := certificates[0].ID
	t.Run("Get by ID for first certificate", func(t *testing.T) {
		certificateByID, _, err := Get(context.Background(), service, firstCertID)
		if err != nil {
			t.Fatalf("Error getting enrollment certificate by ID %s: %v", firstCertID, err)
		}
		if certificateByID.ID != firstCertID {
			t.Errorf("Enrollment certificate ID does not match: expected %s, got %s", firstCertID, certificateByID.ID)
		}
	})

	// Negative Test: Try to retrieve a certificate with a non-existent ID
	nonExistentID := "non_existent_id"
	t.Run("Get by non-existent ID", func(t *testing.T) {
		_, _, err := Get(context.Background(), service, nonExistentID)
		if err == nil {
			t.Errorf("Expected error when getting by non-existent ID, got nil")
		}
	})
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

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
				certificate, _, err := GetByName(context.Background(), service, variation)
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
