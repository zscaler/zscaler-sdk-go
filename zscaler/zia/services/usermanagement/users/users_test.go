package users

/*
import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/departments"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/groups"
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

func TestUserManagement(t *testing.T) {
	// Step 1: Create a random user name and other test data
	name := "tests-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	updateComments := "tests-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	email := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	// Step 2: Create the general ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	departments, err := departments.GetAll(context.Background(), service)
	if err != nil || len(departments) == 0 {
		t.Fatalf("Error retrieving departments or no departments found: %v", err)
	}

		groups, err := groups.GetAllGroups(context.Background(), service, nil)
	if err != nil || len(departments) == 0 {
		t.Fatalf("Error retrieving departments or no departments found: %v", err)
	}

	// Step 5: Prepare a random password and create a user payload
	rPassword := tests.TestPassword(20)
	user := Users{
		Name:     name,
		Email:    email + "@securitygeek.io",
		Password: rPassword,
		Comments: updateComments,
		Groups: []common.IDNameExtensions{
			{
				ID: groups[0].ID, // Associating the first group for simplicity
			},
		},
		Department: &common.UserDepartment{
			ID: departments[0].ID, // Associating the first department for simplicity
		},
	}

	// Step 6: Test resource creation
	var createdResource *Users

	err = retryOnConflict(func() error {
		createdResource, err = Create(context.Background(), service, &user)
		return err
	})
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	time.Sleep(5 * time.Second)

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created user name '%s', but got '%s'", name, createdResource.Name)
	}

	//Step 7: Test resource retrieval by ID (Instead of Name)
	retrievedResource, err := tryRetrieveResource(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving user by ID: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved user ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}

	// Log the names for debugging purposes
	t.Logf("Created user name (raw): '%s'", createdResource.Name)
	t.Logf("Retrieved user name (raw): '%s'", retrievedResource.Name)

	// Compare the raw names directly
	if retrievedResource.Name != createdResource.Name {
		t.Errorf("Expected retrieved user name '%s', but got '%s'", createdResource.Name, retrievedResource.Name)
	}

	//Step 8: Test resource update
	retrievedResource.Comments = updateComments
	err = retryOnConflict(func() error {
		_, _, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating user: %v", err)
	}

	// Step 9: Verify the update
	updatedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving updated user: %v", err)
	}
	if updatedResource.Comments != updateComments {
		t.Errorf("Expected updated user comment '%s', but got '%s'", updateComments, updatedResource.Comments)
	}

	// Step 10: Test retrieving all users (by ID)
	allUsers, err := GetAllUsers(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving all users: %v", err)
	}
	if len(allUsers) == 0 {
		t.Fatal("Expected non-empty list of users, but got none")
	}
	found := false
	for _, u := range allUsers {
		if u.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected created user to be in the list, but it wasn't")
	}

	// Step 11: Test resource removal by ID
	err = retryOnConflict(func() error {
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted user, but got nil")
	}
}

// Helper function to retrieve a resource with retry mechanism
func tryRetrieveResource(ctx context.Context, service *zscaler.Service, id int) (*Users, error) {
	var resource *Users
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(ctx, service, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}

		// Handle RESOURCE_NOT_FOUND errors by retrying
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			log.Printf("Attempt %d: User not found yet, retrying in %v...", i+1, retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}
*/
