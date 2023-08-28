package dlpdictionaries

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
)

// clean all resources
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources() // clean up at the beginning
}

func teardown() {
	cleanResources() // clean up at the end
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	if !present {
		return true // default value
	}
	shouldClean, err := strconv.ParseBool(val)
	if err != nil {
		return true // default to cleaning if the value is not parseable
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
		_, _ = service.DeleteDlpDictionary(r.ID)
	}
}

func TestDLPDictionaries(t *testing.T) {
	cleanResources()                // At the start of the test
	defer t.Cleanup(cleanResources) // Will be called at the end

	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	dictionary := DlpDictionary{
		Name:                  name,
		Description:           name,
		DictionaryType:        "PATTERNS_AND_PHRASES",
		CustomPhraseMatchType: "MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY",
		Phrases: []Phrases{
			{
				Action: "PHRASE_COUNT_TYPE_ALL",
				Phrase: "YourPhrase",
			},
		},
		Patterns: []Patterns{
			{
				Action:  "PATTERN_COUNT_TYPE_UNIQUE",
				Pattern: "YourPattern",
			},
		},
	}

	// Test resource creation
	createdResource, _, err := service.Create(&dictionary)
	// Check if the request was successful
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 || createdResource.Name != name {
		t.Fatalf("Creation: Expected resource with ID and name '%s', got ID %d and name '%s'", name, createdResource.ID, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, err := service.Get(createdResource.ID)
	if err != nil || retrievedResource.ID != createdResource.ID || retrievedResource.Name != name {
		t.Errorf("Retrieval: Expected resource with ID %d and name '%s', got ID %d and name '%s'", createdResource.ID, name, retrievedResource.ID, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Name = updateName
	_, _, err = service.Update(createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Update: Error updating resource: %v", err)
	}
	updatedResource, err := service.Get(createdResource.ID)
	if err != nil || updatedResource.Name != updateName {
		t.Errorf("Update: Expected updated resource with name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, err = service.GetByName(updateName)
	if err != nil || retrievedResource.ID != createdResource.ID || retrievedResource.Name != updateName {
		t.Errorf("GetByName: Expected resource with ID %d and name '%s', got ID %d and name '%s'", createdResource.ID, updateName, retrievedResource.ID, retrievedResource.Name)
	}

	// Resources Retrieval
	resources, err := service.GetAll()
	if err != nil || len(resources) == 0 {
		t.Errorf("GetAll: Error retrieving resources or received empty slice: %v", err)
	}
	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetAll: Expected resources to contain ID %d, but it didn't", createdResource.ID)
	}

	// Resource Removal
	_, err = service.DeleteDlpDictionary(createdResource.ID)
	if err != nil {
		t.Errorf("Delete: Error deleting resource: %v", err)
	}

	// Confirm Deletion
	_, err = service.Get(createdResource.ID)
	if err == nil {
		t.Errorf("Post-Delete Retrieval: Expected error retrieving deleted resource, but got nil")
	}
}
