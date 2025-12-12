// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

func TestZIACommon_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IDNameExtensions JSON marshaling", func(t *testing.T) {
		idName := common.IDNameExtensions{
			ID:   12345,
			Name: "Test Resource",
			Extensions: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
		}

		data, err := json.Marshal(idName)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Test Resource"`)
		assert.Contains(t, string(data), `"extensions"`)
	})

	t.Run("IDName JSON marshaling", func(t *testing.T) {
		idName := common.IDName{
			ID:     12345,
			Name:   "Test Resource",
			Parent: "Parent Resource",
		}

		data, err := json.Marshal(idName)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"parent":"Parent Resource"`)
	})

	t.Run("IDNameExternalID JSON marshaling", func(t *testing.T) {
		idName := common.IDNameExternalID{
			ID:         12345,
			Name:       "Test Resource",
			ExternalID: "ext-12345",
		}

		data, err := json.Marshal(idName)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"externalId":"ext-12345"`)
	})

	t.Run("UserGroups JSON marshaling", func(t *testing.T) {
		group := common.UserGroups{
			ID:              12345,
			Name:            "Engineering",
			IdpID:           100,
			Comments:        "Engineering team",
			IsSystemDefined: "false",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"idp_id":100`)
	})

	t.Run("UserDepartment JSON marshaling", func(t *testing.T) {
		dept := common.UserDepartment{
			ID:       12345,
			Name:     "Engineering",
			IdpID:    100,
			Comments: "Engineering department",
			Deleted:  false,
		}

		data, err := json.Marshal(dept)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Engineering"`)
	})

	t.Run("DeviceGroups JSON marshaling", func(t *testing.T) {
		dg := common.DeviceGroups{
			ID:   12345,
			Name: "Mobile Devices",
		}

		data, err := json.Marshal(dg)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Mobile Devices"`)
	})

	t.Run("CommonNSS JSON marshaling", func(t *testing.T) {
		nss := common.CommonNSS{
			ID:          12345,
			PID:         100,
			Name:        "NSS Server",
			Description: "NSS server description",
			Deleted:     false,
			GetlID:      200,
		}

		data, err := json.Marshal(nss)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"pid":100`)
	})

	t.Run("ZPAAppSegments JSON marshaling", func(t *testing.T) {
		segment := common.ZPAAppSegments{
			ID:         12345,
			Name:       "ZPA Segment",
			ExternalID: "zpa-ext-123",
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"externalId":"zpa-ext-123"`)
	})

	t.Run("DatacenterSearchParameters structure", func(t *testing.T) {
		params := common.DatacenterSearchParameters{
			RoutableIP:                true,
			WithinCountryOnly:         true,
			IncludePrivateServiceEdge: false,
			IncludeCurrentVips:        true,
			SourceIp:                  "10.0.0.1",
			Latitude:                  37.7749,
			Longitude:                 -122.4194,
			Subcloud:                  "subcloud1",
		}

		assert.True(t, params.RoutableIP)
		assert.True(t, params.WithinCountryOnly)
		assert.Equal(t, "10.0.0.1", params.SourceIp)
	})

	t.Run("Order JSON marshaling", func(t *testing.T) {
		order := common.Order{
			On: "name",
			By: "asc",
		}

		data, err := json.Marshal(order)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"on":"name"`)
		assert.Contains(t, string(data), `"by":"asc"`)
	})

	t.Run("DataConsumed JSON marshaling", func(t *testing.T) {
		dc := common.DataConsumed{
			Min: 100,
			Max: 1000,
		}

		data, err := json.Marshal(dc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"min":100`)
		assert.Contains(t, string(data), `"max":1000`)
	})
}

func TestZIACommon_SortParams(t *testing.T) {
	t.Parallel()

	t.Run("GetSortParams with both parameters", func(t *testing.T) {
		params := common.GetSortParams(common.NameSortField, common.ASCSortOrder)
		assert.Contains(t, params, "sortBy=name")
		assert.Contains(t, params, "sortOrder=asc")
	})

	t.Run("GetSortParams with only sortBy", func(t *testing.T) {
		params := common.GetSortParams(common.IDSortField, "")
		assert.Equal(t, "sortBy=id", params)
	})

	t.Run("GetSortParams with only sortOrder", func(t *testing.T) {
		params := common.GetSortParams("", common.DESCSortOrder)
		assert.Equal(t, "sortOrder=desc", params)
	})

	t.Run("GetPageSize returns correct value", func(t *testing.T) {
		pageSize := common.GetPageSize()
		assert.Equal(t, 1000, pageSize)
	})
}

