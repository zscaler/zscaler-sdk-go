// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/administration"
)

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

