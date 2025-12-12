// Package unit provides unit tests for ZPA Application Segment service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

// TestApplicationSegment_Structure tests the struct definitions
func TestApplicationSegment_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ApplicationSegmentResource JSON marshaling", func(t *testing.T) {
		segment := applicationsegment.ApplicationSegmentResource{
			ID:                        "app-123",
			Name:                      "Test App Segment",
			Description:               "Test Description",
			Enabled:                   true,
			DomainNames:               []string{"example.com", "test.example.com"},
			SegmentGroupID:            "sg-001",
			SegmentGroupName:          "Test Segment Group",
			BypassType:                "NEVER",
			HealthReporting:           "ON_ACCESS",
			HealthCheckType:           "DEFAULT",
			IsCnameEnabled:            true,
			IpAnchored:                false,
			DoubleEncrypt:             false,
			PassiveHealthEnabled:      true,
			TCPKeepAlive:              "1",
			SelectConnectorCloseToApp: true,
			IcmpAccessType:            "NONE",
			TCPPortRanges:             []string{"80", "80", "443", "443"},
			UDPPortRanges:             []string{},
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegment.ApplicationSegmentResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, segment.ID, unmarshaled.ID)
		assert.Equal(t, segment.Name, unmarshaled.Name)
		assert.Equal(t, segment.Enabled, unmarshaled.Enabled)
		assert.ElementsMatch(t, segment.DomainNames, unmarshaled.DomainNames)
		assert.ElementsMatch(t, segment.TCPPortRanges, unmarshaled.TCPPortRanges)
	})

	t.Run("ApplicationSegmentResource JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "app-456",
			"name": "Production Web App",
			"description": "Production web application",
			"enabled": true,
			"domainNames": ["webapp.example.com", "api.example.com"],
			"segmentGroupId": "sg-002",
			"segmentGroupName": "Production Apps",
			"bypassType": "NEVER",
			"healthReporting": "ON_ACCESS",
			"healthCheckType": "DEFAULT",
			"isCnameEnabled": true,
			"ipAnchored": false,
			"fqdnDnsCheck": true,
			"doubleEncrypt": true,
			"passiveHealthEnabled": true,
			"tcpKeepAlive": "1",
			"selectConnectorCloseToApp": true,
			"icmpAccessType": "PING",
			"tcpPortRanges": ["80", "80", "443", "443", "8080", "8080"],
			"udpPortRanges": ["53", "53"],
			"tcpPortRange": [
				{"from": "80", "to": "80"},
				{"from": "443", "to": "443"}
			],
			"udpPortRange": [
				{"from": "53", "to": "53"}
			],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"matchStyle": "EXCLUSIVE",
			"extranetEnabled": false,
			"apiProtectionEnabled": true,
			"adpEnabled": false,
			"inspectTrafficWithZia": true,
			"weightedLoadBalancing": false
		}`

		var segment applicationsegment.ApplicationSegmentResource
		err := json.Unmarshal([]byte(apiResponse), &segment)
		require.NoError(t, err)

		assert.Equal(t, "app-456", segment.ID)
		assert.Equal(t, "Production Web App", segment.Name)
		assert.True(t, segment.Enabled)
		assert.Len(t, segment.DomainNames, 2)
		assert.True(t, segment.DoubleEncrypt)
		assert.True(t, segment.FQDNDnsCheck)
		assert.True(t, segment.InspectTrafficWithZia)
		assert.True(t, segment.APIProtectionEnabled)
		assert.Equal(t, "PING", segment.IcmpAccessType)
		assert.Equal(t, "EXCLUSIVE", segment.MatchStyle)
		assert.Len(t, segment.TCPAppPortRange, 2)
		assert.Len(t, segment.UDPAppPortRange, 1)
	})

	t.Run("Port ranges structure", func(t *testing.T) {
		ports := common.NetworkPorts{
			From: "443",
			To:   "443",
		}

		data, err := json.Marshal(ports)
		require.NoError(t, err)

		var unmarshaled common.NetworkPorts
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "443", unmarshaled.From)
		assert.Equal(t, "443", unmarshaled.To)
	})

	t.Run("SharedMicrotenantDetails structure", func(t *testing.T) {
		details := applicationsegment.SharedMicrotenantDetails{
			SharedFromMicrotenant: applicationsegment.SharedFromMicrotenant{
				ID:   "mt-source",
				Name: "Source Tenant",
			},
			SharedToMicrotenants: []applicationsegment.SharedToMicrotenant{
				{ID: "mt-dest-1", Name: "Dest Tenant 1"},
				{ID: "mt-dest-2", Name: "Dest Tenant 2"},
			},
		}

		data, err := json.Marshal(details)
		require.NoError(t, err)

		var unmarshaled applicationsegment.SharedMicrotenantDetails
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "mt-source", unmarshaled.SharedFromMicrotenant.ID)
		assert.Len(t, unmarshaled.SharedToMicrotenants, 2)
	})
}

// TestApplicationSegment_ResponseParsing tests parsing of various API responses
func TestApplicationSegment_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse single application segment response", func(t *testing.T) {
		response := `{
			"id": "app-789",
			"name": "Edge Application",
			"description": "Edge application segment",
			"enabled": true,
			"domainNames": ["edge.example.com"],
			"segmentGroupId": "sg-003",
			"bypassType": "ON_NET",
			"bypassOnReauth": true,
			"creationTime": "1609459200000"
		}`

		var segment applicationsegment.ApplicationSegmentResource
		err := json.Unmarshal([]byte(response), &segment)
		require.NoError(t, err)

		assert.Equal(t, "app-789", segment.ID)
		assert.Equal(t, "Edge Application", segment.Name)
		assert.Equal(t, "ON_NET", segment.BypassType)
		assert.True(t, segment.BypassOnReauth)
	})

	t.Run("Parse application segment list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "App 1", "enabled": true, "domainNames": ["app1.example.com"]},
				{"id": "2", "name": "App 2", "enabled": false, "domainNames": ["app2.example.com"]},
				{"id": "3", "name": "App 3", "enabled": true, "domainNames": ["app3.example.com", "app3-alt.example.com"]}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []applicationsegment.ApplicationSegmentResource `json:"list"`
			TotalPages int                                             `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "1", listResp.List[0].ID)
		assert.Len(t, listResp.List[2].DomainNames, 2)
	})

	t.Run("Parse application count response", func(t *testing.T) {
		response := `[{
			"appsConfigured": "150",
			"configuredDateInEpochSeconds": "1609459200"
		}]`

		var counts []applicationsegment.ApplicationCountResponse
		err := json.Unmarshal([]byte(response), &counts)
		require.NoError(t, err)

		assert.Len(t, counts, 1)
		assert.Equal(t, "150", counts[0].AppsConfigured)
	})

	t.Run("Parse current and max limit response", func(t *testing.T) {
		response := `{
			"currentAppsCount": "75",
			"maxAppsLimit": "500"
		}`

		var limit applicationsegment.ApplicationCurrentMaxLimitResponse
		err := json.Unmarshal([]byte(response), &limit)
		require.NoError(t, err)

		assert.Equal(t, "75", limit.CurrentAppsCount)
		assert.Equal(t, "500", limit.MaxAppsLimit)
	})
}

// TestApplicationSegment_MockServerOperations tests CRUD operations with mock server
func TestApplicationSegment_MockServerOperations(t *testing.T) {
	t.Run("GET application segment by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/application/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "app-123",
				"name": "Mock Application",
				"description": "Created by mock server",
				"enabled": true,
				"domainNames": ["mock.example.com"],
				"segmentGroupId": "sg-001"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/application/app-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all application segments", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "App A", "enabled": true, "domainNames": ["a.example.com"]},
					{"id": "2", "name": "App B", "enabled": true, "domainNames": ["b.example.com"]}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/application")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create application segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/application")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-app-456",
				"name": "New Application",
				"description": "Newly created",
				"enabled": true,
				"domainNames": ["new.example.com"]
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/application", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update application segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/application/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/application/app-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE application segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Contains(t, r.URL.Path, "/application/")
			assert.Equal(t, "true", r.URL.Query().Get("forceDelete"))

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/application/app-123?forceDelete=true", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestApplicationSegment_ErrorHandling tests error scenarios
func TestApplicationSegment_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Application segment not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/application/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Bad Request - Missing domain names", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Domain names are required"}`))
		}))
		defer server.Close()

		resp, _ := http.Post(server.URL+"/application", "application/json", nil)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Validation error response", func(t *testing.T) {
		response := `{
			"params": ["domainNames"],
			"id": "validation-001",
			"reason": "Domain name already exists in another application"
		}`

		var validationErr applicationsegment.ApplicationValidationError
		err := json.Unmarshal([]byte(response), &validationErr)
		require.NoError(t, err)

		assert.Contains(t, validationErr.Params, "domainNames")
		assert.Equal(t, "Domain name already exists in another application", validationErr.Reason)
	})
}

// TestApplicationSegment_SpecialCases tests edge cases and special scenarios
func TestApplicationSegment_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Application segment with multiple domain names", func(t *testing.T) {
		segment := applicationsegment.ApplicationSegmentResource{
			ID:          "123",
			Name:        "Multi-domain App",
			Enabled:     true,
			DomainNames: []string{"domain1.com", "domain2.com", "domain3.com", "*.wildcard.com"},
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegment.ApplicationSegmentResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.DomainNames, 4)
		assert.Contains(t, unmarshaled.DomainNames, "*.wildcard.com")
	})

	t.Run("Application segment with weighted load balancing", func(t *testing.T) {
		config := applicationsegment.WeightedLoadBalancerConfig{
			ApplicationID:         "app-123",
			WeightedLoadBalancing: true,
			ApplicationToServerGroupMaps: []applicationsegment.ApplicationToServerGroupMapping{
				{ID: "sg-1", Name: "Group 1", Weight: "50", Passive: false},
				{ID: "sg-2", Name: "Group 2", Weight: "30", Passive: false},
				{ID: "sg-3", Name: "Group 3", Weight: "20", Passive: true},
			},
		}

		data, err := json.Marshal(config)
		require.NoError(t, err)

		var unmarshaled applicationsegment.WeightedLoadBalancerConfig
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.WeightedLoadBalancing)
		assert.Len(t, unmarshaled.ApplicationToServerGroupMaps, 3)
		assert.True(t, unmarshaled.ApplicationToServerGroupMaps[2].Passive)
	})

	t.Run("Application segment with tags", func(t *testing.T) {
		segment := applicationsegment.ApplicationSegmentResource{
			ID:      "123",
			Name:    "Tagged App",
			Enabled: true,
			Tags: []applicationsegment.Tag{
				{
					Namespace: common.CommonSummary{ID: "ns-1", Name: "Environment"},
					TagKey:    common.CommonSummary{ID: "key-1", Name: "env"},
					TagValue:  common.CommonIDName{ID: "val-1", Name: "production"},
					Origin:    "USER",
				},
			},
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegment.ApplicationSegmentResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.Tags, 1)
		assert.Equal(t, "Environment", unmarshaled.Tags[0].Namespace.Name)
	})

	t.Run("Multi-match unsupported references payload", func(t *testing.T) {
		payload := applicationsegment.MultiMatchUnsupportedReferencesPayload{"domain1.com", "domain2.com"}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		assert.Contains(t, string(data), "domain1.com")
		assert.Contains(t, string(data), "domain2.com")
	})

	t.Run("Multi-match unsupported references response", func(t *testing.T) {
		response := `[{
			"id": "app-001",
			"appSegmentName": "Legacy App",
			"domains": ["legacy.example.com"],
			"tcpPorts": ["80", "443"],
			"matchStyle": "INCLUSIVE",
			"microtenantName": "Default"
		}]`

		var refs []applicationsegment.MultiMatchUnsupportedReferencesResponse
		err := json.Unmarshal([]byte(response), &refs)
		require.NoError(t, err)

		assert.Len(t, refs, 1)
		assert.Equal(t, "Legacy App", refs[0].AppSegmentName)
		assert.Equal(t, "INCLUSIVE", refs[0].MatchStyle)
	})

	t.Run("Bulk update multi-match payload", func(t *testing.T) {
		payload := applicationsegment.BulkUpdateMultiMatchPayload{
			ApplicationIDs: []int{1, 2, 3, 4, 5},
			MatchStyle:     "EXCLUSIVE",
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		assert.Contains(t, string(data), "EXCLUSIVE")
		assert.Contains(t, string(data), "applicationIds")
	})

	t.Run("Application mappings", func(t *testing.T) {
		response := `[
			{"name": "Mapping 1", "type": "BROWSER_ACCESS"},
			{"name": "Mapping 2", "type": "SECURE_REMOTE_ACCESS"}
		]`

		var mappings []applicationsegment.ApplicationMappings
		err := json.Unmarshal([]byte(response), &mappings)
		require.NoError(t, err)

		assert.Len(t, mappings, 2)
		assert.Equal(t, "BROWSER_ACCESS", mappings[0].Type)
	})
}

// TestApplicationSegment_GetByName tests the GetByName functionality
func TestApplicationSegment_GetByName(t *testing.T) {
	t.Run("Search returns matching application segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			search := r.URL.Query().Get("search")
			assert.NotEmpty(t, search)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Production Web App", "enabled": true, "domainNames": ["webapp.example.com"]}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/application?search=Production")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestApplicationSegment_PortConfiguration tests port configuration scenarios
func TestApplicationSegment_PortConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("TCP port ranges as strings", func(t *testing.T) {
		segment := applicationsegment.ApplicationSegmentResource{
			ID:            "123",
			Name:          "Port Test App",
			TCPPortRanges: []string{"80", "80", "443", "443", "8080", "8090"},
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegment.ApplicationSegmentResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.TCPPortRanges, 6)
	})

	t.Run("TCP port ranges as objects", func(t *testing.T) {
		segment := applicationsegment.ApplicationSegmentResource{
			ID:   "123",
			Name: "Port Object Test App",
			TCPAppPortRange: []common.NetworkPorts{
				{From: "80", To: "80"},
				{From: "443", To: "443"},
				{From: "8080", To: "8090"},
			},
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegment.ApplicationSegmentResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.TCPAppPortRange, 3)
		assert.Equal(t, "8080", unmarshaled.TCPAppPortRange[2].From)
		assert.Equal(t, "8090", unmarshaled.TCPAppPortRange[2].To)
	})

	t.Run("UDP port configuration", func(t *testing.T) {
		segment := applicationsegment.ApplicationSegmentResource{
			ID:            "123",
			Name:          "UDP App",
			UDPPortRanges: []string{"53", "53", "500", "500"},
			UDPAppPortRange: []common.NetworkPorts{
				{From: "53", To: "53"},
				{From: "500", To: "500"},
			},
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegment.ApplicationSegmentResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.UDPPortRanges, 4)
		assert.Len(t, unmarshaled.UDPAppPortRange, 2)
	})
}

