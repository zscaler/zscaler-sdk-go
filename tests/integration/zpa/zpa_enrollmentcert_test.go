package integration

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/enrollmentcert"
)

func TestEnrollmentCert(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := enrollmentcert.New(client)

	certificates, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting enrollment certificates: %v", err)
		return
	}
	if len(certificates) == 0 {
		t.Errorf("No enrollment certificate found")
		return
	}
	name := certificates[0].Name
	t.Log("Getting enrollment certificate by name:" + name)
	certificate, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting enrollment certificatee by name: %v", err)
		return
	}
	if certificate.Name != name {
		t.Errorf("enrollment certificate name does not match: expected %s, got %s", name, certificate.Name)
		return
	}
}
