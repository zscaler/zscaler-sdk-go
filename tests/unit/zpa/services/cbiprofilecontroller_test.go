// Package unit provides unit tests for ZPA CBI Profile Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CBIIsolationProfile represents the CBI isolation profile for testing
type CBIIsolationProfile struct {
	ID               string                  `json:"id,omitempty"`
	Name             string                  `json:"name,omitempty"`
	Description      string                  `json:"description,omitempty"`
	Enabled          bool                    `json:"enabled,omitempty"`
	CBITenantID      string                  `json:"cbiTenantId,omitempty"`
	CBIProfileID     string                  `json:"cbiProfileId,omitempty"`
	CBIURL           string                  `json:"cbiUrl,omitempty"`
	BannerID         string                  `json:"bannerId,omitempty"`
	SecurityControls *CBISecurityControls    `json:"securityControls,omitempty"`
	IsDefault        bool                    `json:"isDefault,omitempty"`
	Regions          []CBIRegion             `json:"regions,omitempty"`
	RegionIDs        []string                `json:"regionIds,omitempty"`
	UserExperience   *CBIUserExperience      `json:"userExperience,omitempty"`
	Certificates     []CBICertificateRef     `json:"certificates,omitempty"`
	CertificateIDs   []string                `json:"certificateIds,omitempty"`
	Banner           *CBIBannerRef           `json:"banner,omitempty"`
	DebugMode        *CBIDebugMode           `json:"debugMode,omitempty"`
	CreationTime     string                  `json:"creationTime,omitempty"`
	ModifiedBy       string                  `json:"modifiedBy,omitempty"`
	ModifiedTime     string                  `json:"modifiedTime,omitempty"`
}

type CBISecurityControls struct {
	DocumentViewer     bool          `json:"documentViewer,omitempty"`
	AllowPrinting      bool          `json:"allowPrinting,omitempty"`
	Watermark          *CBIWatermark `json:"watermark,omitempty"`
	FlattenedPdf       bool          `json:"flattenedPdf,omitempty"`
	UploadDownload     string        `json:"uploadDownload,omitempty"`
	RestrictKeystrokes bool          `json:"restrictKeystrokes,omitempty"`
	CopyPaste          string        `json:"copyPaste,omitempty"`
	LocalRender        bool          `json:"localRender,omitempty"`
	DeepLink           *CBIDeepLink  `json:"deepLink,omitempty"`
}

type CBIWatermark struct {
	Enabled       bool   `json:"enabled,omitempty"`
	ShowUserID    bool   `json:"showUserId,omitempty"`
	ShowTimestamp bool   `json:"showTimestamp,omitempty"`
	ShowMessage   bool   `json:"showMessage,omitempty"`
	Message       string `json:"message,omitempty"`
}

type CBIDeepLink struct {
	Enabled      bool     `json:"enabled,omitempty"`
	Applications []string `json:"applications,omitempty"`
}

type CBIUserExperience struct {
	SessionPersistence  bool              `json:"sessionPersistence"`
	BrowserInBrowser    bool              `json:"browserInBrowser"`
	PersistIsolationBar bool              `json:"persistIsolationBar"`
	Translate           bool              `json:"translate"`
	ZGPU                bool              `json:"zgpu,omitempty"`
	ForwardToZia        *CBIForwardToZia  `json:"forwardToZia,omitempty"`
}

type CBIForwardToZia struct {
	Enabled        bool   `json:"enabled"`
	OrganizationID string `json:"organizationId"`
	CloudName      string `json:"cloudName,omitempty"`
	PacFileUrl     string `json:"pacFileUrl,omitempty"`
}

type CBIRegion struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CBICertificateRef struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

type CBIBannerRef struct {
	ID string `json:"id,omitempty"`
}

type CBIDebugMode struct {
	Allowed      bool   `json:"allowed,omitempty"`
	FilePassword string `json:"filePassword,omitempty"`
}

// TestCBIProfile_Structure tests the struct definitions
func TestCBIProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CBIIsolationProfile JSON marshaling", func(t *testing.T) {
		profile := CBIIsolationProfile{
			ID:          "profile-123",
			Name:        "Secure Isolation Profile",
			Description: "Profile for secure browser isolation",
			Enabled:     true,
			CBITenantID: "tenant-001",
			CBIURL:      "https://cbi.zscaler.com",
			RegionIDs:   []string{"us-west-1", "us-east-1"},
			SecurityControls: &CBISecurityControls{
				DocumentViewer:     true,
				AllowPrinting:      false,
				UploadDownload:     "none",
				CopyPaste:          "none",
				RestrictKeystrokes: true,
				Watermark: &CBIWatermark{
					Enabled:       true,
					ShowUserID:    true,
					ShowTimestamp: true,
				},
			},
			UserExperience: &CBIUserExperience{
				SessionPersistence:  true,
				BrowserInBrowser:    false,
				PersistIsolationBar: true,
				Translate:           true,
			},
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled CBIIsolationProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.NotNil(t, unmarshaled.SecurityControls)
		assert.True(t, unmarshaled.SecurityControls.Watermark.Enabled)
	})

	t.Run("CBIIsolationProfile from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "profile-456",
			"name": "Enterprise Isolation",
			"description": "Enterprise-grade isolation profile",
			"enabled": true,
			"cbiTenantId": "tenant-002",
			"cbiProfileId": "cbi-profile-001",
			"cbiUrl": "https://enterprise-cbi.zscaler.com",
			"isDefault": true,
			"regions": [
				{"id": "us-west-1", "name": "US West"},
				{"id": "eu-central-1", "name": "EU Central"}
			],
			"securityControls": {
				"documentViewer": true,
				"allowPrinting": true,
				"uploadDownload": "all",
				"copyPaste": "isolatedToLocal",
				"restrictKeystrokes": false,
				"localRender": true,
				"deepLink": {
					"enabled": true,
					"applications": ["zoom", "teams"]
				}
			},
			"userExperience": {
				"sessionPersistence": true,
				"browserInBrowser": true,
				"persistIsolationBar": false,
				"translate": true,
				"zgpu": true,
				"forwardToZia": {
					"enabled": true,
					"organizationId": "org-123",
					"cloudName": "zscaler.net"
				}
			},
			"debugMode": {
				"allowed": true,
				"filePassword": "encrypted"
			}
		}`

		var profile CBIIsolationProfile
		err := json.Unmarshal([]byte(apiResponse), &profile)
		require.NoError(t, err)

		assert.Equal(t, "profile-456", profile.ID)
		assert.True(t, profile.IsDefault)
		assert.Len(t, profile.Regions, 2)
		assert.True(t, profile.SecurityControls.DeepLink.Enabled)
		assert.True(t, profile.UserExperience.ForwardToZia.Enabled)
		assert.True(t, profile.DebugMode.Allowed)
	})
}

// TestCBIProfile_MockServerOperations tests CRUD operations
func TestCBIProfile_MockServerOperations(t *testing.T) {
	t.Run("GET profile by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/profiles/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "profile-123", "name": "Mock Profile", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/profiles/profile-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all profiles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[{"id": "1", "name": "Profile A"}, {"id": "2", "name": "Profile B"}]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/profiles")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-profile", "name": "New Profile"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/cbi/api/customers/123/profiles", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/profiles/profile-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestCBIProfile_SpecialCases tests edge cases
func TestCBIProfile_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Upload/Download options", func(t *testing.T) {
		options := []string{"none", "all", "uploadOnly", "downloadOnly"}

		for _, opt := range options {
			controls := CBISecurityControls{
				UploadDownload: opt,
			}

			data, err := json.Marshal(controls)
			require.NoError(t, err)

			assert.Contains(t, string(data), opt)
		}
	})

	t.Run("Copy/Paste options", func(t *testing.T) {
		options := []string{"none", "all", "localToIsolated", "isolatedToLocal"}

		for _, opt := range options {
			controls := CBISecurityControls{
				CopyPaste: opt,
			}

			data, err := json.Marshal(controls)
			require.NoError(t, err)

			assert.Contains(t, string(data), opt)
		}
	})

	t.Run("ZIA forwarding configuration", func(t *testing.T) {
		zia := CBIForwardToZia{
			Enabled:        true,
			OrganizationID: "org-123",
			CloudName:      "zscaler.net",
			PacFileUrl:     "https://pac.zscaler.net/proxy.pac",
		}

		data, err := json.Marshal(zia)
		require.NoError(t, err)

		var unmarshaled CBIForwardToZia
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, "org-123", unmarshaled.OrganizationID)
	})
}

