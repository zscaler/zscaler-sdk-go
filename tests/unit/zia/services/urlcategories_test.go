// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlcategories"
)

func TestURLCategory_Structure(t *testing.T) {
	t.Parallel()

	t.Run("URLCategory JSON marshaling", func(t *testing.T) {
		cat := urlcategories.URLCategory{
			ID:                              "CUSTOM_01",
			ConfiguredName:                  "Blocked Sites",
			Keywords:                        []string{"gambling", "casino"},
			KeywordsRetainingParentCategory: []string{"adult"},
			Urls:                            []string{"*.badsite.com", "malware.example.com"},
			DBCategorizedUrls:               []string{"known-bad.com"},
			CustomCategory:                  true,
			Editable:                        true,
			Description:                     "Custom blocked sites",
			Type:                            "URL_CATEGORY",
			SuperCategory:                   "USER_DEFINED",
			IPRanges:                        []string{"192.168.1.0/24"},
		}

		data, err := json.Marshal(cat)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"CUSTOM_01"`)
		assert.Contains(t, string(data), `"configuredName":"Blocked Sites"`)
		assert.Contains(t, string(data), `"customCategory":true`)
	})

	t.Run("URLCategory JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "ADULT_CONTENT",
			"configuredName": "Adult Content",
			"keywords": ["adult", "xxx"],
			"keywordsRetainingParentCategory": [],
			"urls": [],
			"dbCategorizedUrls": ["adult-site.com"],
			"customCategory": false,
			"editable": false,
			"description": "Adult content websites",
			"type": "URL_CATEGORY",
			"superCategory": "ADULT_AND_MATURE",
			"urlKeywordCounts": {
				"totalUrlCount": 1000,
				"retainParentUrlCount": 50,
				"totalKeywordCount": 100,
				"retainParentKeywordCount": 10
			},
			"customUrlsCount": 0,
			"urlsRetainingParentCategoryCount": 0,
			"ipRanges": [],
			"ipRangesRetainingParentCategory": [],
			"customIpRangesCount": 0,
			"ipRangesRetainingParentCategoryCount": 0
		}`

		var cat urlcategories.URLCategory
		err := json.Unmarshal([]byte(jsonData), &cat)
		require.NoError(t, err)

		assert.Equal(t, "ADULT_CONTENT", cat.ID)
		assert.False(t, cat.CustomCategory)
		assert.NotNil(t, cat.URLKeywordCounts)
		assert.Equal(t, 1000, cat.URLKeywordCounts.TotalURLCount)
	})

	t.Run("Scopes JSON marshaling", func(t *testing.T) {
		scope := urlcategories.Scopes{
			Type: "LOCATION",
			ScopeEntities: []common.IDNameExtensions{
				{ID: 1, Name: "HQ"},
			},
		}

		data, err := json.Marshal(scope)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"Type":"LOCATION"`)
	})

	t.Run("URLKeywordCounts JSON marshaling", func(t *testing.T) {
		counts := urlcategories.URLKeywordCounts{
			TotalURLCount:            1000,
			RetainParentURLCount:     50,
			TotalKeywordCount:        100,
			RetainParentKeywordCount: 10,
		}

		data, err := json.Marshal(counts)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"totalUrlCount":1000`)
		assert.Contains(t, string(data), `"retainParentUrlCount":50`)
	})

	t.Run("URLQuota JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"uniqueUrlsProvisioned": 5000,
			"remainingUrlsQuota": 20000
		}`

		var quota urlcategories.URLQuota
		err := json.Unmarshal([]byte(jsonData), &quota)
		require.NoError(t, err)

		assert.Equal(t, 5000, quota.UniqueUrlsProvisioned)
		assert.Equal(t, 20000, quota.RemainingUrlsQuota)
	})

	t.Run("URLClassification JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"url": "example.com",
			"urlClassifications": ["BUSINESS", "TECHNOLOGY"],
			"urlClassificationsWithSecurityAlert": [],
			"application": "OFFICE_365"
		}`

		var classification urlcategories.URLClassification
		err := json.Unmarshal([]byte(jsonData), &classification)
		require.NoError(t, err)

		assert.Equal(t, "example.com", classification.URL)
		assert.Len(t, classification.URLClassifications, 2)
		assert.Equal(t, "OFFICE_365", classification.Application)
	})

	t.Run("URLReview JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"url": "newsite.com",
			"domainType": "FULL_DOMAIN",
			"matches": [
				{"id": "UNCATEGORIZED", "name": "Uncategorized"}
			]
		}`

		var review urlcategories.URLReview
		err := json.Unmarshal([]byte(jsonData), &review)
		require.NoError(t, err)

		assert.Equal(t, "newsite.com", review.URL)
		assert.Equal(t, "FULL_DOMAIN", review.DomainType)
		assert.Len(t, review.Matches, 1)
	})
}

func TestURLCategories_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse URL categories list", func(t *testing.T) {
		jsonResponse := `[
			{"id": "ADULT_CONTENT", "configuredName": "Adult Content", "customCategory": false},
			{"id": "GAMBLING", "configuredName": "Gambling", "customCategory": false},
			{"id": "CUSTOM_01", "configuredName": "My Custom", "customCategory": true}
		]`

		var cats []urlcategories.URLCategory
		err := json.Unmarshal([]byte(jsonResponse), &cats)
		require.NoError(t, err)

		assert.Len(t, cats, 3)
		assert.True(t, cats[2].CustomCategory)
	})

	t.Run("Parse URL lookup results", func(t *testing.T) {
		jsonResponse := `[
			{
				"url": "google.com",
				"urlClassifications": ["SEARCH_ENGINES"],
				"urlClassificationsWithSecurityAlert": []
			},
			{
				"url": "malware.com",
				"urlClassifications": ["MALWARE"],
				"urlClassificationsWithSecurityAlert": ["MALWARE"]
			}
		]`

		var classifications []urlcategories.URLClassification
		err := json.Unmarshal([]byte(jsonResponse), &classifications)
		require.NoError(t, err)

		assert.Len(t, classifications, 2)
		assert.Len(t, classifications[1].URLClassificationsWithSecurityAlert, 1)
	})
}

