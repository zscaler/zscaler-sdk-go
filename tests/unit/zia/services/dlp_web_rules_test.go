// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

func TestDLPWebRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebDLPRules JSON marshaling", func(t *testing.T) {
		rule := dlp_web_rules.WebDLPRules{
			ID:                       12345,
			Order:                    1,
			Rank:                     7,
			Name:                     "Block PII Upload",
			Description:              "Block uploads containing PII",
			Protocols:                []string{"HTTPS_RULE", "HTTP_RULE"},
			FileTypes:                []string{"ALL_DOCUMENT", "ALL_SPREADSHEET"},
			CloudApplications:        []string{"GOOGLE_DRIVE", "DROPBOX"},
			MinSize:                  1024,
			Action:                   "BLOCK",
			State:                    "ENABLED",
			OcrEnabled:               true,
			DLPDownloadScanEnabled:   true,
			ZCCNotificationsEnabled:  true,
			WithoutContentInspection: false,
			Severity:                 "RULE_SEVERITY_HIGH",
			UserRiskScoreLevels:      []string{"HIGH", "CRITICAL"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"action":"BLOCK"`)
		assert.Contains(t, string(data), `"ocrEnabled":true`)
	})

	t.Run("WebDLPRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"order": 2,
			"name": "Monitor Financial Data",
			"description": "Monitor financial data uploads",
			"protocols": ["HTTPS_RULE"],
			"fileTypes": ["ALL_SPREADSHEET"],
			"cloudApplications": ["OFFICE365"],
			"action": "ALLOW",
			"state": "ENABLED",
			"matchOnly": true,
			"ocrEnabled": false,
			"dlpDownloadScanEnabled": true,
			"externalAuditorEmail": "audit@company.com",
			"severity": "RULE_SEVERITY_MEDIUM",
			"locations": [
				{"id": 100, "name": "HQ"}
			],
			"departments": [
				{"id": 200, "name": "Finance"}
			],
			"dlpEngines": [
				{"id": 300, "name": "PCI DSS Engine"}
			],
			"workloadGroups": [
				{"id": 400, "name": "Production"}
			]
		}`

		var rule dlp_web_rules.WebDLPRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.True(t, rule.MatchOnly)
		assert.Equal(t, "audit@company.com", rule.ExternalAuditorEmail)
		assert.Len(t, rule.Locations, 1)
		assert.Len(t, rule.DLPEngines, 1)
	})

	t.Run("Receiver JSON marshaling", func(t *testing.T) {
		receiver := dlp_web_rules.Receiver{
			ID:   100,
			Name: "Incident Receiver",
			Type: "ZIA",
		}

		data, err := json.Marshal(receiver)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":100`)
		assert.Contains(t, string(data), `"type":"ZIA"`)
	})

	t.Run("SubRule JSON marshaling", func(t *testing.T) {
		subRule := dlp_web_rules.SubRule{
			ID: 12345,
		}

		data, err := json.Marshal(subRule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
	})
}

func TestDLPWebRules_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse DLP web rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Rule 1", "action": "BLOCK", "state": "ENABLED"},
			{"id": 2, "name": "Rule 2", "action": "ALLOW", "state": "ENABLED"},
			{"id": 3, "name": "Rule 3", "action": "CAUTION", "state": "DISABLED"}
		]`

		var rules []dlp_web_rules.WebDLPRules
		err := json.Unmarshal([]byte(jsonResponse), &rules)
		require.NoError(t, err)

		assert.Len(t, rules, 3)
		assert.Equal(t, "CAUTION", rules[2].Action)
	})

	t.Run("Parse rule with sub-rules", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Parent Rule",
			"action": "BLOCK",
			"subRules": [
				{"id": 101},
				{"id": 102}
			]
		}`

		var rule dlp_web_rules.WebDLPRules
		err := json.Unmarshal([]byte(jsonResponse), &rule)
		require.NoError(t, err)

		assert.Len(t, rule.SubRules, 2)
	})
}

