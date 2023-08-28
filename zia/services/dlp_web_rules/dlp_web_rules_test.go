package dlp_web_rules

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/zia/services/rule_labels"
)

// clean all resources
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
	if !present {
		return true
	}
	shouldClean, err := strconv.ParseBool(val)
	if err != nil {
		return true
	}
	log.Printf("ZSCALER_SDK_TEST_SWEEP value: %v", shouldClean)
	return shouldClean
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
	resources, _ := service.GetAll()
	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		_, _ = service.Delete(r.ID)
	}
}

func TestDLPWebRule(t *testing.T) {
	cleanResources()                // At the start of the test
	defer t.Cleanup(cleanResources) // Will be called at the end

	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// create rule label for testing
	ruleLabelService := rule_labels.New(client)
	ruleLabel, _, err := ruleLabelService.Create(&rule_labels.RuleLabels{
		Name:        name,
		Description: name,
	})
	if err != nil {
		t.Fatalf("Error creating rule label for testing: %v", err)
	}

	// Ensure the rule label is cleaned up at the end of this test
	defer func() {
		_, err := ruleLabelService.Delete(ruleLabel.ID)
		if err != nil {
			t.Errorf("Error deleting rule label: %v", err)
		}
	}()

	service := New(client)
	rule := WebDLPRules{
		Name:                     name,
		Description:              name,
		Order:                    1,
		Rank:                     7,
		State:                    "ENABLED",
		Action:                   "BLOCK",
		OcrEnabled:               true,
		ZscalerIncidentReceiver:  true,
		WithoutContentInspection: false,
		Protocols:                []string{"FTP_RULE", "HTTPS_RULE", "HTTP_RULE"},
		CloudApplications:        []string{"WINDOWS_LIVE_HOTMAIL"},
		FileTypes:                []string{"WINDOWS_META_FORMAT", "BITMAP", "JPEG", "PNG", "TIFF"},
		Labels: []common.IDNameExtensions{
			{
				ID: ruleLabel.ID,
			},
		},
	}

	// Test resource creation
	createdResource, err := service.Create(&rule)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	// Other assertions based on the creation result
	if createdResource.ID == 0 {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Name = updateName
	_, err = service.Update(createdResource.ID, retrievedResource)
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
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedByNameResource, err := service.GetByName(updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedByNameResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedByNameResource.ID)
	}
	if retrievedByNameResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedByNameResource.Name)
	}

	// Test resources retrieval
	allResources, err := service.GetAll()
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(allResources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
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

	// Test resource removal
	_, err = service.Delete(createdResource.ID)
	if err != nil {
		t.Fatalf("Error deleting resource: %v", err)
	}

	// Test resource retrieval after deletion
	_, err = service.Get(createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}
