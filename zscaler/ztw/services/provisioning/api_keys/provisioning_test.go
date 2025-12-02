package api_keys

/*
import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon/services"
)

func TestZConProvisioningAPIKey(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewZConClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Test 1: Retrieve all API keys
	apiKeys, err := GetAll(service)
	if err != nil {
		t.Fatalf("Error getting api keys: %v", err)
	}
	if len(apiKeys) == 0 {
		t.Fatal("No api keys found")
	}

	// Selecting a key for further tests
	testKey := apiKeys[0]

	// Test 2: GetPartnerAPIKey with includePartnerKey set to true and false
	_, err = GetPartnerAPIKey(service, testKey.KeyValue, true)
	if err != nil {
		t.Errorf("Error getting api key with includePartnerKey true: %v", err)
	}

	_, err = GetPartnerAPIKey(service, testKey.KeyValue, false)
	if err != nil {
		t.Errorf("Error getting api key with includePartnerKey false: %v", err)
	}

	// Test 3: Regenerate API key
	regeneratedKey, err := Create(service, nil, false, &testKey.ID)
	if err != nil {
		t.Fatalf("Error regenerating api key: %v", err)
	}

	// Check that a new keyValue is returned
	if regeneratedKey.KeyValue == testKey.KeyValue {
		t.Errorf("API key was not regenerated: expected a different keyValue")
	}

	// Optionally, check if the lastModifiedTime has been updated
	if regeneratedKey.LastModifiedTime <= testKey.LastModifiedTime {
		t.Errorf("API key regeneration did not update lastModifiedTime")
	}
}
*/
