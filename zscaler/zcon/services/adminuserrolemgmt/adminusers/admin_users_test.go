package adminusers

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon/services/adminuserrolemgmt/adminroles"
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

	client, err := tests.NewZConClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	resources, err := GetAllAdminUsers(service)
	if err != nil {
		log.Printf("Error retrieving resources during cleanup: %v", err)
		return
	}

	for _, r := range resources {
		if strings.HasPrefix(r.UserName, "tests-") {
			_, err := DeleteAdminUser(service, r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func TestZCONUserManagement(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateComments := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	email := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZConClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	roles, err := adminroles.GetAllAdminRoles(service)
	if err != nil || len(roles) == 0 {
		t.Fatalf("Error retrieving roles or no roles found: %v", err)
	}
	// Generate random complex password for admin user account
	rPassword := generateComplexPassword(12)

	admin := AdminUsers{
		UserName:                    name + name,
		LoginName:                   name + "@securitygeek.io",
		Email:                       email + "@securitygeek.io",
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
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
	}
	expectedLoginName := name + "@securitygeek.io"
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
	_, err = GetAdminUsers(service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
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
