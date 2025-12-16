// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlpdictionaries"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDLPDictionaries_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	dictID := 12345
	path := "/zia/api/v1/dlpDictionaries/12345"

	server.On("GET", path, common.SuccessResponse(dlpdictionaries.DlpDictionary{
		ID:             dictID,
		Name:           "Custom SSN Dictionary",
		Description:    "Detects SSN patterns",
		Custom:         true,
		DictionaryType: "PATTERNS_AND_PHRASES",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlpdictionaries.Get(context.Background(), service, dictID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, dictID, result.ID)
	assert.Equal(t, "Custom SSN Dictionary", result.Name)
	assert.True(t, result.Custom)
}

func TestDLPDictionaries_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dlpDictionaries"

	server.On("GET", path, common.SuccessResponse([]dlpdictionaries.DlpDictionary{
		{ID: 1, Name: "SSN", Custom: false, DictionaryType: "PREDEFINED"},
		{ID: 2, Name: "Credit Card", Custom: false, DictionaryType: "PREDEFINED"},
		{ID: 3, Name: "Custom PII", Custom: true, DictionaryType: "PATTERNS_AND_PHRASES"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlpdictionaries.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestDLPDictionaries_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dlpDictionaries"

	server.On("POST", path, common.SuccessResponse(dlpdictionaries.DlpDictionary{
		ID:             99999,
		Name:           "New Dictionary",
		Custom:         true,
		DictionaryType: "PATTERNS_AND_PHRASES",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newDict := &dlpdictionaries.DlpDictionary{
		Name:           "New Dictionary",
		DictionaryType: "PATTERNS_AND_PHRASES",
	}

	result, _, err := dlpdictionaries.Create(context.Background(), service, newDict)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestDLPDictionaries_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	dictID := 12345
	path := "/zia/api/v1/dlpDictionaries/12345"

	server.On("PUT", path, common.SuccessResponse(dlpdictionaries.DlpDictionary{
		ID:   dictID,
		Name: "Updated Dictionary",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateDict := &dlpdictionaries.DlpDictionary{
		ID:   dictID,
		Name: "Updated Dictionary",
	}

	result, _, err := dlpdictionaries.Update(context.Background(), service, dictID, updateDict)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Dictionary", result.Name)
}

func TestDLPDictionaries_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	dictID := 12345
	path := "/zia/api/v1/dlpDictionaries/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = dlpdictionaries.DeleteDlpDictionary(context.Background(), service, dictID)

	require.NoError(t, err)
}

func TestDLPDictionaries_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	dictName := "Custom SSN Dictionary"
	path := "/zia/api/v1/dlpDictionaries"

	server.On("GET", path, common.SuccessResponse([]dlpdictionaries.DlpDictionary{
		{ID: 1, Name: "Other Dictionary", Custom: true},
		{ID: 2, Name: dictName, Custom: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlpdictionaries.GetByName(context.Background(), service, dictName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, dictName, result.Name)
}

func TestDLPDictionaries_GetPredefinedIdentifiers_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	dictName := "SSN_DICTIONARY"
	dictID := 12345

	// Mock GetAll which is called by GetByName
	server.On("GET", "/zia/api/v1/dlpDictionaries", common.SuccessResponse([]dlpdictionaries.DlpDictionary{
		{ID: dictID, Name: dictName, Custom: false},
	}))

	// Mock the predefinedIdentifiers endpoint
	path := "/zia/api/v1/dlpDictionaries/12345/predefinedIdentifiers"
	server.On("GET", path, common.SuccessResponse([]string{"SSN", "PASSPORT", "DRIVER_LICENSE"}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, returnedID, err := dlpdictionaries.GetPredefinedIdentifiers(context.Background(), service, dictName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Equal(t, dictID, returnedID)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDLPDictionaries_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DlpDictionary JSON marshaling", func(t *testing.T) {
		dict := dlpdictionaries.DlpDictionary{
			ID:                    12345,
			Name:                  "Custom SSN Dictionary",
			Description:           "Detects SSN patterns",
			ConfidenceThreshold:   "CONFIDENCE_LEVEL_HIGH",
			CustomPhraseMatchType: "MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY",
			Custom:                true,
			ThresholdType:         "UNIQUE_COUNT",
			DictionaryType:        "PATTERNS_AND_PHRASES",
			Proximity:             50,
			Phrases: []dlpdictionaries.Phrases{
				{Action: "PHRASE_COUNT_TYPE_UNIQUE", Phrase: "social security"},
			},
			Patterns: []dlpdictionaries.Patterns{
				{Action: "PATTERN_COUNT_TYPE_UNIQUE", Pattern: "\\d{3}-\\d{2}-\\d{4}"},
			},
		}

		data, err := json.Marshal(dict)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Custom SSN Dictionary"`)
		assert.Contains(t, string(data), `"custom":true`)
	})

	t.Run("DlpDictionary JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Credit Card Dictionary",
			"description": "Detects credit card numbers",
			"confidenceThreshold": "CONFIDENCE_LEVEL_MEDIUM",
			"custom": false,
			"dictionaryType": "PREDEFINED",
			"proximityLengthEnabled": true,
			"includeBinNumbers": true,
			"binNumbers": [400000, 500000, 600000],
			"phrases": [
				{"action": "PHRASE_COUNT_TYPE_ALL", "phrase": "credit card"}
			],
			"patterns": [
				{"action": "PATTERN_COUNT_TYPE_ALL", "pattern": "\\d{4}-\\d{4}-\\d{4}-\\d{4}"}
			],
			"exactDataMatchDetails": [
				{
					"dictionaryEdmMappingId": 100,
					"schemaId": 200,
					"primaryFields": [1, 2],
					"secondaryFields": [3, 4]
				}
			]
		}`

		var dict dlpdictionaries.DlpDictionary
		err := json.Unmarshal([]byte(jsonData), &dict)
		require.NoError(t, err)

		assert.Equal(t, 54321, dict.ID)
		assert.False(t, dict.Custom)
		assert.True(t, dict.IncludeBinNumbers)
		assert.Len(t, dict.BinNumbers, 3)
		assert.Len(t, dict.EDMMatchDetails, 1)
	})
}

func TestDLPDictionaries_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse dictionaries list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "SSN", "custom": false, "dictionaryType": "PREDEFINED"},
			{"id": 2, "name": "Credit Card", "custom": false, "dictionaryType": "PREDEFINED"},
			{"id": 3, "name": "Custom PII", "custom": true, "dictionaryType": "PATTERNS_AND_PHRASES"}
		]`

		var dicts []dlpdictionaries.DlpDictionary
		err := json.Unmarshal([]byte(jsonResponse), &dicts)
		require.NoError(t, err)

		assert.Len(t, dicts, 3)
		assert.True(t, dicts[2].Custom)
	})
}
