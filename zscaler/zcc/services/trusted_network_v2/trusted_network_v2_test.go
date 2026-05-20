package trusted_network_v2

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestTrustedNetworkV2(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updatedDNSSearchDomain := "updated-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha) + ".example.com"

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	network := TrustedNetworkV2{
		Name:   name,
		Active: true,

		DNSServerIPs: []string{
			"8.8.8.8",
		},
		DNSSearchDomains: []string{
			"example.com",
		},
		TrustedSubnetIPs: []string{
			"10.0.0.0/24",
		},
		TrustedGatewayIPs: []string{
			"10.0.0.1",
		},
	}

	createdResource, _, err := Create(context.Background(), service, &network)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-zero, but got 0")
	}
	if !strings.EqualFold(createdResource.Name, name) {
		t.Errorf("Expected created trusted network name '%s', but got '%s'", name, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if !strings.EqualFold(retrievedResource.Name, name) {
		t.Errorf("Expected retrieved trusted network name '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.DNSSearchDomains = []string{updatedDNSSearchDomain}
	if _, _, err := Update(context.Background(), service, createdResource.ID, retrievedResource); err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving updated resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if len(updatedResource.DNSSearchDomains) == 0 || updatedResource.DNSSearchDomains[0] != updatedDNSSearchDomain {
		t.Errorf("Expected updated DNS search domain '%s', but got '%v'", updatedDNSSearchDomain, updatedResource.DNSSearchDomains)
	}

	// Test resource partial update (PATCH) — only mutate trustedGatewayIps.
	patchedGateway := "10.0.0.2"
	patchBody := &TrustedNetworkV2{
		TrustedGatewayIPs: []string{patchedGateway},
	}
	patchedResource, _, err := PartialUpdate(context.Background(), service, createdResource.ID, patchBody)
	if err != nil {
		t.Fatalf("Error partially updating resource: %v", err)
	}
	if patchedResource.ID != createdResource.ID {
		t.Errorf("Expected partially updated resource ID '%d', but got '%d'", createdResource.ID, patchedResource.ID)
	}

	patchedFromGet, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving partially updated resource: %v", err)
	}
	if len(patchedFromGet.TrustedGatewayIPs) == 0 || patchedFromGet.TrustedGatewayIPs[0] != patchedGateway {
		t.Errorf("Expected partially updated trusted gateway IP '%s', but got '%v'", patchedGateway, patchedFromGet.TrustedGatewayIPs)
	}
	if len(patchedFromGet.DNSSearchDomains) == 0 || patchedFromGet.DNSSearchDomains[0] != updatedDNSSearchDomain {
		t.Errorf("Expected PATCH to preserve DNS search domain '%s', but got '%v'", updatedDNSSearchDomain, patchedFromGet.DNSSearchDomains)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetByName(context.Background(), service, name)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if len(retrievedResource.DNSSearchDomains) == 0 || retrievedResource.DNSSearchDomains[0] != updatedDNSSearchDomain {
		t.Errorf("Expected retrieved resource DNS search domain '%s', but got '%v'", updatedDNSSearchDomain, retrievedResource.DNSSearchDomains)
	}

	// Test resources retrieval (v2 paginated list)
	resources, err := GetAll(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
	}

	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}

	// Test resource removal
	if _, err := Delete(context.Background(), service, createdResource.ID); err != nil {
		t.Fatalf("Error deleting resource: %v", err)
	}

	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
