package users

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/departments"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/groups"
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

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources()
}

func teardown() {
	cleanResources()
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true")) // simplified for clarity
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	client, err := tests.NewZiaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)
	resources, err := service.GetAllUsers()
	if err != nil {
		log.Printf("Error retrieving resources during cleanup: %v", err)
		return
	}

	for _, r := range resources {
		if strings.HasPrefix(r.Name, "tests-") {
			_, err := service.Delete(r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func TestUserManagement(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	updateComments := "tests-" + acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	email := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rPassword := acctest.RandString(10)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}
	// For groups, you might want to make sure you get a list of groups instead, as per your structure.
	// Assuming a function similar to GetAll for departments exists for groups:
	groupService := groups.New(client)
	groups, err := groupService.GetAllGroups()
	if err != nil || len(groups) == 0 {
		t.Fatalf("Error retrieving groups or no groups found: %v", err)
	}
	// Retrieve the first department and group (for simplicity; you can enhance this to retrieve by name or other criteria)
	departmentService := departments.New(client)
	departments, err := departmentService.GetAll()
	if err != nil || len(departments) == 0 {
		t.Fatalf("Error retrieving departments or no departments found: %v", err)
	}
	service := New(client)

	user := Users{
		Name:     name,
		Email:    email + "@bd-hashicorp.com",
		Password: rPassword,
		Comments: updateComments,
		Groups: []common.IDNameExtensions{
			{
				ID: groups[0].ID, // For simplicity, we are associating the first group; you can extend this to multiple groups
			},
		},
		Department: &common.UserDepartment{
			ID: departments[0].ID, // Associating the first department for simplicity
		},
	}

	var createdResource *Users
	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = service.Create(&user)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created user '%s', but got '%s'", name, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved user '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Comments = updateComments
	err = retryOnConflict(func() error {
		_, _, err = service.Update(createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := service.Get(createdResource.ID)
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
	retrievedResource, err = service.GetUserByName(name) // Name should be prefixed with "tests-"
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
	resources, err := service.GetAllUsers()
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
		_, delErr := service.Delete(createdResource.ID)
		return delErr
	})
	_, err = service.Get(createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *Service, id int) (*Users, error) {
	var resource *Users
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = s.Get(id)
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
	service := New(client)

	_, err = service.Get(0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Delete(0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.Update(0, &Users{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.GetUserByName("non-existent-name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
