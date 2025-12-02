package risk_profiles

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	maxRetries    = 3
	retryInterval = 2 * time.Second
)

// Constants for conflict retries
const (
	maxConflictRetries    = 5
	conflictRetryInterval = 1 * time.Second
)

func retryOnConflict(operation func() error) error {
	var lastErr error
	for i := 0; i < maxConflictRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

func TestRiskProfiles(t *testing.T) {
	tests.ResetTestNameCounter()
	name := tests.GetTestName("tests-riskprof")
	updateName := tests.GetTestName("tests-riskprof")

	client, err := tests.NewVCRTestClient(t, "risk_profiles", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	profile := RiskProfiles{
		ProfileName:               name,
		Status:                    "SANCTIONED",
		RiskIndex:                 []int{1, 2, 3, 4, 5},
		PasswordStrength:          "GOOD",
		PoorItemsOfService:        "YES",
		AdminAuditLogs:            "YES",
		DataBreach:                "YES",
		SourceIpRestrictions:      "YES",
		FileSharing:               "YES",
		MfaSupport:                "YES",
		SslPinned:                 "YES",
		Certifications:            []string{"AICPA", "CCPA", "CISP"},
		DataEncryptionInTransit:   []string{"SSLV2", "SSLV3", "TLSV1_0", "TLSV1_1", "TLSV1_2", "TLSV1_3", "UN_KNOWN"},
		HttpSecurityHeaders:       "YES",
		Evasive:                   "YES",
		DnsCaaPolicy:              "YES",
		SslCertValidity:           "YES",
		WeakCipherSupport:         "YES",
		Vulnerability:             "YES",
		VulnerableToHeartBleed:    "YES",
		SslCertKeySize:            "BITS_2048",
		VulnerableToPoodle:        "YES",
		SupportForWaf:             "YES",
		VulnerabilityDisclosure:   "YES",
		DomainKeysIdentifiedMail:  "YES",
		MalwareScanningForContent: "YES",
		DomainBasedMessageAuth:    "YES",
		SenderPolicyFramework:     "YES",
		RemoteScreenSharing:       "YES",
		VulnerableToLogJam:        "YES",
		ProfileType:               "CLOUD_APPLICATIONS",
	}

	var createdResource *RiskProfiles

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, _, err = Create(context.Background(), service, &profile)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	// Other assertions based on the creation result
	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.ProfileName != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.ProfileName)
	}

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.ProfileName != name {
		t.Errorf("Expected retrieved bandwidth control classes '%s', but got '%s'", name, retrievedResource.ProfileName)
	}

	// Test resource update
	retrievedResource.ProfileName = updateName
	err = retryOnConflict(func() error {
		_, _, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.ProfileName != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.ProfileName)
	}

	// Test resource retrieval by name
	retrievedByNameResource, err := GetByName(context.Background(), service, updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedByNameResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedByNameResource.ID)
	}
	if retrievedByNameResource.ProfileName != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedByNameResource.ProfileName)
	}

	// Test resources retrieval
	allResources, err := GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(allResources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
	}

	// check if the created resource is in the list
	found := false
	for _, resource := range allResources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}

	// Introduce a delay before deleting
	time.Sleep(5 * time.Second) // sleep for 5 seconds

	// Test resource removal
	err = retryOnConflict(func() error {
		_, getErr := Get(context.Background(), service, createdResource.ID)
		if getErr != nil {
			return fmt.Errorf("Resource %d may have already been deleted: %v", createdResource.ID, getErr)
		}
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id int) (*RiskProfiles, error) {
	var resource *RiskProfiles
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(context.Background(), s, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func TestRetrieveNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "risk_profiles", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = Get(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "risk_profiles", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = Delete(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "risk_profiles", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, _, err = Update(context.Background(), service, 0, &RiskProfiles{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "risk_profiles", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
