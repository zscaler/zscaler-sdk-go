// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlpdictionaries"
)

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

	t.Run("Phrases JSON marshaling", func(t *testing.T) {
		phrase := dlpdictionaries.Phrases{
			Action: "PHRASE_COUNT_TYPE_UNIQUE",
			Phrase: "confidential",
		}

		data, err := json.Marshal(phrase)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"action":"PHRASE_COUNT_TYPE_UNIQUE"`)
		assert.Contains(t, string(data), `"phrase":"confidential"`)
	})

	t.Run("Patterns JSON marshaling", func(t *testing.T) {
		pattern := dlpdictionaries.Patterns{
			Action:  "PATTERN_COUNT_TYPE_ALL",
			Pattern: "\\b[A-Z]{2}\\d{6}\\b",
		}

		data, err := json.Marshal(pattern)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"action":"PATTERN_COUNT_TYPE_ALL"`)
		assert.Contains(t, string(data), `"pattern"`)
	})

	t.Run("EDMMatchDetails JSON marshaling", func(t *testing.T) {
		edm := dlpdictionaries.EDMMatchDetails{
			DictionaryEdmMappingID: 100,
			SchemaID:               200,
			PrimaryFields:          []int{1, 2, 3},
			SecondaryFields:        []int{4, 5},
			SecondaryFieldMatchOn:  "MATCHON_ANY",
		}

		data, err := json.Marshal(edm)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"dictionaryEdmMappingId":100`)
		assert.Contains(t, string(data), `"schemaId":200`)
	})

	t.Run("IDMProfileMatchAccuracy JSON marshaling", func(t *testing.T) {
		idm := dlpdictionaries.IDMProfileMatchAccuracy{
			MatchAccuracy: "LOW",
		}

		data, err := json.Marshal(idm)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"matchAccuracy":"LOW"`)
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

	t.Run("Parse hierarchical dictionary", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Hierarchical Dict",
			"hierarchicalDictionary": true,
			"hierarchicalIdentifiers": ["US_SSN", "US_DRIVER_LICENSE", "US_PASSPORT"]
		}`

		var dict dlpdictionaries.DlpDictionary
		err := json.Unmarshal([]byte(jsonResponse), &dict)
		require.NoError(t, err)

		assert.True(t, dict.HierarchicalDictionary)
		assert.Len(t, dict.HierarchicalIdentifiers, 3)
	})
}

