package enrollmentcert

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
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
	requiredNames := []string{"Root", "Client", "Connector", "Service Edge"}
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
