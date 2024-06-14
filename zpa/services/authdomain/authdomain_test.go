package authdomain

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestGetAllAuthDomains(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := services.New(client)

	domains, resp, err := GetAllAuthDomains(service)
	if err != nil {
		t.Fatalf("Failed to fetch authentication domains: %v", err)
	}
	if resp == nil {
		t.Fatalf("Expected a non-nil response")
	}
	if len(domains.AuthDomains) == 0 {
		t.Fatalf("Expected to retrieve at least one authentication domain, but got none")
	}

	t.Logf("Retrieved %d authentication domains", len(domains.AuthDomains))
	for _, domain := range domains.AuthDomains {
		t.Logf("Auth Domain: %s", domain)
	}
}
