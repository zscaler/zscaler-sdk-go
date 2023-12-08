package admins

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/adminuserrolemgmt/roles"
)

// Constants for conflict retries
const (
	maxRetries            = 3
	retryInterval         = 2 * time.Second
	maxConflictRetries    = 5
	conflictRetryInterval = 1 * time.Second
	passwordCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	specialCharacters     = "!@#$%^&*()-_+=<>?"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateComplexPassword(length int) string {
	if length < 4 { // 4 is the minimum to satisfy all the criteria
		length = 12
	}

	password := make([]byte, length)

	// Ensure password meets complexity requirements
	password[0] = byte(passwordCharset[rand.Intn(len(passwordCharset))])
	password[1] = byte(passwordCharset[rand.Intn(26)])                       // Lowercase letter
	password[2] = byte(passwordCharset[rand.Intn(26)+26])                    // Uppercase letter
	password[3] = byte(passwordCharset[rand.Intn(10)+52])                    // Digit
	password[4] = byte(specialCharacters[rand.Intn(len(specialCharacters))]) // Special character

	for i := 5; i < length; i++ {
		password[i] = byte(passwordCharset[rand.Intn(len(passwordCharset))])
	}

	rand.Shuffle(length, func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}

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
	resources, err := service.GetAllAdminUsers()
	if err != nil {
		log.Printf("Error retrieving resources during cleanup: %v", err)
		return
	}

	for _, r := range resources {
		if strings.HasPrefix(r.UserName, "tests-") {
			_, err := service.DeleteAdminUser(r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func TestUserManagement(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateComments := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	email := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	require.NoError(t, err, "Error creating client")

	roleService := roles.New(client)
	roles, err := roleService.GetAllAdminRoles()
	require.NoError(t, err, "Error retrieving roles")
	require.NotEmpty(t, roles, "No roles found")

	// Generate random complex password for admin user account
	rPassword := generateComplexPassword(12)

	service := New(client)
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
		createdResource, err = service.CreateAdminUser(admin)
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
		_, err = service.UpdateAdminUser(createdResource.ID, *retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := service.GetAdminUsers(createdResource.ID)
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
	retrievedResource, err = service.GetAdminUsersByLoginName(expectedLoginName) // Name should be prefixed with "tests-"
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
	resources, err := service.GetAllAdminUsers()
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
		_, delErr := service.DeleteAdminUser(createdResource.ID)
		return delErr
	})
	require.NoError(t, err, "Should not error when deleting")

	_, err = service.GetAdminUsers(createdResource.ID)
	require.Error(t, err, "Expected error retrieving deleted resource")
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *Service, id int) (*AdminUsers, error) {
	var resource *AdminUsers
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = s.GetAdminUsers(id)
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

	_, err = service.GetAdminUsers(0)
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

	_, err = service.DeleteAdminUser(0)
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

	_, err = service.UpdateAdminUser(0, AdminUsers{})
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

	_, err = service.GetAdminByUsername("non-existent-name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
