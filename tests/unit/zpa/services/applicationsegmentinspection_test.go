// Package unit provides unit tests for ZPA Application Segment Inspection service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AppSegmentInspection represents the inspection app segment for testing
type AppSegmentInspection struct {
	ID                        string                      `json:"id,omitempty"`
	Name                      string                      `json:"name,omitempty"`
	Description               string                      `json:"description,omitempty"`
	SegmentGroupID            string                      `json:"segmentGroupId,omitempty"`
	SegmentGroupName          string                      `json:"segmentGroupName,omitempty"`
	BypassType                string                      `json:"bypassType,omitempty"`
	DomainNames               []string                    `json:"domainNames,omitempty"`
	Enabled                   bool                        `json:"enabled"`
	AdpEnabled                bool                        `json:"adpEnabled,omitempty"`
	ICMPAccessType            string                      `json:"icmpAccessType,omitempty"`
	PassiveHealthEnabled      bool                        `json:"passiveHealthEnabled,omitempty"`
	FQDNDnsCheck              bool                        `json:"fqdnDnsCheck"`
	APIProtectionEnabled      bool                        `json:"apiProtectionEnabled"`
	MatchStyle                string                      `json:"matchStyle,omitempty"`
	SelectConnectorCloseToApp bool                        `json:"selectConnectorCloseToApp"`
	DoubleEncrypt             bool                        `json:"doubleEncrypt"`
	HealthCheckType           string                      `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                        `json:"isCnameEnabled"`
	IPAnchored                bool                        `json:"ipAnchored"`
	HealthReporting           string                      `json:"healthReporting,omitempty"`
	TCPKeepAlive              string                      `json:"tcpKeepAlive,omitempty"`
	TCPPortRanges             []string                    `json:"tcpPortRanges,omitempty"`
	UDPPortRanges             []string                    `json:"udpPortRanges,omitempty"`
	TCPProtocols              []string                    `json:"tcpProtocols"`
	InspectionAppDto          []InspectionAppDto          `json:"inspectionApps,omitempty"`
	CommonAppsDto             InspectionCommonAppsDto     `json:"commonAppsDto,omitempty"`
	MicroTenantID             string                      `json:"microtenantId,omitempty"`
	MicroTenantName           string                      `json:"microtenantName,omitempty"`
	CreationTime              string                      `json:"creationTime,omitempty"`
	ModifiedBy                string                      `json:"modifiedBy,omitempty"`
	ModifiedTime              string                      `json:"modifiedTime,omitempty"`
}

// InspectionAppDto represents inspection app configuration
type InspectionAppDto struct {
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId,omitempty"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled"`
	ApplicationPort     string   `json:"applicationPort,omitempty"`
	ApplicationProtocol string   `json:"applicationProtocol,omitempty"`
	CertificateID       string   `json:"certificateId,omitempty"`
	CertificateName     string   `json:"certificateName,omitempty"`
	Domain              string   `json:"domain,omitempty"`
	Protocols           []string `json:"protocols,omitempty"`
	TrustUntrustedCert  bool     `json:"trustUntrustedCert"`
}

// InspectionCommonAppsDto represents common apps configuration
type InspectionCommonAppsDto struct {
	AppsConfig         []InspectionAppsConfig `json:"appsConfig,omitempty"`
	DeletedInspectApps []string               `json:"deletedInspectApps,omitempty"`
}

// InspectionAppsConfig represents app config for inspection
type InspectionAppsConfig struct {
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId,omitempty"`
	InspectAppID        string   `json:"inspectAppId"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled"`
	AdpEnabled          bool     `json:"adpEnabled"`
	AllowOptions        bool     `json:"allowOptions"`
	AppTypes            []string `json:"appTypes,omitempty"`
	ApplicationPort     string   `json:"applicationPort,omitempty"`
	ApplicationProtocol string   `json:"applicationProtocol,omitempty"`
	CertificateID       string   `json:"certificateId,omitempty"`
	Domain              string   `json:"domain,omitempty"`
	Hidden              bool     `json:"hidden"`
	TrustUntrustedCert  bool     `json:"trustUntrustedCert"`
}

// TestAppSegmentInspection_Structure tests the struct definitions
func TestAppSegmentInspection_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentInspection JSON marshaling", func(t *testing.T) {
		inspection := AppSegmentInspection{
			ID:               "insp-123",
			Name:             "Web Inspection App",
			Description:      "Inspection for web application",
			SegmentGroupID:   "sg-001",
			SegmentGroupName: "Inspection Apps",
			BypassType:       "NEVER",
			DomainNames:      []string{"webapp.example.com"},
			Enabled:          true,
			HealthReporting:  "ON_ACCESS",
			TCPPortRanges:    []string{"443", "443"},
			TCPProtocols:     []string{"HTTPS"},
			InspectionAppDto: []InspectionAppDto{
				{
					ID:                  "ia-001",
					Name:                "Inspection Sub-App",
					ApplicationPort:     "443",
					ApplicationProtocol: "HTTPS",
					Enabled:             true,
				},
			},
		}

		data, err := json.Marshal(inspection)
		require.NoError(t, err)

		var unmarshaled AppSegmentInspection
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, inspection.ID, unmarshaled.ID)
		assert.Equal(t, inspection.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.InspectionAppDto, 1)
	})

	t.Run("AppSegmentInspection from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "insp-456",
			"name": "DLP Inspection App",
			"description": "DLP inspection application",
			"segmentGroupId": "sg-002",
			"segmentGroupName": "DLP Apps",
			"bypassType": "NEVER",
			"domainNames": ["dlp.example.com"],
			"enabled": true,
			"adpEnabled": true,
			"passiveHealthEnabled": true,
			"fqdnDnsCheck": true,
			"apiProtectionEnabled": false,
			"matchStyle": "EXCLUSIVE",
			"selectConnectorCloseToApp": true,
			"doubleEncrypt": false,
			"healthCheckType": "DEFAULT",
			"isCnameEnabled": true,
			"ipAnchored": false,
			"healthReporting": "ON_ACCESS",
			"tcpPortRanges": ["443", "443"],
			"tcpProtocols": ["HTTPS"],
			"inspectionApps": [
				{
					"id": "ia-001",
					"appId": "insp-456",
					"name": "DLP Sub-App",
					"enabled": true,
					"applicationPort": "443",
					"applicationProtocol": "HTTPS",
					"domain": "dlp.example.com",
					"trustUntrustedCert": false
				}
			],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var inspection AppSegmentInspection
		err := json.Unmarshal([]byte(apiResponse), &inspection)
		require.NoError(t, err)

		assert.Equal(t, "insp-456", inspection.ID)
		assert.Equal(t, "DLP Inspection App", inspection.Name)
		assert.True(t, inspection.AdpEnabled)
		assert.True(t, inspection.Enabled)
		assert.Len(t, inspection.InspectionAppDto, 1)
		assert.Equal(t, "HTTPS", inspection.InspectionAppDto[0].ApplicationProtocol)
	})

	t.Run("InspectionAppsConfig structure", func(t *testing.T) {
		config := InspectionAppsConfig{
			ID:                  "config-001",
			AppID:               "insp-001",
			InspectAppID:        "ia-001",
			Name:                "Config App",
			Enabled:             true,
			AdpEnabled:          true,
			AllowOptions:        true,
			AppTypes:            []string{"INSPECT"},
			ApplicationPort:     "443",
			ApplicationProtocol: "HTTPS",
			TrustUntrustedCert:  false,
		}

		data, err := json.Marshal(config)
		require.NoError(t, err)

		var unmarshaled InspectionAppsConfig
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, config.InspectAppID, unmarshaled.InspectAppID)
		assert.True(t, unmarshaled.AdpEnabled)
	})
}

// TestAppSegmentInspection_MockServerOperations tests CRUD operations
func TestAppSegmentInspection_MockServerOperations(t *testing.T) {
	t.Run("GET inspection app by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/application/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "insp-123",
				"name": "Mock Inspection App",
				"enabled": true,
				"inspectionApps": [{"id": "ia-1"}]
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/application/insp-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create inspection app", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-insp-456",
				"name": "New Inspection App",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update inspection app", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/application/insp-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE inspection app with forceDelete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "true", r.URL.Query().Get("forceDelete"))
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/application/insp-123?forceDelete=true", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestAppSegmentInspection_SpecialCases tests edge cases
func TestAppSegmentInspection_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("ADP enabled inspection app", func(t *testing.T) {
		inspection := AppSegmentInspection{
			ID:         "insp-123",
			Name:       "ADP App",
			AdpEnabled: true,
			Enabled:    true,
		}

		data, err := json.Marshal(inspection)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"adpEnabled":true`)
	})

	t.Run("CommonAppsDto with deleted apps", func(t *testing.T) {
		commonApps := InspectionCommonAppsDto{
			AppsConfig: []InspectionAppsConfig{
				{ID: "config-1", Name: "Active App", Enabled: true},
			},
			DeletedInspectApps: []string{"ia-deleted-1", "ia-deleted-2"},
		}

		data, err := json.Marshal(commonApps)
		require.NoError(t, err)

		var unmarshaled InspectionCommonAppsDto
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.DeletedInspectApps, 2)
	})

	t.Run("Multiple inspection protocols", func(t *testing.T) {
		inspection := AppSegmentInspection{
			ID:           "insp-123",
			Name:         "Multi-Protocol Inspection",
			TCPProtocols: []string{"HTTPS", "HTTP"},
		}

		data, err := json.Marshal(inspection)
		require.NoError(t, err)

		var unmarshaled AppSegmentInspection
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.TCPProtocols, 2)
	})
}

