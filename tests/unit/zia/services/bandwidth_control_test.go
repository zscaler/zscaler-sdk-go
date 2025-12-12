// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_classes"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_control_rules"
)

func TestBandwidthClasses_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BandwidthClasses JSON marshaling", func(t *testing.T) {
		class := bandwidth_classes.BandwidthClasses{
			ID:                       12345,
			Name:                     "Streaming Media Class",
			Type:                     "CUSTOM",
			IsNameL10nTag:            false,
			WebApplications:          []string{"YOUTUBE", "NETFLIX"},
			UrlCategories:            []string{"STREAMING_MEDIA"},
			NetworkApplications:      []string{"HTTP", "HTTPS"},
			NetworkServices:          []string{"TCP_443", "TCP_80"},
			ApplicationServiceGroups: []string{"VIDEO_STREAMING"},
		}

		data, err := json.Marshal(class)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"CUSTOM"`)
	})

	t.Run("BandwidthClasses JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Cloud Storage Class",
			"type": "PREDEFINED",
			"isNameL10nTag": true,
			"webApplications": ["GOOGLE_DRIVE", "DROPBOX", "BOX"],
			"urlCategories": ["CLOUD_STORAGE"],
			"urls": ["*.googleapis.com", "*.dropbox.com"],
			"applications": ["GOOGLE_DRIVE", "DROPBOX"]
		}`

		var class bandwidth_classes.BandwidthClasses
		err := json.Unmarshal([]byte(jsonData), &class)
		require.NoError(t, err)

		assert.Equal(t, 54321, class.ID)
		assert.True(t, class.IsNameL10nTag)
		assert.Len(t, class.WebApplications, 3)
	})
}

func TestBandwidthControlRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BandwidthControlRules JSON marshaling", func(t *testing.T) {
		rule := bandwidth_control_rules.BandwidthControlRules{
			ID:                12345,
			Name:              "Limit Streaming",
			Order:             1,
			State:             "ENABLED",
			Description:       "Limit streaming media bandwidth",
			MaxBandwidth:      5000000, // 5 Mbps
			MinBandwidth:      1000000, // 1 Mbps
			Rank:              7,
			AccessControl:     "READ_WRITE",
			DefaultRule:       false,
			Protocols:         []string{"HTTPS_RULE", "HTTP_RULE"},
			DeviceTrustLevels: []string{"HIGH_TRUST"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"maxBandwidth":5000000`)
		assert.Contains(t, string(data), `"minBandwidth":1000000`)
	})

	t.Run("BandwidthControlRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Priority Traffic",
			"order": 1,
			"state": "ENABLED",
			"description": "Prioritize business traffic",
			"maxBandwidth": 0,
			"minBandwidth": 10000000,
			"rank": 5,
			"accessControl": "READ_ONLY",
			"defaultRule": false,
			"protocols": ["HTTPS_RULE"],
			"lastModifiedTime": 1699000000,
			"lastModifiedBy": {
				"id": 100,
				"name": "admin@company.com"
			},
			"bandwidthClasses": [
				{"id": 200, "name": "Business Apps"}
			],
			"locations": [
				{"id": 300, "name": "HQ"}
			],
			"locationGroups": [
				{"id": 400, "name": "North America"}
			]
		}`

		var rule bandwidth_control_rules.BandwidthControlRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, 10000000, rule.MinBandwidth)
		assert.NotNil(t, rule.LastModifiedBy)
		assert.Len(t, rule.BandwidthClasses, 1)
	})

	t.Run("BandwidthControlRules default rule", func(t *testing.T) {
		jsonData := `{
			"id": 1,
			"name": "Default Rule",
			"order": 100,
			"state": "ENABLED",
			"defaultRule": true,
			"maxBandwidth": 100000000
		}`

		var rule bandwidth_control_rules.BandwidthControlRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.True(t, rule.DefaultRule)
	})
}

func TestBandwidthControl_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse bandwidth classes list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Video Streaming", "type": "PREDEFINED"},
			{"id": 2, "name": "Cloud Storage", "type": "PREDEFINED"},
			{"id": 3, "name": "Custom Class", "type": "CUSTOM"}
		]`

		var classes []bandwidth_classes.BandwidthClasses
		err := json.Unmarshal([]byte(jsonResponse), &classes)
		require.NoError(t, err)

		assert.Len(t, classes, 3)
		assert.Equal(t, "CUSTOM", classes[2].Type)
	})

	t.Run("Parse bandwidth control rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Rule 1", "state": "ENABLED", "maxBandwidth": 10000000},
			{"id": 2, "name": "Rule 2", "state": "ENABLED", "maxBandwidth": 5000000},
			{"id": 3, "name": "Default", "state": "ENABLED", "defaultRule": true}
		]`

		var rules []bandwidth_control_rules.BandwidthControlRules
		err := json.Unmarshal([]byte(jsonResponse), &rules)
		require.NoError(t, err)

		assert.Len(t, rules, 3)
		assert.True(t, rules[2].DefaultRule)
	})
}

