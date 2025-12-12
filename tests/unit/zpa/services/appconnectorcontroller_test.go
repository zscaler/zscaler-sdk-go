// Package unit provides unit tests for ZPA App Connector Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AppConnector represents the app connector structure for testing
type AppConnector struct {
	ID                               string                 `json:"id,omitempty"`
	Name                             string                 `json:"name,omitempty"`
	Description                      string                 `json:"description,omitempty"`
	Enabled                          bool                   `json:"enabled,omitempty"`
	AppConnectorGroupID              string                 `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName            string                 `json:"appConnectorGroupName,omitempty"`
	ControlChannelStatus             string                 `json:"controlChannelStatus,omitempty"`
	CurrentVersion                   string                 `json:"currentVersion,omitempty"`
	ExpectedVersion                  string                 `json:"expectedVersion,omitempty"`
	ExpectedUpgradeTime              string                 `json:"expectedUpgradeTime,omitempty"`
	Fingerprint                      string                 `json:"fingerprint,omitempty"`
	IPACL                            string                 `json:"ipAcl,omitempty"`
	IssuedCertID                     string                 `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                 `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                 `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                 `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                 `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastUpgradeTime                  string                 `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                 `json:"latitude,omitempty"`
	Location                         string                 `json:"location,omitempty"`
	Longitude                        string                 `json:"longitude,omitempty"`
	Platform                         string                 `json:"platform,omitempty"`
	PlatformDetail                   string                 `json:"platformDetail,omitempty"`
	PreviousVersion                  string                 `json:"previousVersion,omitempty"`
	PrivateIP                        string                 `json:"privateIp,omitempty"`
	PublicIP                         string                 `json:"publicIp,omitempty"`
	ProvisioningKeyID                string                 `json:"provisioningKeyId"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName"`
	RuntimeOS                        string                 `json:"runtimeOS,omitempty"`
	SargeVersion                     string                 `json:"sargeVersion,omitempty"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                 `json:"upgradeStatus,omitempty"`
	CtrlBrokerName                   string                 `json:"ctrlBrokerName,omitempty"`
	CreationTime                     string                 `json:"creationTime,omitempty"`
	ModifiedBy                       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                 `json:"modifiedTime,omitempty"`
	MicroTenantID                    string                 `json:"microtenantId,omitempty"`
	MicroTenantName                  string                 `json:"microtenantName,omitempty"`
	ApplicationStartTime             string                 `json:"applicationStartTime,omitempty"`
	AssistantVersion                 AssistantVersion       `json:"assistantVersion,omitempty"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert,omitempty"`
}

// AssistantVersion represents the assistant version structure
type AssistantVersion struct {
	ID                       string `json:"id,omitempty"`
	ApplicationStartTime     string `json:"applicationStartTime,omitempty"`
	AppConnectorGroupID      string `json:"appConnectorGroupId,omitempty"`
	BrokerId                 string `json:"brokerId,omitempty"`
	CreationTime             string `json:"creationTime,omitempty"`
	CtrlChannelStatus        string `json:"ctrlChannelStatus,omitempty"`
	CurrentVersion           string `json:"currentVersion,omitempty"`
	DisableAutoUpdate        bool   `json:"disableAutoUpdate,omitempty"`
	ExpectedVersion          string `json:"expectedVersion,omitempty"`
	LastBrokerConnectTime    string `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerDisconnectTime string `json:"lastBrokerDisconnectTime,omitempty"`
	LastUpgradedTime         string `json:"lastUpgradedTime,omitempty"`
	LoneWarrior              bool   `json:"loneWarrior,omitempty"`
	ModifiedBy               string `json:"modifiedBy,omitempty"`
	ModifiedTime             string `json:"modifiedTime,omitempty"`
	Latitude                 string `json:"latitude,omitempty"`
	Longitude                string `json:"longitude,omitempty"`
	MtunnelID                string `json:"mtunnelId,omitempty"`
	Platform                 string `json:"platform,omitempty"`
	PlatformDetail           string `json:"platformDetail,omitempty"`
	PreviousVersion          string `json:"previousVersion,omitempty"`
	PrivateIP                string `json:"privateIp,omitempty"`
	PublicIP                 string `json:"publicIp,omitempty"`
	RestartTimeInSec         string `json:"restartTimeInSec,omitempty"`
	RuntimeOS                string `json:"runtimeOS,omitempty"`
	SargeVersion             string `json:"sargeVersion,omitempty"`
	SystemStartTime          string `json:"systemStartTime,omitempty"`
	UpgradeAttempt           string `json:"upgradeAttempt,omitempty"`
	UpgradeStatus            string `json:"upgradeStatus,omitempty"`
	UpgradeNowOnce           bool   `json:"upgradeNowOnce,omitempty"`
}

// BulkDeleteRequest represents bulk delete request
type BulkDeleteRequest struct {
	IDs []string `json:"ids"`
}

// TestAppConnector_Structure tests the struct definitions
func TestAppConnector_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppConnector JSON marshaling", func(t *testing.T) {
		connector := AppConnector{
			ID:                    "conn-123",
			Name:                  "Production Connector 1",
			Description:           "Primary production connector",
			Enabled:               true,
			AppConnectorGroupID:   "acg-001",
			AppConnectorGroupName: "Production Group",
			ControlChannelStatus:  "ZPN_STATUS_ONLINE",
			CurrentVersion:        "24.1.2.123",
			Platform:              "el8",
			PlatformDetail:        "Rocky Linux 8.9",
			PrivateIP:             "10.0.0.5",
			PublicIP:              "203.0.113.5",
			RuntimeOS:             "Linux",
			Latitude:              "37.7749",
			Longitude:             "-122.4194",
			Location:              "San Francisco, CA",
			ProvisioningKeyID:     "pk-001",
			ProvisioningKeyName:   "Production Key",
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var unmarshaled AppConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, connector.ID, unmarshaled.ID)
		assert.Equal(t, connector.Name, unmarshaled.Name)
		assert.Equal(t, connector.ControlChannelStatus, unmarshaled.ControlChannelStatus)
		assert.Equal(t, connector.CurrentVersion, unmarshaled.CurrentVersion)
		assert.Equal(t, connector.Platform, unmarshaled.Platform)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("AppConnector JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "conn-456",
			"name": "Edge Connector",
			"description": "Edge location connector",
			"enabled": true,
			"appConnectorGroupId": "acg-002",
			"appConnectorGroupName": "Edge Group",
			"controlChannelStatus": "ZPN_STATUS_ONLINE",
			"currentVersion": "24.1.2.125",
			"expectedVersion": "24.1.2.130",
			"expectedUpgradeTime": "1612137600000",
			"fingerprint": "ABC123DEF456",
			"issuedCertId": "cert-001",
			"lastBrokerConnectTime": "1612137600000",
			"lastBrokerConnectTimeDuration": "2h30m",
			"lastBrokerDisconnectTime": "1612051200000",
			"lastBrokerDisconnectTimeDuration": "5m",
			"lastUpgradeTime": "1609459200000",
			"latitude": "40.7128",
			"longitude": "-74.0060",
			"location": "New York, NY",
			"platform": "el9",
			"platformDetail": "RHEL 9.2",
			"previousVersion": "24.1.2.120",
			"privateIp": "192.168.1.100",
			"publicIp": "198.51.100.1",
			"provisioningKeyId": "pk-002",
			"provisioningKeyName": "Edge Key",
			"runtimeOS": "Linux",
			"sargeVersion": "1.0.0",
			"upgradeAttempt": "0",
			"upgradeStatus": "NOT_SCHEDULED",
			"ctrlBrokerName": "broker-1",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"applicationStartTime": "1609459200000",
			"assistantVersion": {
				"id": "av-001",
				"currentVersion": "24.1.2.125",
				"ctrlChannelStatus": "ONLINE",
				"platform": "el9"
			}
		}`

		var connector AppConnector
		err := json.Unmarshal([]byte(apiResponse), &connector)
		require.NoError(t, err)

		assert.Equal(t, "conn-456", connector.ID)
		assert.Equal(t, "Edge Connector", connector.Name)
		assert.Equal(t, "ZPN_STATUS_ONLINE", connector.ControlChannelStatus)
		assert.Equal(t, "24.1.2.125", connector.CurrentVersion)
		assert.Equal(t, "24.1.2.130", connector.ExpectedVersion)
		assert.Equal(t, "el9", connector.Platform)
		assert.Equal(t, "RHEL 9.2", connector.PlatformDetail)
		assert.Equal(t, "NOT_SCHEDULED", connector.UpgradeStatus)
		assert.NotEmpty(t, connector.AssistantVersion.ID)
	})

	t.Run("AssistantVersion structure", func(t *testing.T) {
		av := AssistantVersion{
			ID:                    "av-001",
			AppConnectorGroupID:   "acg-001",
			CurrentVersion:        "24.1.2.125",
			CtrlChannelStatus:     "ONLINE",
			DisableAutoUpdate:     false,
			ExpectedVersion:       "24.1.2.130",
			Platform:              "el9",
			PrivateIP:             "10.0.0.5",
			PublicIP:              "203.0.113.5",
			RuntimeOS:             "Linux",
			UpgradeStatus:         "NOT_SCHEDULED",
		}

		data, err := json.Marshal(av)
		require.NoError(t, err)

		var unmarshaled AssistantVersion
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, av.ID, unmarshaled.ID)
		assert.Equal(t, av.CurrentVersion, unmarshaled.CurrentVersion)
		assert.False(t, unmarshaled.DisableAutoUpdate)
	})

	t.Run("BulkDeleteRequest structure", func(t *testing.T) {
		req := BulkDeleteRequest{
			IDs: []string{"conn-1", "conn-2", "conn-3"},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled BulkDeleteRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.IDs, 3)
	})
}

// TestAppConnector_ResponseParsing tests parsing of various API responses
func TestAppConnector_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse app connector list response", func(t *testing.T) {
		response := `{
			"list": [
				{
					"id": "1",
					"name": "Connector 1",
					"enabled": true,
					"controlChannelStatus": "ZPN_STATUS_ONLINE",
					"currentVersion": "24.1.2.123"
				},
				{
					"id": "2",
					"name": "Connector 2",
					"enabled": true,
					"controlChannelStatus": "ZPN_STATUS_OFFLINE",
					"currentVersion": "24.1.2.120"
				},
				{
					"id": "3",
					"name": "Connector 3",
					"enabled": false,
					"controlChannelStatus": "ZPN_STATUS_MAINTENANCE",
					"currentVersion": "24.1.2.125"
				}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []AppConnector `json:"list"`
			TotalPages int            `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "ZPN_STATUS_ONLINE", listResp.List[0].ControlChannelStatus)
		assert.Equal(t, "ZPN_STATUS_OFFLINE", listResp.List[1].ControlChannelStatus)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestAppConnector_MockServerOperations tests CRUD operations with mock server
func TestAppConnector_MockServerOperations(t *testing.T) {
	t.Run("GET app connector by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/connector/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "conn-123",
				"name": "Mock Connector",
				"enabled": true,
				"controlChannelStatus": "ZPN_STATUS_ONLINE",
				"currentVersion": "24.1.2.123"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/connector/conn-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all app connectors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Connector A", "enabled": true},
					{"id": "2", "name": "Connector B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/connector")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update app connector", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/connector/")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/connector/conn-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE app connector", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Contains(t, r.URL.Path, "/connector/")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/connector/conn-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("POST bulk delete app connectors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/bulkDelete")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/connector/bulkDelete", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestAppConnector_ErrorHandling tests error scenarios
func TestAppConnector_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 App Connector Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "App connector not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/connector/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Bad Request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Invalid connector ID"}`))
		}))
		defer server.Close()

		resp, _ := http.Get(server.URL + "/connector/invalid")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// TestAppConnector_SpecialCases tests edge cases
func TestAppConnector_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Control channel status values", func(t *testing.T) {
		statuses := []string{
			"ZPN_STATUS_ONLINE",
			"ZPN_STATUS_OFFLINE",
			"ZPN_STATUS_MAINTENANCE",
			"ZPN_STATUS_UNREACHABLE",
			"ZPN_STATUS_UPGRADE_IN_PROGRESS",
			"ZPN_STATUS_UPGRADE_FAILED",
		}

		for _, status := range statuses {
			connector := AppConnector{
				ID:                   "conn-" + status,
				Name:                 status + " Connector",
				ControlChannelStatus: status,
			}

			data, err := json.Marshal(connector)
			require.NoError(t, err)

			assert.Contains(t, string(data), status)
		}
	})

	t.Run("Upgrade status values", func(t *testing.T) {
		upgradeStatuses := []string{
			"NOT_SCHEDULED",
			"SCHEDULED",
			"IN_PROGRESS",
			"COMPLETED",
			"FAILED",
			"PENDING_RESTART",
		}

		for _, status := range upgradeStatuses {
			connector := AppConnector{
				ID:            "conn-" + status,
				Name:          status + " Connector",
				UpgradeStatus: status,
			}

			data, err := json.Marshal(connector)
			require.NoError(t, err)

			assert.Contains(t, string(data), status)
		}
	})

	t.Run("Platform types", func(t *testing.T) {
		platforms := []string{"el7", "el8", "el9", "ubuntu", "docker", "aws", "azure", "gcp"}

		for _, platform := range platforms {
			connector := AppConnector{
				ID:       "conn-" + platform,
				Name:     platform + " Connector",
				Platform: platform,
			}

			data, err := json.Marshal(connector)
			require.NoError(t, err)

			assert.Contains(t, string(data), platform)
		}
	})

	t.Run("Connector with enrollment cert", func(t *testing.T) {
		connector := AppConnector{
			ID:           "conn-123",
			Name:         "Enrolled Connector",
			IssuedCertID: "cert-001",
			EnrollmentCert: map[string]interface{}{
				"id":       "cert-001",
				"name":     "Connector Certificate",
				"issuedTo": "Connector",
			},
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var unmarshaled AppConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.NotNil(t, unmarshaled.EnrollmentCert)
	})

	t.Run("Connector geo location", func(t *testing.T) {
		connector := AppConnector{
			ID:        "conn-123",
			Name:      "Geo Connector",
			Latitude:  "37.7749",
			Longitude: "-122.4194",
			Location:  "San Francisco, CA, USA",
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var unmarshaled AppConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "37.7749", unmarshaled.Latitude)
		assert.Equal(t, "-122.4194", unmarshaled.Longitude)
	})

	t.Run("Disabled connector", func(t *testing.T) {
		connector := AppConnector{
			ID:      "conn-123",
			Name:    "Disabled Connector",
			Enabled: false,
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		// Note: With omitempty, false boolean won't be in output
		var unmarshaled AppConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.False(t, unmarshaled.Enabled)
	})
}

