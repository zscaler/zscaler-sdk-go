package admins

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/adminuserrolemgmt/roles"
)

// Constants for conflict retries
const (
	maxRetries            = 3
	retryInterval         = 2 * time.Second
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

func TestUserManagement(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateComments := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	email := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	require.NoError(t, err, "Error creating client")

	service := services.New(client)

	roles, err := roles.GetAllAdminRoles(service)

	require.NoError(t, err, "Error retrieving roles")
	require.NotEmpty(t, roles, "No roles found")

	// Generate random complex password for admin user account
	rPassword := tests.TestPassword(20)

	admin := AdminUsers{
		UserName:                    name + name,
		LoginName:                   name + "@bd-hashicorp.com",
		Email:                       email + "@bd-hashicorp.com",
		Comments:                    updateComments,
		Password:                    rPassword,
		IsPasswordLoginAllowed:      true,
		IsSecurityReportCommEnabled: true,
		IsServiceUpdateCommEnabled:  true,
		IsProductUpdateCommEnabled:  true,
		IsPasswordExpired:           false,
		IsExecMobileAppEnabled:      false,
		Role: &Role{
			ID: roles[0].ID, // Associating the first role for simplicity
		},
	}

	var createdResource *AdminUsers
	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = CreateAdminUser(service, admin)
		require.NoError(t, err, "Creating a new admin user should not error")
		return err
	})
	require.NoError(t, err, "Error making POST request")
	require.NotZero(t, createdResource.ID, "Expected created resource ID to be non-empty")

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
	}
	expectedLoginName := name + "@bd-hashicorp.com"
	if createdResource.LoginName != expectedLoginName {
		t.Errorf("Expected created admin user '%s', but got '%s'", expectedLoginName, createdResource.LoginName)
	}

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.LoginName != expectedLoginName {
		t.Errorf("Expected retrieved user '%s', but got '%s'", expectedLoginName, retrievedResource.LoginName)
	}

	// Test resource update
	retrievedResource.Comments = updateComments
	err = retryOnConflict(func() error {
		_, err = UpdateAdminUser(service, createdResource.ID, *retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := GetAdminUsers(service, createdResource.ID)
	require.NoError(t, err, "Could not get admin user by ID")
	assert.NotNil(t, updatedResource.Disabled, "admin user disabled is missing")

	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Comments != updateComments {
		t.Errorf("Expected retrieved updated resource comment '%s', but got '%s'", updateComments, updatedResource.Comments)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetAdminUsersByLoginName(service, expectedLoginName) // Name should be prefixed with "tests-"
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Comments != updateComments {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", updateComments, createdResource.Comments)
	}
	// Test resources retrieval
	resources, err := GetAllAdminUsers(service)
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
	}
	// check if the created resource is in the list
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
	err = retryOnConflict(func() error {
		_, delErr := DeleteAdminUser(service, createdResource.ID)
		return delErr
	})
	require.NoError(t, err, "Should not error when deleting")

	// Confirm that the user has been deleted
	_, err = GetAdminUsers(service, createdResource.ID)
	if err != nil {
		if strings.Contains(err.Error(), "resource not found") || strings.Contains(err.Error(), "does not exist") {
			// User deletion confirmed, no further operations on this user
			log.Println("User deletion confirmed. No further operations will be performed on this user.")
		} else {
			t.Fatalf("Unexpected error after deletion: %v", err)
		}
	} else {
		t.Fatal("User still exists after deletion")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *services.Service, id int) (*AdminUsers, error) {
	var resource *AdminUsers
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = GetAdminUsers(s, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = GetAdminUsers(service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = DeleteAdminUser(service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = UpdateAdminUser(service, 0, AdminUsers{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = GetAdminByUsername(service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
