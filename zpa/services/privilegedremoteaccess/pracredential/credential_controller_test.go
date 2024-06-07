package pracredential

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/authdomain"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/microtenants"
)

func TestCredentialController(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	service := New(client)

	credController := Credential{
		Name:           name,
		Description:    name,
		CredentialType: "USERNAME_PASSWORD",
		UserName:       name,
		Password:       tests.TestPassword(10), // Ensuring the password length is within constraints
		UserDomain:     "acme.com",
	}

	// Test resource creation
	createdResource, _, err := service.Create(&credController)

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Name = updateName
	retrievedResource.Password = tests.TestPassword(10) // Ensure the password is not being reset during update

	_, err = service.Update(createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, _, err = service.GetByName(updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
	}

	// Test resources retrieval
	resources, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	}

	// Check if the created resource is in the list
	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	}

	// Test resource removal
	_, err = service.Delete(createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = service.Get(createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.Get("non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Delete("non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Update("non_existent_id", &Credential{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.GetByName("non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}

// /////// Credential Move to Microtenant Test /////
func TestPRACredentialMove(t *testing.T) {
	// Generate base random strings
	baseName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	baseDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 1: Create a new Microtenant
	authDomainService := authdomain.New(client)
	authDomainList, _, err := authDomainService.GetAllAuthDomains()
	if err != nil {
		t.Errorf("Error getting auth domains: %v", err)
		return
	}
	if len(authDomainList.AuthDomains) == 0 {
		t.Error("Expected retrieved auth domains to be non-empty, but got empty slice")
		return
	}

	microtenantService := microtenants.New(client)
	newMicrotenant := microtenants.MicroTenant{
		Name:                       baseName + "-microtenant",
		Description:                baseDescription + "-microtenant",
		Enabled:                    true,
		PrivilegedApprovalsEnabled: true,
		CriteriaAttribute:          "AuthDomain",
		CriteriaAttributeValues:    []string{authDomainList.AuthDomains[0]},
	}
	createdMicrotenant, _, err := microtenantService.Create(newMicrotenant)
	if err != nil {
		t.Fatalf("Failed to create microtenant: %v", err)
	}

	// Ensure the microtenant is deleted at the end of the test
	defer func() {
		_, err := microtenantService.Delete(createdMicrotenant.ID)
		if err != nil {
			t.Errorf("Error deleting microtenant: %v", err)
		}
	}()

	microtenantID := createdMicrotenant.ID

	// Step 2: Create a local credential
	credentialService := New(client)
	credController := Credential{
		Name:           baseName + "-credential",
		Description:    baseDescription + "-credential",
		CredentialType: "USERNAME_PASSWORD",
		UserName:       baseName + "-user",
		Password:       tests.TestPassword(10), // Ensuring the password length is within constraints
		UserDomain:     "acme.com",
	}

	createdCredential, _, err := credentialService.Create(&credController)
	if err != nil {
		t.Fatalf("Failed to create credential: %v", err)
	}

	// Step 3: Move the credential to the microtenant
	resp, err := credentialService.CredentialMove(createdCredential.ID, microtenantID)
	if err != nil {
		t.Fatalf("Error moving credential to microtenant: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Failed to move credential to microtenant, status code: %d", resp.StatusCode)
	}
}
