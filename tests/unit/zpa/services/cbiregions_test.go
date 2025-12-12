// Package unit provides unit tests for ZPA CBI Regions service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CBIRegions represents the CBI region for testing
type CBIRegions struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TestCBIRegions_Structure tests the struct definitions
func TestCBIRegions_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CBIRegions JSON marshaling", func(t *testing.T) {
		region := CBIRegions{
			ID:   "us-west-1",
			Name: "US West (Oregon)",
		}

		data, err := json.Marshal(region)
		require.NoError(t, err)

		var unmarshaled CBIRegions
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, region.ID, unmarshaled.ID)
		assert.Equal(t, region.Name, unmarshaled.Name)
	})

	t.Run("CBIRegions from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "eu-central-1",
			"name": "EU Central (Frankfurt)"
		}`

		var region CBIRegions
		err := json.Unmarshal([]byte(apiResponse), &region)
		require.NoError(t, err)

		assert.Equal(t, "eu-central-1", region.ID)
		assert.Equal(t, "EU Central (Frankfurt)", region.Name)
	})
}

// TestCBIRegions_ResponseParsing tests parsing of API responses
func TestCBIRegions_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse regions list response", func(t *testing.T) {
		response := `[
			{"id": "us-west-1", "name": "US West (Oregon)"},
			{"id": "us-east-1", "name": "US East (Virginia)"},
			{"id": "eu-west-1", "name": "EU West (Ireland)"},
			{"id": "eu-central-1", "name": "EU Central (Frankfurt)"},
			{"id": "ap-southeast-1", "name": "Asia Pacific (Singapore)"},
			{"id": "ap-northeast-1", "name": "Asia Pacific (Tokyo)"}
		]`

		var regions []CBIRegions
		err := json.Unmarshal([]byte(response), &regions)
		require.NoError(t, err)

		assert.Len(t, regions, 6)
		assert.Equal(t, "us-west-1", regions[0].ID)
		assert.Equal(t, "Asia Pacific (Tokyo)", regions[5].Name)
	})
}

// TestCBIRegions_MockServerOperations tests operations
func TestCBIRegions_MockServerOperations(t *testing.T) {
	t.Run("GET all regions", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/regions")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[
				{"id": "us-west-1", "name": "US West"},
				{"id": "eu-central-1", "name": "EU Central"}
			]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/regions")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestCBIRegions_AllRegions tests various region identifiers
func TestCBIRegions_AllRegions(t *testing.T) {
	t.Parallel()

	t.Run("Common CBI region identifiers", func(t *testing.T) {
		commonRegions := []struct {
			id   string
			name string
		}{
			{"us-west-1", "US West"},
			{"us-east-1", "US East"},
			{"eu-west-1", "EU West"},
			{"eu-central-1", "EU Central"},
			{"ap-southeast-1", "Asia Pacific Southeast"},
			{"ap-northeast-1", "Asia Pacific Northeast"},
			{"ap-south-1", "Asia Pacific South"},
			{"sa-east-1", "South America East"},
		}

		for _, r := range commonRegions {
			region := CBIRegions{
				ID:   r.id,
				Name: r.name,
			}

			data, err := json.Marshal(region)
			require.NoError(t, err)

			var unmarshaled CBIRegions
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, r.id, unmarshaled.ID)
		}
	})
}

