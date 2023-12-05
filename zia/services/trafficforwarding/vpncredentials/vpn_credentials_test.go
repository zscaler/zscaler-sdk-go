package vpncredentials

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
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
	resources, err := service.GetAll()
	if err != nil {
		log.Printf("Error retrieving resources during cleanup: %v", err)
		return
	}

	for _, r := range resources {
		if strings.HasPrefix(r.Comments, "tests-") {
			err := service.Delete(r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func TestTrafficForwardingVPNCreds(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.239.0/24")
	comment := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateComment := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	staticipsService := staticips.New(client)
	staticIP, _, err := staticipsService.Create(&staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	if err != nil {
		t.Fatalf("Creating static ip failed: %v", err)
	}

	defer func() {
		_, err := staticipsService.Delete(staticIP.ID)
		if err != nil {
			t.Errorf("Deleting static ip failed: %v", err)
		}
	}()

	service := New(client)
	cred := VPNCredentials{
		Type:         "IP",
		IPAddress:    ipAddress,
		Comments:     comment,
		PreSharedKey: "newPassword123!",
	}

	var createdResource *VPNCredentials

	err = retryOnConflict(func() error {
		createdResource, _, err = service.Create(&cred)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
	}

	if createdResource.Comments != comment {
		t.Errorf("Expected created resource comment '%s', but got '%s'", comment, createdResource.Comments)
	}

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}

	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}

	if retrievedResource.Comments != comment {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", comment, retrievedResource.Comments)
	}

	retrievedResource.Comments = updateComment
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

	if updatedResource.Comments != updateComment {
		t.Errorf("Expected retrieved updated resource comment '%s', but got '%s'", updateComment, updatedResource.Comments)
	}

	retrievedResource, err = service.GetVPNByType("IP")
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}

	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}

	if retrievedResource.Comments != updateComment {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", updateComment, retrievedResource.Comments)
	}

	resources, err := service.GetAll()
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

	err = retryOnConflict(func() error {
		return service.Delete(createdResource.ID)
	})
	if err != nil {
		t.Fatalf("Error deleting resource: %v", err)
	}

	_, err = service.Get(createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *Service, id int) (*VPNCredentials, error) {
	var resource *VPNCredentials
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
