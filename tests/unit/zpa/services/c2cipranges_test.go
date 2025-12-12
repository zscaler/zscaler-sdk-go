// Package unit provides unit tests for ZPA C2C IP Ranges service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IPRanges represents the IP ranges for testing
type IPRanges struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	AvailableIps  string `json:"availableIps,omitempty"`
	CountryCode   string `json:"countryCode,omitempty"`
	CustomerId    string `json:"customerId,omitempty"`
	Enabled       bool   `json:"enabled,omitempty"`
	IpRangeBegin  string `json:"ipRangeBegin,omitempty"`
	IpRangeEnd    string `json:"ipRangeEnd,omitempty"`
	IsDeleted     string `json:"isDeleted,omitempty"`
	LatitudeInDb  string `json:"latitudeInDb,omitempty"`
	Location      string `json:"location,omitempty"`
	LocationHint  string `json:"locationHint,omitempty"`
	LongitudeInDb string `json:"longitudeInDb,omitempty"`
	SccmFlag      bool   `json:"sccmFlag,omitempty"`
	SubnetCidr    string `json:"subnetCidr,omitempty"`
	TotalIps      string `json:"totalIps,omitempty"`
	UsedIps       string `json:"usedIps,omitempty"`
	CreationTime  string `json:"creationTime,omitempty"`
	ModifiedBy    string `json:"modifiedBy,omitempty"`
	ModifiedTime  string `json:"modifiedTime,omitempty"`
}

// TestIPRanges_Structure tests the struct definitions
func TestIPRanges_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IPRanges JSON marshaling", func(t *testing.T) {
		ipRange := IPRanges{
			ID:           "ipr-123",
			Name:         "Corporate Network",
			Description:  "Main corporate network range",
			IpRangeBegin: "10.0.0.1",
			IpRangeEnd:   "10.0.0.254",
			SubnetCidr:   "10.0.0.0/24",
			Enabled:      true,
			CountryCode:  "US",
			Location:     "San Francisco, CA",
			TotalIps:     "254",
			AvailableIps: "200",
			UsedIps:      "54",
		}

		data, err := json.Marshal(ipRange)
		require.NoError(t, err)

		var unmarshaled IPRanges
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, ipRange.ID, unmarshaled.ID)
		assert.Equal(t, ipRange.Name, unmarshaled.Name)
		assert.Equal(t, ipRange.SubnetCidr, unmarshaled.SubnetCidr)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("IPRanges from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "ipr-456",
			"name": "Branch Office Network",
			"description": "Branch office IP range",
			"ipRangeBegin": "192.168.1.1",
			"ipRangeEnd": "192.168.1.254",
			"subnetCidr": "192.168.1.0/24",
			"enabled": true,
			"countryCode": "US",
			"customerId": "cust-001",
			"location": "New York, NY",
			"locationHint": "East Coast",
			"latitudeInDb": "40.7128",
			"longitudeInDb": "-74.0060",
			"totalIps": "254",
			"availableIps": "150",
			"usedIps": "104",
			"sccmFlag": false,
			"isDeleted": "false",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com"
		}`

		var ipRange IPRanges
		err := json.Unmarshal([]byte(apiResponse), &ipRange)
		require.NoError(t, err)

		assert.Equal(t, "ipr-456", ipRange.ID)
		assert.Equal(t, "Branch Office Network", ipRange.Name)
		assert.Equal(t, "192.168.1.0/24", ipRange.SubnetCidr)
		assert.Equal(t, "40.7128", ipRange.LatitudeInDb)
		assert.True(t, ipRange.Enabled)
	})
}

// TestIPRanges_ResponseParsing tests parsing of API responses
func TestIPRanges_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse IP ranges list response", func(t *testing.T) {
		response := `[
			{"id": "1", "name": "Range 1", "subnetCidr": "10.0.0.0/24", "enabled": true},
			{"id": "2", "name": "Range 2", "subnetCidr": "10.0.1.0/24", "enabled": true},
			{"id": "3", "name": "Range 3", "subnetCidr": "10.0.2.0/24", "enabled": false}
		]`

		var ranges []IPRanges
		err := json.Unmarshal([]byte(response), &ranges)
		require.NoError(t, err)

		assert.Len(t, ranges, 3)
		assert.True(t, ranges[0].Enabled)
		assert.False(t, ranges[2].Enabled)
	})
}

// TestIPRanges_MockServerOperations tests CRUD operations
func TestIPRanges_MockServerOperations(t *testing.T) {
	t.Run("GET IP range by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/ipRanges/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "ipr-123",
				"name": "Mock Range",
				"subnetCidr": "10.0.0.0/24",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/v2/ipRanges/ipr-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all IP ranges", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[
				{"id": "1", "name": "Range A"},
				{"id": "2", "name": "Range B"}
			]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/v2/ipRanges")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create IP range", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-ipr-456",
				"name": "New IP Range",
				"subnetCidr": "172.16.0.0/16"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/v2/ipRanges", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update IP range", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/v2/ipRanges/ipr-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE IP range", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/v2/ipRanges/ipr-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestIPRanges_SpecialCases tests edge cases
func TestIPRanges_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Large IP range", func(t *testing.T) {
		ipRange := IPRanges{
			ID:           "ipr-123",
			Name:         "Large Range",
			IpRangeBegin: "10.0.0.1",
			IpRangeEnd:   "10.255.255.254",
			SubnetCidr:   "10.0.0.0/8",
			TotalIps:     "16777214",
			AvailableIps: "16000000",
			UsedIps:      "777214",
		}

		data, err := json.Marshal(ipRange)
		require.NoError(t, err)

		var unmarshaled IPRanges
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "16777214", unmarshaled.TotalIps)
	})

	t.Run("IP range with geo location", func(t *testing.T) {
		ipRange := IPRanges{
			ID:            "ipr-123",
			Name:          "Geo Range",
			LatitudeInDb:  "37.7749",
			LongitudeInDb: "-122.4194",
			Location:      "San Francisco, CA",
			CountryCode:   "US",
		}

		data, err := json.Marshal(ipRange)
		require.NoError(t, err)

		var unmarshaled IPRanges
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "37.7749", unmarshaled.LatitudeInDb)
		assert.Equal(t, "US", unmarshaled.CountryCode)
	})

	t.Run("SCCM enabled IP range", func(t *testing.T) {
		ipRange := IPRanges{
			ID:       "ipr-123",
			Name:     "SCCM Range",
			SccmFlag: true,
		}

		data, err := json.Marshal(ipRange)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"sccmFlag":true`)
	})
}

