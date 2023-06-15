package integration

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/bacertificate"
)

func TestBACertificates(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := bacertificate.New(client)

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
