package bacertificate

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestBACertificates(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	certificates, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting browser access certificates: %v", err)
		return
	}
	if len(certificates) == 0 {
		t.Errorf("No browser access certificate found")
		return
	}
	name := certificates[0].Name
	t.Log("Getting browser access certificate by name:" + name)
	certificate, _, err := service.GetIssuedByName(name)
	if err != nil {
		t.Errorf("Error getting browser access certificate by name: %v", err)
		return
	}
	if certificate.Name != name {
		t.Errorf("browser access certificate name does not match: expected %s, got %s", name, certificate.Name)
		return
	}
}

func TestCertificatesExpiration(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	certificates, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting browser access certificates: %v", err)
		return
	}

	// Flag to track if we find any expired certificate
	anyExpired := false

	// Iterate over each certificate and check the status
	for _, cert := range certificates {
		if cert.Status == "Expired" {
			anyExpired = true
			break
		}
	}

	if anyExpired {
		t.Errorf("Found an expired browser access certificate.")
		return
	}

	// If no expired certificates are found, log a success message and complete the test.
	t.Log("No expired browser access certificates found.")
}
