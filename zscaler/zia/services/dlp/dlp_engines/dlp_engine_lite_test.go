package dlp_engines

/*
import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDLPEngineLite_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	engines, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting predefined engine name: %v", err)
		return
	}
	if len(engines) == 0 {
		t.Errorf("No predefined engine name found")
		return
	}
	name := engines[0].PredefinedEngineName
	t.Log("Getting predefined engine name by name:" + name)
	engine, err := GetByPredefinedEngine(service, name)
	if err != nil {
		t.Errorf("Error getting predefined engine by name: %v", err)
		return
	}
	if engine.PredefinedEngineName != name {
		t.Errorf("predefined engine name does not match: expected %s, got %s", name, engine.PredefinedEngineName)
		return
	}
	// Negative Test: Try to retrieve an predefined engine name with a non-existent name
	nonExistentName := "ThisPredefinedEngineDoesNotExist"
	_, err = GetByPredefinedEngine(service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestGetById(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Get all engines to find a valid ID
	engines, err := GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error getting all dlp predefined engine: %v", err)
	}
	if len(engines) == 0 {
		t.Fatalf("No dlp predefined engines found for testing")
	}

	// Choose the first engines's ID for testing
	testID := engines[0].ID

	// Retrieve the engine by ID
	engine, err := Get(context.Background(), service, testID)
	if err != nil {
		t.Errorf("Error retrieving dlp predefined engine with ID %d: %v", testID, err)
		return
	}

	// Verify the retrieved engine
	if engine == nil {
		t.Errorf("No engine returned for ID %d", testID)
		return
	}

	if engine.ID != testID {
		t.Errorf("Retrieved engine ID mismatch: expected %d, got %d", testID, engine.ID)
	}
}

func TestGetByLiteID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Retrieve all DLP engines
	engines, err := GetAllEngineLite(service)
	if err != nil {
		t.Fatalf("Error getting all DLP engines: %v", err)
	}
	if len(engines) == 0 {
		t.Log("No DLP engines found")
		return
	}

	// Use the first engine for testing
	firstEngine := engines[0]

	// Test GetEngineLiteID
	t.Run("GetEngineLiteID", func(t *testing.T) {
		engineByID, err := GetEngineLiteID(service, firstEngine.ID)
		if err != nil {
			t.Errorf("Error getting DLP engine by ID: %v", err)
			return
		}
		if engineByID == nil || engineByID.ID != firstEngine.ID {
			t.Errorf("DLP engine ID does not match: expected %d, got %d", firstEngine.ID, engineByID.ID)
		}
	})
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	engines, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting predefined engine: %v", err)
		return
	}
	if len(engines) == 0 {
		t.Errorf("No predefined engine found")
		return
	}

	// Validate predefined engine
	for _, engine := range engines {
		// Checking if essential fields are not empty
		if engine.ID == 0 {
			t.Errorf("predefined engine ID is empty")
		}
		if !engine.CustomDlpEngine && engine.PredefinedEngineName == "" {
			t.Errorf("predefined engine Name is empty for predefined engine with ID: %d", engine.ID)
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Assuming a group with the name "EXTERNAL" exists
	knownName := "EXTERNAL"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		engine, err := GetByPredefinedEngine(service, variation)
		if err != nil {
			t.Errorf("Error getting predefined engine with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if engine.PredefinedEngineName != knownName {
			t.Errorf("Expected group name to be '%s' for variation '%s', but got '%s'", knownName, variation, engine.PredefinedEngineName)
		}
	}
}
*/
