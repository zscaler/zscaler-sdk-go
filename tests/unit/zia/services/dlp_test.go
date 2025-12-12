// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_engines"
)

func TestDLPEngines_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DLPEngines JSON marshaling", func(t *testing.T) {
		engine := dlp_engines.DLPEngines{
			ID:                   12345,
			Name:                 "Custom DLP Engine",
			Description:          "Custom engine for PII detection",
			PredefinedEngineName: "",
			EngineExpression:     "((D63.S > 1))",
			CustomDlpEngine:      true,
		}

		data, err := json.Marshal(engine)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Custom DLP Engine"`)
		assert.Contains(t, string(data), `"customDlpEngine":true`)
	})

	t.Run("DLPEngines JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "PCI DSS Engine",
			"description": "Detects payment card data",
			"predefinedEngineName": "PCI_DSS",
			"engineExpression": "((D1.S > 0) AND (D2.S > 0))",
			"customDlpEngine": false
		}`

		var engine dlp_engines.DLPEngines
		err := json.Unmarshal([]byte(jsonData), &engine)
		require.NoError(t, err)

		assert.Equal(t, 54321, engine.ID)
		assert.Equal(t, "PCI_DSS", engine.PredefinedEngineName)
		assert.False(t, engine.CustomDlpEngine)
	})
}

func TestDLP_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse DLP engines list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "HIPAA Engine", "predefinedEngineName": "HIPAA", "customDlpEngine": false},
			{"id": 2, "name": "PCI DSS Engine", "predefinedEngineName": "PCI_DSS", "customDlpEngine": false},
			{"id": 3, "name": "Custom SSN", "customDlpEngine": true}
		]`

		var engines []dlp_engines.DLPEngines
		err := json.Unmarshal([]byte(jsonResponse), &engines)
		require.NoError(t, err)

		assert.Len(t, engines, 3)
		assert.True(t, engines[2].CustomDlpEngine)
	})

	t.Run("Parse DLP engine lite", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Engine 1", "predefinedEngineName": "BUILT_IN_1"},
			{"id": 2, "name": "Engine 2", "predefinedEngineName": "BUILT_IN_2"}
		]`

		var engines []dlp_engines.DLPEngines
		err := json.Unmarshal([]byte(jsonResponse), &engines)
		require.NoError(t, err)

		assert.Len(t, engines, 2)
		assert.Equal(t, "BUILT_IN_1", engines[0].PredefinedEngineName)
	})
}

