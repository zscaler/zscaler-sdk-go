package dlpdictionaries

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
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

func TestDLPDictionaries(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "dlpdictionaries", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	name := tests.GetTestName("tests-dlpdict")
	updateName := tests.GetTestName("tests-dlpdict")

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

	var createdResource *DlpDictionary

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, _, err = Create(context.Background(), service, &dictionary)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created dlp dictionary '%s', but got '%s'", name, createdResource.Name)
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
		t.Errorf("Expected retrieved dlp dictionary '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update - create a clean update object
	updateDict := &DlpDictionary{
		ID:                    retrievedResource.ID,
		Name:                  updateName,
		Description:           updateName,
		DictionaryType:        retrievedResource.DictionaryType,
		CustomPhraseMatchType: retrievedResource.CustomPhraseMatchType,
		Phrases:               retrievedResource.Phrases,
		Patterns:              retrievedResource.Patterns,
	}
	err = retryOnConflict(func() error {
		_, _, err = Update(context.Background(), service, createdResource.ID, updateDict)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name (use updated name)
	retrievedResource, err = GetByName(context.Background(), service, updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedResource.Name)
	}

	// Test resources retrieval
	resources, err := GetAll(context.Background(), service)
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
		_, delErr := DeleteDlpDictionary(context.Background(), service, createdResource.ID)
		return delErr
	})
	if err != nil {
		t.Fatalf("Error deleting resource: %v", err)
	}

	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}

	// Test predefined identifiers
	t.Run("Predefined Identifiers", func(t *testing.T) {
		dictionaryName := "CRED_LEAKAGE"
		identifiers, dictionaryID, err := GetPredefinedIdentifiers(context.Background(), service, dictionaryName)
		require.NoError(t, err)
		assert.NotZero(t, dictionaryID)
		assert.NotEmpty(t, identifiers)
		fmt.Printf("Dictionary ID: %d\n", dictionaryID)
		fmt.Printf("Predefined Identifiers: %v\n", identifiers)
	})

	// Test error cases
	t.Run("Error fetching dictionary by name", func(t *testing.T) {
		_, _, err := GetPredefinedIdentifiers(context.Background(), service, "InvalidDictionaryName")
		assert.Error(t, err)
	})

	t.Run("Retrieve non-existent resource", func(t *testing.T) {
		_, err := Get(context.Background(), service, 999999999)
		if err == nil {
			t.Error("Expected error retrieving non-existent resource, but got nil")
		}
	})

	t.Run("Delete non-existent resource", func(t *testing.T) {
		_, err := DeleteDlpDictionary(context.Background(), service, 999999999)
		if err == nil {
			t.Error("Expected error deleting non-existent resource, but got nil")
		}
	})

	t.Run("Update non-existent resource", func(t *testing.T) {
		_, _, err := Update(context.Background(), service, 999999999, &DlpDictionary{})
		if err == nil {
			t.Error("Expected error updating non-existent resource, but got nil")
		}
	})

	t.Run("GetByName non-existent resource", func(t *testing.T) {
		_, err := GetByName(context.Background(), service, "non_existent_name")
		if err == nil {
			t.Error("Expected error retrieving resource by non-existent name, but got nil")
		}
	})
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id int) (*DlpDictionary, error) {
	var resource *DlpDictionary
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(context.Background(), s, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}
