// Package services provides unit tests for ZDX users service
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/users"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestUsers_GetUser_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/users/12345"

	server.On("GET", path, common.SuccessResponse(users.User{
		ID:    12345,
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Devices: []users.Devices{
			{
				ID:   1001,
				Name: "LAPTOP-001",
				UserLocation: []users.UserLocation{
					{ID: "loc-1", City: "San Jose", State: "CA", Country: "US"},
				},
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := users.GetUser(context.Background(), service, "12345")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 12345, result.ID)
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "john.doe@example.com", result.Email)
	assert.Len(t, result.Devices, 1)
}

func TestUsers_GetAllUsers_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/users"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"users": []users.User{
			{ID: 1, Name: "User 1", Email: "user1@example.com"},
			{ID: 2, Name: "User 2", Email: "user2@example.com"},
			{ID: 3, Name: "User 3", Email: "user3@example.com"},
		},
		"next_offset": nil,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := users.GetAllUsers(context.Background(), service, users.GetUsersFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "User 1", result[0].Name)
}

func TestUsers_GetUser_WithDevices_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/users/99999"

	server.On("GET", path, common.SuccessResponse(users.User{
		ID:    99999,
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",
		Devices: []users.Devices{
			{
				ID:   2001,
				Name: "DESKTOP-001",
				UserLocation: []users.UserLocation{
					{ID: "geo-1", City: "New York", State: "NY", Country: "US", GeoLat: 40.7128, GeoLong: -74.0060},
				},
				ZSLocation: []users.ZSLocation{
					{ID: 100, Name: "NYC Data Center"},
				},
			},
			{
				ID:   2002,
				Name: "LAPTOP-002",
				UserLocation: []users.UserLocation{
					{ID: "geo-2", City: "Boston", State: "MA", Country: "US"},
				},
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := users.GetUser(context.Background(), service, "99999")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Devices, 2)
	assert.Equal(t, "DESKTOP-001", result.Devices[0].Name)
	assert.Equal(t, "New York", result.Devices[0].UserLocation[0].City)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("User JSON marshaling", func(t *testing.T) {
		user := users.User{
			ID:    12345,
			Name:  "Test User",
			Email: "test@example.com",
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Test User"`)
		assert.Contains(t, string(data), `"email":"test@example.com"`)
	})

	t.Run("User JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "Another User",
			"email": "another@example.com"
		}`

		var user users.User
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, 67890, user.ID)
		assert.Equal(t, "Another User", user.Name)
		assert.Equal(t, "another@example.com", user.Email)
	})

	t.Run("Devices JSON marshaling", func(t *testing.T) {
		device := users.Devices{
			ID:   1001,
			Name: "WORKSTATION-001",
			UserLocation: []users.UserLocation{
				{ID: "loc-1", City: "San Jose", State: "CA", Country: "US"},
			},
			ZSLocation: []users.ZSLocation{
				{ID: 100, Name: "West Coast DC"},
			},
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":1001`)
		assert.Contains(t, string(data), `"name":"WORKSTATION-001"`)
		assert.Contains(t, string(data), `"geo_loc"`)
		assert.Contains(t, string(data), `"zs_loc"`)
	})

	t.Run("UserLocation JSON marshaling", func(t *testing.T) {
		location := users.UserLocation{
			ID:           "geo-12345",
			City:         "San Francisco",
			State:        "California",
			Country:      "United States",
			GeoLat:       37.7749,
			GeoLong:      -122.4194,
			GeoDetection: "IP",
		}

		data, err := json.Marshal(location)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"geo-12345"`)
		assert.Contains(t, string(data), `"city":"San Francisco"`)
		assert.Contains(t, string(data), `"country":"United States"`)
		assert.Contains(t, string(data), `"geo_lat":37.7749`)
	})

	t.Run("ZSLocation JSON marshaling", func(t *testing.T) {
		zsLoc := users.ZSLocation{
			ID:   500,
			Name: "Zscaler West Coast",
		}

		data, err := json.Marshal(zsLoc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":500`)
		assert.Contains(t, string(data), `"name":"Zscaler West Coast"`)
	})

	t.Run("Complete User with devices JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 11111,
			"name": "Complex User",
			"email": "complex@example.com",
			"devices": [
				{
					"id": 2001,
					"name": "LAPTOP-XYZ",
					"geo_loc": [
						{"id": "geo-1", "city": "Austin", "state": "TX", "country": "US", "geo_lat": 30.2672, "geo_long": -97.7431}
					],
					"zs_loc": [
						{"id": 200, "name": "Texas DC"}
					]
				}
			]
		}`

		var user users.User
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, 11111, user.ID)
		assert.Len(t, user.Devices, 1)
		assert.Equal(t, "LAPTOP-XYZ", user.Devices[0].Name)
		assert.Equal(t, "Austin", user.Devices[0].UserLocation[0].City)
		assert.Equal(t, "Texas DC", user.Devices[0].ZSLocation[0].Name)
	})
}

func TestUsers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse users list response", func(t *testing.T) {
		jsonResponse := `{
			"users": [
				{"id": 1, "name": "User 1", "email": "user1@test.com"},
				{"id": 2, "name": "User 2", "email": "user2@test.com"},
				{"id": 3, "name": "User 3", "email": "user3@test.com"}
			],
			"next_offset": "page2"
		}`

		var response struct {
			Users      []users.User `json:"users"`
			NextOffset string       `json:"next_offset"`
		}
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Users, 3)
		assert.Equal(t, "page2", response.NextOffset)
	})

	t.Run("Parse empty users list", func(t *testing.T) {
		jsonResponse := `{
			"users": [],
			"next_offset": null
		}`

		var response struct {
			Users      []users.User `json:"users"`
			NextOffset interface{}  `json:"next_offset"`
		}
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Empty(t, response.Users)
	})
}
