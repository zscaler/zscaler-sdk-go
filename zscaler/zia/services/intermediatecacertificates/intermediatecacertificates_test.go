package intermediatecacertificates

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestIntermediateCertificate_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Test 1: GetAll
	certificates, err := GetAll(service)
	if err != nil {
		t.Fatalf("Error getting intermediate certificates: %v", err)
		return
	}
	if len(certificates) == 0 {
		t.Fatal("No intermediate certificate found")
		return
	}
	t.Logf("Found %d intermediate certificates", len(certificates))

	// Test 2: GetByName
	name := certificates[0].Name
	certificateByName, err := GetByName(service, name)
	if err != nil {
		t.Fatalf("Error getting intermediate certificate by name %s: %v", name, err)
		return
	}
	if certificateByName.Name != name {
		t.Errorf("Intermediate certificate name mismatch: expected %s, got %s", name, certificateByName.Name)
		return
	}
	t.Logf("Successfully retrieved certificate by name: %s", name)

	// Test 3: GetCertificate
	certID := certificates[0].ID
	certificateByID, err := GetCertificate(service, certID)
	if err != nil {
		t.Fatalf("Error getting intermediate certificate by ID %d: %v", certID, err)
		return
	}
	if certificateByID.ID != certID {
		t.Errorf("Intermediate certificate ID mismatch: expected %d, got %d", certID, certificateByID.ID)
		return
	}
	t.Logf("Successfully retrieved certificate by ID: %d", certID)

	// Test 4: GetIntCAReadyToUse
	readyToUseCerts, err := GetIntCAReadyToUse(service)
	if err != nil {
		t.Fatalf("Error getting intermediate CA ready to use: %v", err)
		return
	}
	if len(readyToUseCerts) == 0 {
		t.Fatal("No ready-to-use intermediate certificate found")
		return
	}
	t.Logf("Successfully retrieved ready-to-use intermediate certificate: %s", readyToUseCerts[0].Name)
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Assuming a certificate with the name "Zscaler Intermediate CA Certificate" exists
	knownName := "Zscaler Intermediate CA Certificate"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		certificate, err := GetByName(service, variation)
		if err != nil {
			t.Errorf("Error getting certificate with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the certificate's actual name matches the known name
		if certificate.Name != knownName {
			t.Errorf("Expected certificate name to be '%s' for variation '%s', but got '%s'", knownName, variation, certificate.Name)
		}
	}
}
