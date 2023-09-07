package idpcontroller

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
}
