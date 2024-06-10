package policysetcontroller

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestNonExistentResourceOperations(t *testing.T) {
	policyTypes := []string{
		"ACCESS_POLICY", "TIMEOUT_POLICY", "CLIENT_FORWARDING_POLICY", "ISOLATION_POLICY",
		"INSPECTION_POLICY", "CREDENTIAL_POLICY", "CAPABILITIES_POLICY",
		"CLIENTLESS_SESSION_PROTECTION_POLICY", "REDIRECTION_POLICY", "SIEM_POLICY",
	}

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client) // Ensure this New() function is correctly pointing to your service constructor

	for _, policyType := range policyTypes {
		t.Run(policyType, func(t *testing.T) {
			// Testing retrieve non-existent resource
			_, _, err = GetPolicyRule(service, policyType, "non_existent_id")
			if err == nil {
				t.Errorf("Expected error retrieving non-existent resource for policyType %s, but got nil", policyType)
			}

			// Testing delete non-existent resource
			_, err = Delete(service, policyType, "non_existent_id")
			if err == nil {
				t.Errorf("Expected error deleting non-existent resource for policyType %s, but got nil", policyType)
			}

			// Testing update non-existent resource
			_, err = UpdateRule(service, policyType, "non_existent_id", &PolicyRule{})
			if err == nil {
				t.Errorf("Expected error updating non-existent resource for policyType %s, but got nil", policyType)
			}

			// Testing get by name non-existent resource
			_, _, err = GetByNameAndType(service, policyType, "non_existent_name")
			if err == nil {
				t.Errorf("Expected error retrieving resource by non-existent name for policyType %s, but got nil", policyType)
			}
		})
	}
}
