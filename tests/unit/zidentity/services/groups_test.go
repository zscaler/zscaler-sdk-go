// Package services provides unit tests for ZIdentity services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/groups"
)

func TestGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Groups JSON marshaling", func(t *testing.T) {
		group := groups.Groups{
			ID:                        "group-123",
			Name:                      "Engineering Team",
			Description:               "Engineering department group",
			Source:                    "OKTA",
			IsDynamicGroup:            false,
			DynamicGroup:              false,
			AdminEntitlementEnabled:   true,
			ServiceEntitlementEnabled: true,
			IDP: &common.IDNameDisplayName{
				ID:          "idp-001",
				Name:        "Okta",
				DisplayName: "Okta Identity Provider",
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"group-123"`)
		assert.Contains(t, string(data), `"name":"Engineering Team"`)
		assert.Contains(t, string(data), `"source":"OKTA"`)
		assert.Contains(t, string(data), `"adminEntitlementEnabled":true`)
	})

	t.Run("Groups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "group-456",
			"name": "Sales Team",
			"description": "Sales department group",
			"source": "AZURE_AD",
			"isDynamicGroup": true,
			"dynamicGroup": true,
			"adminEntitlementEnabled": false,
			"serviceEntitlementEnabled": true,
			"idp": {
				"id": "idp-002",
				"name": "Azure AD",
				"displayName": "Azure Active Directory"
			}
		}`

		var group groups.Groups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, "group-456", group.ID)
		assert.Equal(t, "Sales Team", group.Name)
		assert.Equal(t, "AZURE_AD", group.Source)
		assert.True(t, group.IsDynamicGroup)
		assert.True(t, group.DynamicGroup)
		assert.False(t, group.AdminEntitlementEnabled)
		assert.NotNil(t, group.IDP)
		assert.Equal(t, "Azure AD", group.IDP.Name)
	})

	t.Run("UserID JSON marshaling", func(t *testing.T) {
		userID := groups.UserID{
			ID: "user-12345",
		}

		data, err := json.Marshal(userID)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"user-12345"`)
	})

	t.Run("Groups without IDP", func(t *testing.T) {
		group := groups.Groups{
			ID:          "group-local",
			Name:        "Local Group",
			Description: "Locally managed group",
			Source:      "LOCAL",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"Local Group"`)
		assert.Contains(t, string(data), `"source":"LOCAL"`)
		assert.NotContains(t, string(data), `"idp"`)
	})
}

func TestGroups_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse groups list response", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 3,
			"pageOffset": 0,
			"pageSize": 100,
			"next_link": "",
			"records": [
				{
					"id": "group-001",
					"name": "Administrators",
					"source": "LOCAL",
					"adminEntitlementEnabled": true
				},
				{
					"id": "group-002",
					"name": "Developers",
					"source": "OKTA",
					"isDynamicGroup": false
				},
				{
					"id": "group-003",
					"name": "All Employees",
					"source": "AZURE_AD",
					"isDynamicGroup": true,
					"dynamicGroup": true
				}
			]
		}`

		var response common.PaginationResponse[groups.Groups]
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 3, response.ResultsTotal)
		assert.Len(t, response.Records, 3)
		assert.Equal(t, "Administrators", response.Records[0].Name)
		assert.Equal(t, "LOCAL", response.Records[0].Source)
		assert.True(t, response.Records[2].IsDynamicGroup)
	})

	t.Run("Parse single group response", func(t *testing.T) {
		jsonResponse := `{
			"id": "group-detailed",
			"name": "Security Team",
			"description": "Information Security Team",
			"source": "OKTA",
			"isDynamicGroup": false,
			"dynamicGroup": false,
			"adminEntitlementEnabled": true,
			"serviceEntitlementEnabled": true,
			"idp": {
				"id": "idp-okta",
				"name": "Okta",
				"displayName": "Corporate Okta"
			}
		}`

		var group groups.Groups
		err := json.Unmarshal([]byte(jsonResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, "group-detailed", group.ID)
		assert.Equal(t, "Security Team", group.Name)
		assert.Equal(t, "Information Security Team", group.Description)
		assert.True(t, group.AdminEntitlementEnabled)
		assert.NotNil(t, group.IDP)
		assert.Equal(t, "Corporate Okta", group.IDP.DisplayName)
	})

	t.Run("Parse paginated groups response", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 250,
			"pageOffset": 100,
			"pageSize": 100,
			"next_link": "/admin/api/v1/groups?offset=200&limit=100",
			"prev_link": "/admin/api/v1/groups?offset=0&limit=100",
			"records": [
				{"id": "group-101", "name": "Group 101"},
				{"id": "group-102", "name": "Group 102"}
			]
		}`

		var response common.PaginationResponse[groups.Groups]
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 250, response.ResultsTotal)
		assert.Equal(t, 100, response.PageOffset)
		assert.Equal(t, 100, response.PageSize)
		assert.NotEmpty(t, response.NextLink)
		assert.NotEmpty(t, response.PrevLink)
		assert.Len(t, response.Records, 2)
	})
}

