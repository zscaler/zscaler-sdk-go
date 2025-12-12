// Package unit provides unit tests for ZPA Application Segment Browser Access service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BrowserAccess represents the browser access application segment for testing
type BrowserAccess struct {
	ID                        string                       `json:"id,omitempty"`
	Name                      string                       `json:"name,omitempty"`
	Description               string                       `json:"description,omitempty"`
	SegmentGroupID            string                       `json:"segmentGroupId,omitempty"`
	SegmentGroupName          string                       `json:"segmentGroupName,omitempty"`
	BypassType                string                       `json:"bypassType,omitempty"`
	BypassOnReauth            bool                         `json:"bypassOnReauth,omitempty"`
	ExtranetEnabled           bool                         `json:"extranetEnabled"`
	MatchStyle                string                       `json:"matchStyle,omitempty"`
	ConfigSpace               string                       `json:"configSpace,omitempty"`
	DomainNames               []string                     `json:"domainNames,omitempty"`
	Enabled                   bool                         `json:"enabled"`
	PassiveHealthEnabled      bool                         `json:"passiveHealthEnabled"`
	FQDNDnsCheck              bool                         `json:"fqdnDnsCheck"`
	APIProtectionEnabled      bool                         `json:"apiProtectionEnabled"`
	SelectConnectorCloseToApp bool                         `json:"selectConnectorCloseToApp"`
	DoubleEncrypt             bool                         `json:"doubleEncrypt"`
	HealthCheckType           string                       `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                         `json:"isCnameEnabled"`
	IPAnchored                bool                         `json:"ipAnchored"`
	TCPKeepAlive              string                       `json:"tcpKeepAlive,omitempty"`
	IsIncompleteDRConfig      bool                         `json:"isIncompleteDRConfig"`
	UseInDrMode               bool                         `json:"useInDrMode"`
	InspectTrafficWithZia     bool                         `json:"inspectTrafficWithZia"`
	MicroTenantID             string                       `json:"microtenantId,omitempty"`
	MicroTenantName           string                       `json:"microtenantName,omitempty"`
	HealthReporting           string                       `json:"healthReporting,omitempty"`
	ICMPAccessType            string                       `json:"icmpAccessType,omitempty"`
	CreationTime              string                       `json:"creationTime,omitempty"`
	ModifiedBy                string                       `json:"modifiedBy,omitempty"`
	ModifiedTime              string                       `json:"modifiedTime,omitempty"`
	TCPPortRanges             []string                     `json:"tcpPortRanges,omitempty"`
	UDPPortRanges             []string                     `json:"udpPortRanges,omitempty"`
	ClientlessApps            []ClientlessApps             `json:"clientlessApps,omitempty"`
	SharedMicrotenantDetails  BrowserAccessSharedDetails   `json:"sharedMicrotenantDetails,omitempty"`
}

// ClientlessApps represents clientless app configuration
type ClientlessApps struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	AppID               string `json:"appId,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	Cname               string `json:"cname,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	Hidden              bool   `json:"hidden"`
	LocalDomain         string `json:"localDomain,omitempty"`
	Path                string `json:"path,omitempty"`
	TrustUntrustedCert  bool   `json:"trustUntrustedCert"`
	AllowOptions        bool   `json:"allowOptions"`
	ExtDomain           string `json:"extDomain"`
	ExtLabel            string `json:"extLabel"`
	ExtDomainName       string `json:"extDomainName"`
	ExtID               string `json:"extId"`
	CreationTime        string `json:"creationTime,omitempty"`
	ModifiedBy          string `json:"modifiedBy,omitempty"`
	ModifiedTime        string `json:"modifiedTime,omitempty"`
	MicroTenantID       string `json:"microtenantId,omitempty"`
	MicroTenantName     string `json:"microtenantName,omitempty"`
}

// BrowserAccessSharedDetails represents shared microtenant details
type BrowserAccessSharedDetails struct {
	SharedFromMicrotenant BrowserAccessSharedFrom   `json:"sharedFromMicrotenant,omitempty"`
	SharedToMicrotenants  []BrowserAccessSharedTo   `json:"sharedToMicrotenants,omitempty"`
}

type BrowserAccessSharedFrom struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type BrowserAccessSharedTo struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TestBrowserAccess_Structure tests the struct definitions
func TestBrowserAccess_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BrowserAccess JSON marshaling", func(t *testing.T) {
		ba := BrowserAccess{
			ID:               "ba-123",
			Name:             "Internal Portal",
			Description:      "Internal web portal",
			SegmentGroupID:   "sg-001",
			SegmentGroupName: "Internal Apps",
			BypassType:       "NEVER",
			DomainNames:      []string{"portal.internal.example.com"},
			Enabled:          true,
			HealthReporting:  "ON_ACCESS",
			HealthCheckType:  "DEFAULT",
			IsCnameEnabled:   true,
			TCPKeepAlive:     "1",
			TCPPortRanges:    []string{"443", "443"},
			ClientlessApps: []ClientlessApps{
				{
					ID:                  "ca-001",
					Name:                "Portal App",
					ApplicationPort:     "443",
					ApplicationProtocol: "HTTPS",
					CertificateID:       "cert-001",
					Domain:              "portal.internal.example.com",
					Enabled:             true,
				},
			},
		}

		data, err := json.Marshal(ba)
		require.NoError(t, err)

		var unmarshaled BrowserAccess
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, ba.ID, unmarshaled.ID)
		assert.Equal(t, ba.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.ClientlessApps, 1)
		assert.Equal(t, "443", unmarshaled.ClientlessApps[0].ApplicationPort)
	})

	t.Run("BrowserAccess JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "ba-456",
			"name": "HR Portal",
			"description": "Human Resources Portal",
			"segmentGroupId": "sg-002",
			"segmentGroupName": "HR Apps",
			"bypassType": "NEVER",
			"bypassOnReauth": false,
			"extranetEnabled": false,
			"matchStyle": "EXCLUSIVE",
			"configSpace": "DEFAULT",
			"domainNames": ["hr.internal.example.com"],
			"enabled": true,
			"passiveHealthEnabled": true,
			"fqdnDnsCheck": true,
			"apiProtectionEnabled": false,
			"selectConnectorCloseToApp": true,
			"doubleEncrypt": false,
			"healthCheckType": "DEFAULT",
			"isCnameEnabled": true,
			"ipAnchored": false,
			"tcpKeepAlive": "1",
			"isIncompleteDRConfig": false,
			"useInDrMode": false,
			"inspectTrafficWithZia": false,
			"healthReporting": "ON_ACCESS",
			"icmpAccessType": "NONE",
			"tcpPortRanges": ["443", "443"],
			"clientlessApps": [
				{
					"id": "ca-002",
					"name": "HR App",
					"appId": "ba-456",
					"applicationPort": "443",
					"applicationProtocol": "HTTPS",
					"certificateId": "cert-002",
					"certificateName": "HR Cert",
					"cname": "hr-portal.zscaler.com",
					"domain": "hr.internal.example.com",
					"enabled": true,
					"hidden": false,
					"localDomain": "hr.internal.example.com",
					"path": "/",
					"trustUntrustedCert": false,
					"allowOptions": true
				}
			],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com"
		}`

		var ba BrowserAccess
		err := json.Unmarshal([]byte(apiResponse), &ba)
		require.NoError(t, err)

		assert.Equal(t, "ba-456", ba.ID)
		assert.Equal(t, "HR Portal", ba.Name)
		assert.Equal(t, "EXCLUSIVE", ba.MatchStyle)
		assert.True(t, ba.Enabled)
		assert.True(t, ba.PassiveHealthEnabled)
		assert.True(t, ba.SelectConnectorCloseToApp)
		assert.Len(t, ba.ClientlessApps, 1)
		assert.Equal(t, "HTTPS", ba.ClientlessApps[0].ApplicationProtocol)
		assert.True(t, ba.ClientlessApps[0].AllowOptions)
	})

	t.Run("ClientlessApps structure", func(t *testing.T) {
		ca := ClientlessApps{
			ID:                  "ca-001",
			Name:                "Test App",
			AppID:               "ba-001",
			ApplicationPort:     "8443",
			ApplicationProtocol: "HTTPS",
			CertificateID:       "cert-001",
			CertificateName:     "Test Cert",
			Cname:               "test.zscaler.com",
			Domain:              "test.internal.example.com",
			Enabled:             true,
			Hidden:              false,
			LocalDomain:         "test.internal.example.com",
			Path:                "/app",
			TrustUntrustedCert:  true,
			AllowOptions:        true,
		}

		data, err := json.Marshal(ca)
		require.NoError(t, err)

		var unmarshaled ClientlessApps
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, ca.ID, unmarshaled.ID)
		assert.Equal(t, ca.ApplicationProtocol, unmarshaled.ApplicationProtocol)
		assert.True(t, unmarshaled.TrustUntrustedCert)
		assert.True(t, unmarshaled.AllowOptions)
	})
}

// TestBrowserAccess_ResponseParsing tests parsing of various API responses
func TestBrowserAccess_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse browser access list response", func(t *testing.T) {
		response := `{
			"list": [
				{
					"id": "1",
					"name": "App 1",
					"enabled": true,
					"clientlessApps": [{"id": "ca-1", "enabled": true}]
				},
				{
					"id": "2",
					"name": "App 2",
					"enabled": true,
					"clientlessApps": [{"id": "ca-2", "enabled": true}]
				},
				{
					"id": "3",
					"name": "App 3 (No BA)",
					"enabled": true,
					"clientlessApps": []
				}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []BrowserAccess `json:"list"`
			TotalPages int             `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Len(t, listResp.List[0].ClientlessApps, 1)
		assert.Len(t, listResp.List[2].ClientlessApps, 0)
	})
}

// TestBrowserAccess_MockServerOperations tests CRUD operations with mock server
func TestBrowserAccess_MockServerOperations(t *testing.T) {
	t.Run("GET browser access by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/application/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "ba-123",
				"name": "Mock Browser Access",
				"enabled": true,
				"clientlessApps": [{"id": "ca-1", "name": "Mock CA"}]
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/application/ba-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all browser access apps", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "BA App A", "clientlessApps": [{"id": "ca-1"}]},
					{"id": "2", "name": "BA App B", "clientlessApps": [{"id": "ca-2"}]}
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

	t.Run("POST create browser access", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-ba-456",
				"name": "New Browser Access",
				"enabled": true,
				"clientlessApps": [{"id": "ca-new"}]
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update browser access", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/application/ba-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE browser access with forceDelete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "true", r.URL.Query().Get("forceDelete"))
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/application/ba-123?forceDelete=true", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestBrowserAccess_SpecialCases tests edge cases
func TestBrowserAccess_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Multiple clientless apps", func(t *testing.T) {
		ba := BrowserAccess{
			ID:      "ba-123",
			Name:    "Multi-App BA",
			Enabled: true,
			ClientlessApps: []ClientlessApps{
				{ID: "ca-1", Name: "App 1", ApplicationPort: "443", ApplicationProtocol: "HTTPS"},
				{ID: "ca-2", Name: "App 2", ApplicationPort: "8443", ApplicationProtocol: "HTTPS"},
				{ID: "ca-3", Name: "App 3", ApplicationPort: "80", ApplicationProtocol: "HTTP"},
			},
		}

		data, err := json.Marshal(ba)
		require.NoError(t, err)

		var unmarshaled BrowserAccess
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.ClientlessApps, 3)
	})

	t.Run("Application protocols", func(t *testing.T) {
		protocols := []string{"HTTP", "HTTPS"}

		for _, protocol := range protocols {
			ca := ClientlessApps{
				ID:                  "ca-" + protocol,
				Name:                protocol + " App",
				ApplicationProtocol: protocol,
			}

			data, err := json.Marshal(ca)
			require.NoError(t, err)

			assert.Contains(t, string(data), protocol)
		}
	})

	t.Run("Trust untrusted cert enabled", func(t *testing.T) {
		ca := ClientlessApps{
			ID:                 "ca-123",
			Name:               "Self-Signed App",
			TrustUntrustedCert: true,
		}

		data, err := json.Marshal(ca)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"trustUntrustedCert":true`)
	})

	t.Run("Hidden clientless app", func(t *testing.T) {
		ca := ClientlessApps{
			ID:     "ca-123",
			Name:   "Hidden App",
			Hidden: true,
		}

		data, err := json.Marshal(ca)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"hidden":true`)
	})

	t.Run("External domain configuration", func(t *testing.T) {
		ca := ClientlessApps{
			ID:            "ca-123",
			Name:          "External App",
			Domain:        "internal.example.com",
			ExtDomain:     "external.zscaler.com",
			ExtDomainName: "external.zscaler.com",
			ExtLabel:      "ext-label",
			ExtID:         "ext-id-123",
		}

		data, err := json.Marshal(ca)
		require.NoError(t, err)

		var unmarshaled ClientlessApps
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "external.zscaler.com", unmarshaled.ExtDomain)
		assert.Equal(t, "ext-id-123", unmarshaled.ExtID)
	})

	t.Run("Shared microtenant details", func(t *testing.T) {
		ba := BrowserAccess{
			ID:   "ba-123",
			Name: "Shared BA",
			SharedMicrotenantDetails: BrowserAccessSharedDetails{
				SharedFromMicrotenant: BrowserAccessSharedFrom{
					ID:   "mt-source",
					Name: "Source Tenant",
				},
				SharedToMicrotenants: []BrowserAccessSharedTo{
					{ID: "mt-dest-1", Name: "Dest 1"},
					{ID: "mt-dest-2", Name: "Dest 2"},
				},
			},
		}

		data, err := json.Marshal(ba)
		require.NoError(t, err)

		var unmarshaled BrowserAccess
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "mt-source", unmarshaled.SharedMicrotenantDetails.SharedFromMicrotenant.ID)
		assert.Len(t, unmarshaled.SharedMicrotenantDetails.SharedToMicrotenants, 2)
	})
}

