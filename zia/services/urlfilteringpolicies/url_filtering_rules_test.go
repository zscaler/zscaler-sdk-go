package urlfilteringpolicies

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/cloudbrowserisolation"
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
		if strings.HasPrefix(r.Name, "tests-") {
			_, err := service.Delete(r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func tryRetrieveResource(service *Service, id int) (*URLFilteringRule, error) {
	var resource *URLFilteringRule
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = service.Get(id)
		if err == nil {
			return resource, nil
		}
		time.Sleep(retryInterval)
	}
	return nil, err
}

func TestURLFilteringRuleIsolation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-updated-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	cbiService := cloudbrowserisolation.New(client)
	cbiProfileList, err := cbiService.GetAll()
	if err != nil {
		t.Errorf("Error getting cbi profile: %v", err)
		return
	}
	if len(cbiProfileList) == 0 {
		t.Error("Expected retrieved cbi profile to be non-empty, but got empty slice")
	}

	rule := URLFilteringRule{
		Name:           name,
		Description:    name,
		Order:          1,
		Rank:           7,
		State:          "ENABLED",
		Action:         "ISOLATE",
		URLCategories:  []string{"ANY"},
		Protocols:      []string{"HTTPS_RULE", "HTTP_RULE"},
		RequestMethods: []string{"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"},
		UserAgentTypes: []string{"OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"},
		CBIProfileID:   91223046,
		CBIProfile: CBIProfile{
			ProfileSeq: 0,
			ID:         cbiProfileList[0].ID,
			Name:       cbiProfileList[0].Name,
			URL:        cbiProfileList[0].URL,
		},
	}

	// Create a URL Filtering Rule
	createdResource, err := createURLFilteringRule(t, service, &rule)
	if err != nil {
		t.Fatalf("Error creating URL Filtering Rule: %v", err)
	}

	defer cleanupURLFilteringRule(t, service, createdResource.ID)

	// Retrieve and check the URL Filtering Rule
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving URL Filtering Rule: %v", err)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected name '%s', got '%s'", name, retrievedResource.Name)
	}

	// Update the URL Filtering Rule
	retrievedResource.Name = updateName
	err = updateURLFilteringRule(t, service, retrievedResource)
	if err != nil {
		t.Fatalf("Error updating URL Filtering Rule: %v", err)
	}

	// Retrieve and check the updated URL Filtering Rule
	updatedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving updated URL Filtering Rule: %v", err)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected updated name '%s', got '%s'", updateName, updatedResource.Name)
	}
}

// / Testing URL Filtering Rule with BLOCK action
func TestURLFilteringRuleBlock(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-updated-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	rule := URLFilteringRule{
		Name:           name,
		Description:    name,
		Order:          1,
		Rank:           7,
		State:          "ENABLED",
		Action:         "BLOCK",
		URLCategories:  []string{"ANY"},
		Protocols:      []string{"ANY_RULE"},
		RequestMethods: []string{"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"},
	}

	// Create a URL Filtering Rule
	createdResource, err := createURLFilteringRule(t, service, &rule)
	if err != nil {
		t.Fatalf("Error creating URL Filtering Rule: %v", err)
	}

	defer cleanupURLFilteringRule(t, service, createdResource.ID)

	// Retrieve and check the URL Filtering Rule
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving URL Filtering Rule: %v", err)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected name '%s', got '%s'", name, retrievedResource.Name)
	}

	// Update the URL Filtering Rule
	retrievedResource.Name = updateName
	err = updateURLFilteringRule(t, service, retrievedResource)
	if err != nil {
		t.Fatalf("Error updating URL Filtering Rule: %v", err)
	}

	// Retrieve and check the updated URL Filtering Rule
	updatedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving updated URL Filtering Rule: %v", err)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected updated name '%s', got '%s'", updateName, updatedResource.Name)
	}
}

func createURLFilteringRule(t *testing.T, service *Service, rule *URLFilteringRule) (*URLFilteringRule, error) {
	var createdResource *URLFilteringRule
	err := retryOnConflict(func() error {
		var err error
		createdResource, err = service.Create(rule)
		return err
	})
	return createdResource, err
}

func updateURLFilteringRule(t *testing.T, service *Service, resource *URLFilteringRule) error {
	return retryOnConflict(func() error {
		_, _, err := service.Update(resource.ID, resource)
		return err
	})
}

func cleanupURLFilteringRule(t *testing.T, service *Service, id int) {
	err := retryOnConflict(func() error {
		_, err := service.Delete(id)
		return err
	})
	if err != nil {
		t.Errorf("Failed to cleanup URL Filtering Rule: %v", err)
	}
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

	_, _, err = service.Update(0, &URLFilteringRule{})
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

	_, err = service.GetByName("non-existent-name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
