// Package services provides unit tests for ZDX services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/administration"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestAdministration_GetDepartments_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/administration/departments"

	server.On("GET", path, common.SuccessResponse([]administration.Department{
		{ID: 1, Name: "Engineering"},
		{ID: 2, Name: "Sales"},
		{ID: 3, Name: "Marketing"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := administration.GetDepartments(context.Background(), service, administration.GetDepartmentsFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Engineering", result[0].Name)
}

func TestAdministration_GetDepartments_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/administration/departments"

	server.On("GET", path, common.SuccessResponse([]administration.Department{
		{ID: 1, Name: "Engineering"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	filters := administration.GetDepartmentsFilters{
		From:   1699900000,
		To:     1700000000,
		Search: "Engineering",
	}

	result, _, err := administration.GetDepartments(context.Background(), service, filters)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Engineering", result[0].Name)
}

func TestAdministration_GetLocations_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/administration/locations"

	server.On("GET", path, common.SuccessResponse([]administration.Location{
		{ID: 100, Name: "San Jose HQ"},
		{ID: 101, Name: "New York Office"},
		{ID: 102, Name: "London Office"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := administration.GetLocations(context.Background(), service, administration.GetLocationsFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "San Jose HQ", result[0].Name)
	assert.Equal(t, 100, result[0].ID)
}

func TestAdministration_GetLocations_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/administration/locations"

	server.On("GET", path, common.SuccessResponse([]administration.Location{
		{ID: 100, Name: "San Jose HQ"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	filters := administration.GetLocationsFilters{
		From:   1699900000,
		To:     1700000000,
		Search: "San Jose",
		Q:      "HQ",
	}

	result, _, err := administration.GetLocations(context.Background(), service, filters)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "San Jose HQ", result[0].Name)
}

func TestAdministration_GetDepartments_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/administration/departments"

	server.On("GET", path, common.SuccessResponse([]administration.Department{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := administration.GetDepartments(context.Background(), service, administration.GetDepartmentsFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 0)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestAdministration_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Department JSON marshaling", func(t *testing.T) {
		dept := administration.Department{
			ID:   12345,
			Name: "Engineering",
		}

		data, err := json.Marshal(dept)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Engineering"`)
	})

	t.Run("Department JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "Sales"
		}`

		var dept administration.Department
		err := json.Unmarshal([]byte(jsonData), &dept)
		require.NoError(t, err)

		assert.Equal(t, 67890, dept.ID)
		assert.Equal(t, "Sales", dept.Name)
	})

	t.Run("Location JSON marshaling", func(t *testing.T) {
		location := administration.Location{
			ID:   11111,
			Name: "San Jose Headquarters",
		}

		data, err := json.Marshal(location)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":11111`)
		assert.Contains(t, string(data), `"name":"San Jose Headquarters"`)
	})

	t.Run("Location JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 22222,
			"name": "New York Office"
		}`

		var location administration.Location
		err := json.Unmarshal([]byte(jsonData), &location)
		require.NoError(t, err)

		assert.Equal(t, 22222, location.ID)
		assert.Equal(t, "New York Office", location.Name)
	})

	t.Run("GetDepartmentsFilters JSON marshaling", func(t *testing.T) {
		filters := administration.GetDepartmentsFilters{
			From:   1699900000,
			To:     1700000000,
			Search: "Engineering",
		}

		data, err := json.Marshal(filters)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"from":1699900000`)
		assert.Contains(t, string(data), `"to":1700000000`)
		assert.Contains(t, string(data), `"search":"Engineering"`)
	})

	t.Run("GetLocationsFilters JSON marshaling", func(t *testing.T) {
		filters := administration.GetLocationsFilters{
			From:   1699900000,
			To:     1700000000,
			Search: "California",
			Q:      "HQ",
		}

		data, err := json.Marshal(filters)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"from":1699900000`)
		assert.Contains(t, string(data), `"search":"California"`)
		assert.Contains(t, string(data), `"q":"HQ"`)
	})
}

func TestAdministration_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse departments list response", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Engineering"},
			{"id": 2, "name": "Sales"},
			{"id": 3, "name": "Marketing"},
			{"id": 4, "name": "Human Resources"},
			{"id": 5, "name": "Finance"}
		]`

		var departments []administration.Department
		err := json.Unmarshal([]byte(jsonResponse), &departments)
		require.NoError(t, err)

		assert.Len(t, departments, 5)
		assert.Equal(t, "Engineering", departments[0].Name)
		assert.Equal(t, "Sales", departments[1].Name)
		assert.Equal(t, "Finance", departments[4].Name)
	})

	t.Run("Parse locations list response", func(t *testing.T) {
		jsonResponse := `[
			{"id": 100, "name": "San Jose HQ"},
			{"id": 101, "name": "New York Office"},
			{"id": 102, "name": "London Office"},
			{"id": 103, "name": "Tokyo Office"}
		]`

		var locations []administration.Location
		err := json.Unmarshal([]byte(jsonResponse), &locations)
		require.NoError(t, err)

		assert.Len(t, locations, 4)
		assert.Equal(t, "San Jose HQ", locations[0].Name)
		assert.Equal(t, 100, locations[0].ID)
		assert.Equal(t, "Tokyo Office", locations[3].Name)
	})

	t.Run("Parse empty departments response", func(t *testing.T) {
		jsonResponse := `[]`

		var departments []administration.Department
		err := json.Unmarshal([]byte(jsonResponse), &departments)
		require.NoError(t, err)

		assert.Empty(t, departments)
	})
}
