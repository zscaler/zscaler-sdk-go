// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/users"
)

func TestUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("User JSON marshaling", func(t *testing.T) {
		user := users.User{
			ID:    12345,
			Name:  "John Doe",
			Email: "john.doe@example.com",
			Devices: []users.Devices{
				{
					ID:   1001,
					Name: "LAPTOP-001",
					UserLocation: []users.UserLocation{
						{
							ID:           "US-CA-SJC",
							City:         "San Jose",
							State:        "California",
							Country:      "United States",
							GeoLat:       37.3382,
							GeoLong:      -121.8863,
							GeoDetection: "IP",
						},
					},
					ZSLocation: []users.ZSLocation{
						{ID: 100, Name: "San Jose DC"},
					},
				},
			},
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"John Doe"`)
		assert.Contains(t, string(data), `"email":"john.doe@example.com"`)
	})

	t.Run("User JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "Jane Smith",
			"email": "jane.smith@example.com",
			"devices": [
				{
					"id": 2001,
					"name": "DESKTOP-001",
					"geo_loc": [
						{
							"id": "US-NY-NYC",
							"city": "New York",
							"state": "New York",
							"country": "United States",
							"geo_lat": 40.7128,
							"geo_long": -74.0060,
							"geo_detection": "GPS"
						}
					],
					"zs_loc": [
						{"id": 200, "name": "New York DC"}
					]
				}
			]
		}`

		var user users.User
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, 67890, user.ID)
		assert.Equal(t, "Jane Smith", user.Name)
		assert.Equal(t, "jane.smith@example.com", user.Email)
		assert.Len(t, user.Devices, 1)
		assert.Equal(t, "DESKTOP-001", user.Devices[0].Name)
		assert.Len(t, user.Devices[0].UserLocation, 1)
		assert.Equal(t, "New York", user.Devices[0].UserLocation[0].City)
	})

	t.Run("Devices JSON marshaling", func(t *testing.T) {
		device := users.Devices{
			ID:   3001,
			Name: "MACBOOK-001",
			UserLocation: []users.UserLocation{
				{ID: "GB-LDN", City: "London", Country: "United Kingdom"},
			},
			ZSLocation: []users.ZSLocation{
				{ID: 300, Name: "London DC"},
			},
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":3001`)
		assert.Contains(t, string(data), `"name":"MACBOOK-001"`)
	})

	t.Run("UserLocation JSON marshaling", func(t *testing.T) {
		location := users.UserLocation{
			ID:           "JP-TYO",
			City:         "Tokyo",
			State:        "Tokyo",
			Country:      "Japan",
			GeoLat:       35.6762,
			GeoLong:      139.6503,
			GeoDetection: "IP",
		}

		data, err := json.Marshal(location)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"JP-TYO"`)
		assert.Contains(t, string(data), `"city":"Tokyo"`)
		assert.Contains(t, string(data), `"geo_lat":35.6762`)
	})

	t.Run("ZSLocation JSON marshaling", func(t *testing.T) {
		zsLoc := users.ZSLocation{
			ID:   500,
			Name: "Frankfurt DC",
		}

		data, err := json.Marshal(zsLoc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":500`)
		assert.Contains(t, string(data), `"name":"Frankfurt DC"`)
	})

	t.Run("GetUsersFilters JSON marshaling", func(t *testing.T) {
		filters := users.GetUsersFilters{
			From:   1699900000,
			To:     1700000000,
			Loc:    []int{1, 2},
			Dept:   []int{10, 20},
			Geo:    []string{"US-CA", "US-NY"},
			Offset: "page2",
			Limit:  50,
			Q:      "john",
		}

		data, err := json.Marshal(filters)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"from":1699900000`)
		assert.Contains(t, string(data), `"q":"john"`)
	})
}

func TestUsers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse users list response", func(t *testing.T) {
		jsonResponse := `{
			"users": [
				{
					"id": 1,
					"name": "User One",
					"email": "user1@example.com",
					"devices": [
						{"id": 101, "name": "Device 1"}
					]
				},
				{
					"id": 2,
					"name": "User Two",
					"email": "user2@example.com",
					"devices": [
						{"id": 201, "name": "Device 2"},
						{"id": 202, "name": "Device 3"}
					]
				}
			],
			"next_offset": "abc123"
		}`

		var response struct {
			Users      []users.User `json:"users"`
			NextOffset string       `json:"next_offset"`
		}
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Users, 2)
		assert.Equal(t, "abc123", response.NextOffset)
		assert.Equal(t, "User One", response.Users[0].Name)
		assert.Len(t, response.Users[0].Devices, 1)
		assert.Len(t, response.Users[1].Devices, 2)
	})

	t.Run("Parse single user response", func(t *testing.T) {
		jsonResponse := `{
			"id": 12345,
			"name": "Test User",
			"email": "test@example.com",
			"devices": [
				{
					"id": 1001,
					"name": "Work Laptop",
					"geo_loc": [
						{
							"id": "US-CA",
							"city": "San Francisco",
							"state": "California",
							"country": "United States"
						}
					],
					"zs_loc": [
						{"id": 1, "name": "US West"}
					]
				}
			]
		}`

		var user users.User
		err := json.Unmarshal([]byte(jsonResponse), &user)
		require.NoError(t, err)

		assert.Equal(t, 12345, user.ID)
		assert.Equal(t, "Test User", user.Name)
		assert.Len(t, user.Devices, 1)
		assert.Equal(t, "Work Laptop", user.Devices[0].Name)
		assert.Len(t, user.Devices[0].UserLocation, 1)
		assert.Equal(t, "San Francisco", user.Devices[0].UserLocation[0].City)
	})
}

