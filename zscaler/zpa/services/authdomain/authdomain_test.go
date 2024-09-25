package authdomain

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestGetAllAuthDomains(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	domains, resp, err := GetAllAuthDomains(context.Background(), service)
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
