// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

func TestDevicesFilters_Structure(t *testing.T) {
	t.Parallel()

	t.Run("GetDevicesFilters JSON marshaling", func(t *testing.T) {
		filters := devices.GetDevicesFilters{
			UserIDs: []int{1001, 1002, 1003},
			Emails:  []string{"user1@example.com", "user2@example.com"},
			Loc:     []int{100, 200},
			Dept:    []int{10, 20, 30},
			Geo:     []string{"US-CA", "US-NY", "GB-LDN"},
			Offset:  "page2",
			Limit:   50,
		}
		filters.From = 1699900000
		filters.To = 1700000000

		data, err := json.Marshal(filters)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"userids":[1001,1002,1003]`)
		assert.Contains(t, string(data), `"emails":["user1@example.com","user2@example.com"]`)
		assert.Contains(t, string(data), `"loc":[100,200]`)
		assert.Contains(t, string(data), `"from":1699900000`)
		assert.Contains(t, string(data), `"offset":"page2"`)
	})

	t.Run("GetDevicesFilters JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"from": 1699900000,
			"to": 1700000000,
			"userids": [5001, 5002],
			"emails": ["admin@example.com"],
			"loc": [1, 2, 3],
			"dept": [100],
			"geo": ["JP-TYO"],
			"offset": "next-page",
			"limit": 100
		}`

		var filters devices.GetDevicesFilters
		err := json.Unmarshal([]byte(jsonData), &filters)
		require.NoError(t, err)

		assert.Equal(t, 1699900000, filters.From)
		assert.Equal(t, 1700000000, filters.To)
		assert.Equal(t, []int{5001, 5002}, filters.UserIDs)
		assert.Equal(t, []string{"admin@example.com"}, filters.Emails)
		assert.Equal(t, []int{1, 2, 3}, filters.Loc)
		assert.Equal(t, 100, filters.Limit)
	})

	t.Run("GeoLocationFilter JSON marshaling", func(t *testing.T) {
		filter := devices.GeoLocationFilter{
			ParentGeoID: "US",
			Search:      "California",
		}
		filter.From = 1699900000
		filter.To = 1700000000

		data, err := json.Marshal(filter)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"parent_geo_id":"US"`)
		assert.Contains(t, string(data), `"search":"California"`)
		assert.Contains(t, string(data), `"from":1699900000`)
		assert.Contains(t, string(data), `"to":1700000000`)
	})

	t.Run("GeoLocationFilter JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"from": 1699900000,
			"to": 1700000000,
			"parent_geo_id": "US-CA",
			"search": "San"
		}`

		var filter devices.GeoLocationFilter
		err := json.Unmarshal([]byte(jsonData), &filter)
		require.NoError(t, err)

		assert.Equal(t, 1699900000, filter.From)
		assert.Equal(t, 1700000000, filter.To)
		assert.Equal(t, "US-CA", filter.ParentGeoID)
		assert.Equal(t, "San", filter.Search)
	})

	t.Run("Empty filters marshaling", func(t *testing.T) {
		filters := devices.GetDevicesFilters{}

		data, err := json.Marshal(filters)
		require.NoError(t, err)

		// Empty struct should marshal to empty object or object with only zero values
		assert.NotEmpty(t, data)
	})
}

