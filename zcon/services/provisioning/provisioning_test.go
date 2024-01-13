package provisioning

/*
import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestZConProvisioningAPIKey(t *testing.T) {
	client, err := tests.NewZConClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	// Test 1: Retrieve all API keys
	apiKeys, err := service.GetAll()
	if err != nil {
		t.Fatalf("Error getting api keys: %v", err)
	}
	if len(apiKeys) == 0 {
		t.Fatal("No api keys found")
	}

	// Selecting a key for further tests
	testKey := apiKeys[0]

	// Test 2: GetPartnerAPIKey with includePartnerKey set to true and false
	_, err = service.GetPartnerAPIKey(testKey.KeyValue, true)
	if err != nil {
		t.Errorf("Error getting api key with includePartnerKey true: %v", err)
	}

	_, err = service.GetPartnerAPIKey(testKey.KeyValue, false)
	if err != nil {
		t.Errorf("Error getting api key with includePartnerKey false: %v", err)
	}

	/*
		// Test 3: Regenerate API key
		regeneratedKey, err := service.Create(nil, false, &testKey.ID)
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
