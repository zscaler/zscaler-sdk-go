package intermediatecacertificates

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestIntermediateCertificate_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Test 1: GetAll
	certificates, err := service.GetAll()
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
	certificateByName, err := service.GetByName(name)
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
	certificateByID, err := service.GetCertificate(certID)
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
	readyToUseCerts, err := service.GetIntCAReadyToUse()
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
