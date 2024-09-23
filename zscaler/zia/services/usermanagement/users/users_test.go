package users

/*
import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
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
	name := "tests-" + acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	updateComments := "tests-" + acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	email := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	// Step 2: Create the general ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 3: Create a group-specific and department-specific service using the ZIA client
	groupService := groups.New(service.Client)
	departmentService := departments.New(service.Client)
	userService := New(service.Client)

	// Step 4: Retrieve all groups and departments
	groups, err := groupService.GetAllGroups()
	if err != nil || len(groups) == 0 {
		t.Fatalf("Error retrieving groups or no groups found: %v", err)
	}
	departments, err := departmentService.GetAll()
	if err != nil || len(departments) == 0 {
		t.Fatalf("Error retrieving departments or no departments found: %v", err)
	}

	// Step 5: Prepare a random password and create a user payload
	rPassword := tests.TestPassword(20)
	user := Users{
		Name:     name,
		Email:    email + "@bd-hashicorp.com",
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
		createdResource, err = userService.Create(&user)
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

	// Step 7: Test resource retrieval by ID (Instead of Name)
	retrievedResource, err := tryRetrieveResource(userService, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving user by ID: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved user ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != createdResource.Name {
		t.Errorf("Expected retrieved user name '%s', but got '%s'", createdResource.Name, retrievedResource.Name)
	}

	// Step 8: Test resource update
	retrievedResource.Comments = updateComments
	err = retryOnConflict(func() error {
		_, _, err = userService.Update(createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating user: %v", err)
	}

	// Step 9: Verify the update
	updatedResource, err := userService.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving updated user: %v", err)
	}
	if updatedResource.Comments != updateComments {
		t.Errorf("Expected updated user comment '%s', but got '%s'", updateComments, updatedResource.Comments)
	}

	// Step 10: Test retrieving all users (by ID)
	allUsers, err := userService.GetAllUsers()
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
		_, delErr := userService.Delete(createdResource.ID)
		return delErr
	})
	_, err = userService.Get(createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted user, but got nil")
	}
}

// Helper function to retrieve a resource with retry mechanism
func tryRetrieveResource(s *Service, id int) (*Users, error) {
	var resource *Users
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = s.Get(id)
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
