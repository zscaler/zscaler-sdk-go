// Package services provides unit tests for ZDX geo locations service
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestGeoLocations_GetGeoLocations_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/active_geo"

	server.On("GET", path, common.SuccessResponse([]devices.GeoLocation{
		{
			ID:      "US",
			Name:    "United States",
			GeoType: "country",
			Children: []devices.Children{
				{ID: "US-CA", Description: "California", GeoType: "state"},
				{ID: "US-NY", Description: "New York", GeoType: "state"},
			},
		},
		{
			ID:      "GB",
			Name:    "United Kingdom",
			GeoType: "country",
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := devices.GetGeoLocations(context.Background(), service, devices.GeoLocationFilter{})

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "United States", result[0].Name)
	assert.Len(t, result[0].Children, 2)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestGeoLocations_Structure(t *testing.T) {
	t.Parallel()

	t.Run("GeoLocation JSON marshaling", func(t *testing.T) {
		geoLoc := devices.GeoLocation{
			ID:          "US",
			Name:        "United States",
			GeoType:     "country",
			Description: "United States of America",
			Children: []devices.Children{
				{ID: "US-CA", Description: "California", GeoType: "state"},
				{ID: "US-NY", Description: "New York", GeoType: "state"},
				{ID: "US-TX", Description: "Texas", GeoType: "state"},
			},
		}

		data, err := json.Marshal(geoLoc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"US"`)
		assert.Contains(t, string(data), `"name":"United States"`)
		assert.Contains(t, string(data), `"geo_type":"country"`)
		assert.Contains(t, string(data), `"children"`)
	})

	t.Run("GeoLocation JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "US-CA",
			"name": "California",
			"geo_type": "state",
			"description": "State of California",
			"children": [
				{"id": "US-CA-SFO", "description": "San Francisco", "geo_type": "city"},
				{"id": "US-CA-LAX", "description": "Los Angeles", "geo_type": "city"},
				{"id": "US-CA-SJC", "description": "San Jose", "geo_type": "city"}
			]
		}`

		var geoLoc devices.GeoLocation
		err := json.Unmarshal([]byte(jsonData), &geoLoc)
		require.NoError(t, err)

		assert.Equal(t, "US-CA", geoLoc.ID)
		assert.Equal(t, "California", geoLoc.Name)
		assert.Equal(t, "state", geoLoc.GeoType)
		assert.Len(t, geoLoc.Children, 3)
		assert.Equal(t, "US-CA-SFO", geoLoc.Children[0].ID)
	})

	t.Run("Children JSON marshaling", func(t *testing.T) {
		child := devices.Children{
			ID:          "US-CA-SJC",
			Description: "San Jose",
			GeoType:     "city",
		}

		data, err := json.Marshal(child)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"US-CA-SJC"`)
		assert.Contains(t, string(data), `"description":"San Jose"`)
		assert.Contains(t, string(data), `"geo_type":"city"`)
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
	})
}

func TestGeoLocations_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse geo locations list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": "US",
				"name": "United States",
				"geo_type": "country",
				"children": [
					{"id": "US-CA", "description": "California", "geo_type": "state"},
					{"id": "US-NY", "description": "New York", "geo_type": "state"}
				]
			},
			{
				"id": "GB",
				"name": "United Kingdom",
				"geo_type": "country",
				"children": [
					{"id": "GB-ENG", "description": "England", "geo_type": "region"}
				]
			}
		]`

		var geoLocs []devices.GeoLocation
		err := json.Unmarshal([]byte(jsonResponse), &geoLocs)
		require.NoError(t, err)

		assert.Len(t, geoLocs, 2)
		assert.Equal(t, "United States", geoLocs[0].Name)
		assert.Len(t, geoLocs[0].Children, 2)
		assert.Equal(t, "United Kingdom", geoLocs[1].Name)
		assert.Len(t, geoLocs[1].Children, 1)
	})

	t.Run("Parse state level geo locations", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": "US-CA",
				"name": "California",
				"geo_type": "state",
				"children": [
					{"id": "US-CA-SFO", "description": "San Francisco", "geo_type": "city"},
					{"id": "US-CA-LAX", "description": "Los Angeles", "geo_type": "city"},
					{"id": "US-CA-SJC", "description": "San Jose", "geo_type": "city"},
					{"id": "US-CA-SDG", "description": "San Diego", "geo_type": "city"}
				]
			}
		]`

		var geoLocs []devices.GeoLocation
		err := json.Unmarshal([]byte(jsonResponse), &geoLocs)
		require.NoError(t, err)

		assert.Len(t, geoLocs, 1)
		assert.Equal(t, "California", geoLocs[0].Name)
		assert.Equal(t, "state", geoLocs[0].GeoType)
		assert.Len(t, geoLocs[0].Children, 4)
	})
}

