// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
)

func TestApplications_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Apps JSON marshaling", func(t *testing.T) {
		app := applications.Apps{
			ID:         12345,
			Name:       "Microsoft 365",
			Score:      85.5,
			TotalUsers: 1500,
			MostImpactedRegion: &applications.MostImpactedRegion{
				ID:      "US-CA",
				City:    "San Jose",
				Region:  "California",
				Country: "United States",
				GeoType: "city",
			},
			Stats: &applications.Stats{
				ActiveUsers:   1200,
				ActiveDevices: 1400,
				NumPoor:       50,
				NumOkay:       200,
				NumGood:       950,
			},
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Microsoft 365"`)
		assert.Contains(t, string(data), `"score":85.5`)
		assert.Contains(t, string(data), `"total_users":1500`)
	})

	t.Run("Apps JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "Salesforce",
			"score": 92.3,
			"total_users": 800,
			"most_impacted_region": {
				"id": "US-NY",
				"city": "New York",
				"region": "New York",
				"country": "United States",
				"geo_type": "city"
			},
			"stats": {
				"active_users": 750,
				"active_devices": 800,
				"num_poor": 10,
				"num_okay": 50,
				"num_good": 690
			}
		}`

		var app applications.Apps
		err := json.Unmarshal([]byte(jsonData), &app)
		require.NoError(t, err)

		assert.Equal(t, 67890, app.ID)
		assert.Equal(t, "Salesforce", app.Name)
		assert.Equal(t, float32(92.3), app.Score)
		assert.Equal(t, 800, app.TotalUsers)
		assert.NotNil(t, app.MostImpactedRegion)
		assert.Equal(t, "New York", app.MostImpactedRegion.City)
		assert.NotNil(t, app.Stats)
		assert.Equal(t, 750, app.Stats.ActiveUsers)
	})

	t.Run("MostImpactedRegion JSON marshaling", func(t *testing.T) {
		region := applications.MostImpactedRegion{
			ID:      "GB-LDN",
			City:    "London",
			Region:  "England",
			Country: "United Kingdom",
			GeoType: "city",
		}

		data, err := json.Marshal(region)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"GB-LDN"`)
		assert.Contains(t, string(data), `"city":"London"`)
		assert.Contains(t, string(data), `"country":"United Kingdom"`)
	})

	t.Run("Stats JSON marshaling", func(t *testing.T) {
		stats := applications.Stats{
			ActiveUsers:   5000,
			ActiveDevices: 6000,
			NumPoor:       100,
			NumOkay:       500,
			NumGood:       4400,
		}

		data, err := json.Marshal(stats)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"active_users":5000`)
		assert.Contains(t, string(data), `"active_devices":6000`)
		assert.Contains(t, string(data), `"num_poor":100`)
		assert.Contains(t, string(data), `"num_good":4400`)
	})
}

func TestApplications_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse applications list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "Microsoft 365",
				"score": 88.5,
				"total_users": 2000,
				"most_impacted_region": {
					"id": "US-CA",
					"city": "San Francisco"
				},
				"stats": {
					"active_users": 1800,
					"active_devices": 2100,
					"num_poor": 50,
					"num_okay": 150,
					"num_good": 1600
				}
			},
			{
				"id": 2,
				"name": "Slack",
				"score": 95.0,
				"total_users": 1500,
				"stats": {
					"active_users": 1400,
					"active_devices": 1600,
					"num_poor": 10,
					"num_okay": 40,
					"num_good": 1350
				}
			},
			{
				"id": 3,
				"name": "Zoom",
				"score": 78.2,
				"total_users": 1800,
				"most_impacted_region": {
					"id": "DE-BER",
					"city": "Berlin",
					"country": "Germany"
				}
			}
		]`

		var apps []applications.Apps
		err := json.Unmarshal([]byte(jsonResponse), &apps)
		require.NoError(t, err)

		assert.Len(t, apps, 3)
		
		// Check first app
		assert.Equal(t, "Microsoft 365", apps[0].Name)
		assert.Equal(t, float32(88.5), apps[0].Score)
		assert.NotNil(t, apps[0].Stats)
		assert.Equal(t, 1800, apps[0].Stats.ActiveUsers)
		
		// Check second app
		assert.Equal(t, "Slack", apps[1].Name)
		assert.Equal(t, float32(95.0), apps[1].Score)
		assert.Nil(t, apps[1].MostImpactedRegion)
		
		// Check third app
		assert.Equal(t, "Zoom", apps[2].Name)
		assert.NotNil(t, apps[2].MostImpactedRegion)
		assert.Equal(t, "Berlin", apps[2].MostImpactedRegion.City)
	})

	t.Run("Parse single application response", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "SAP",
			"score": 72.5,
			"total_users": 500,
			"most_impacted_region": {
				"id": "JP-TYO",
				"city": "Tokyo",
				"region": "Kanto",
				"country": "Japan",
				"geo_type": "city"
			},
			"stats": {
				"active_users": 480,
				"active_devices": 520,
				"num_poor": 80,
				"num_okay": 120,
				"num_good": 280
			}
		}`

		var app applications.Apps
		err := json.Unmarshal([]byte(jsonResponse), &app)
		require.NoError(t, err)

		assert.Equal(t, 100, app.ID)
		assert.Equal(t, "SAP", app.Name)
		assert.Equal(t, float32(72.5), app.Score)
		assert.Equal(t, "Tokyo", app.MostImpactedRegion.City)
		assert.Equal(t, "Japan", app.MostImpactedRegion.Country)
		assert.Equal(t, 80, app.Stats.NumPoor)
	})
}

