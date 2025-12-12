// Package unit provides unit tests for ZPA Browser Protection service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BrowserProtection represents the browser protection profile for testing
type BrowserProtection struct {
	ID                string                      `json:"id,omitempty"`
	Name              string                      `json:"name,omitempty"`
	Description       string                      `json:"description,omitempty"`
	DefaultCSP        bool                        `json:"defaultCSP,omitempty"`
	CriteriaFlagsMask string                      `json:"criteriaFlagsMask,omitempty"`
	Criteria          BrowserProtectionCriteria   `json:"criteria,omitempty"`
	CreationTime      string                      `json:"creationTime,omitempty"`
	ModifiedBy        string                      `json:"modifiedBy,omitempty"`
	ModifiedTime      string                      `json:"modifiedTime,omitempty"`
}

// BrowserProtectionCriteria represents the criteria
type BrowserProtectionCriteria struct {
	FingerPrintCriteria FingerPrintCriteria `json:"fingerPrintCriteria,omitempty"`
}

// FingerPrintCriteria represents fingerprint criteria
type FingerPrintCriteria struct {
	Browser            BrowserCriteria  `json:"browser,omitempty"`
	CollectLocation    bool             `json:"collect_location,omitempty"`
	FingerprintTimeout string           `json:"fingerprint_timeout,omitempty"`
	Location           LocationCriteria `json:"location,omitempty"`
	System             SystemCriteria   `json:"system,omitempty"`
}

// BrowserCriteria represents browser fingerprint criteria
type BrowserCriteria struct {
	BrowserEng     bool `json:"browser_eng,omitempty"`
	BrowserEngVer  bool `json:"browser_eng_ver,omitempty"`
	BrowserName    bool `json:"browser_name,omitempty"`
	BrowserVersion bool `json:"browser_version,omitempty"`
	Canvas         bool `json:"canvas,omitempty"`
	FlashVer       bool `json:"flash_ver,omitempty"`
	FpUsrAgentStr  bool `json:"fp_usr_agent_str,omitempty"`
	IsCookie       bool `json:"is_cookie,omitempty"`
	IsLocalStorage bool `json:"is_local_storage,omitempty"`
	IsSessStorage  bool `json:"is_sess_storage,omitempty"`
	Ja3            bool `json:"ja3,omitempty"`
	Mime           bool `json:"mime,omitempty"`
	Plugin         bool `json:"plugin,omitempty"`
	SilverlightVer bool `json:"silverlight_ver,omitempty"`
}

// LocationCriteria represents location criteria
type LocationCriteria struct {
	Lat bool `json:"lat,omitempty"`
	Lon bool `json:"lon,omitempty"`
}

// SystemCriteria represents system criteria
type SystemCriteria struct {
	AvailScreenResolution bool `json:"avail_screen_resolution,omitempty"`
	CPUArch               bool `json:"cpu_arch,omitempty"`
	CurrScreenResolution  bool `json:"curr_screen_resolution,omitempty"`
	Font                  bool `json:"font,omitempty"`
	JavaVer               bool `json:"java_ver,omitempty"`
	MobileDevType         bool `json:"mobile_dev_type,omitempty"`
	MonitorMobile         bool `json:"monitor_mobile,omitempty"`
	OSName                bool `json:"os_name,omitempty"`
	OSVersion             bool `json:"os_version,omitempty"`
	SysLang               bool `json:"sys_lang,omitempty"`
	Tz                    bool `json:"tz,omitempty"`
	UsrLang               bool `json:"usr_lang,omitempty"`
}

// TestBrowserProtection_Structure tests the struct definitions
func TestBrowserProtection_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BrowserProtection JSON marshaling", func(t *testing.T) {
		bp := BrowserProtection{
			ID:          "bp-123",
			Name:        "Standard Protection",
			Description: "Standard browser protection profile",
			DefaultCSP:  true,
			Criteria: BrowserProtectionCriteria{
				FingerPrintCriteria: FingerPrintCriteria{
					FingerprintTimeout: "30",
					CollectLocation:    true,
					Browser: BrowserCriteria{
						BrowserName:    true,
						BrowserVersion: true,
						Canvas:         true,
						Ja3:            true,
					},
					System: SystemCriteria{
						OSName:    true,
						OSVersion: true,
						CPUArch:   true,
					},
				},
			},
		}

		data, err := json.Marshal(bp)
		require.NoError(t, err)

		var unmarshaled BrowserProtection
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, bp.ID, unmarshaled.ID)
		assert.Equal(t, bp.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.DefaultCSP)
		assert.True(t, unmarshaled.Criteria.FingerPrintCriteria.Browser.Ja3)
	})

	t.Run("BrowserProtection from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "bp-456",
			"name": "Enhanced Protection",
			"description": "Enhanced browser protection",
			"defaultCSP": false,
			"criteriaFlagsMask": "12345",
			"criteria": {
				"fingerPrintCriteria": {
					"fingerprint_timeout": "60",
					"collect_location": true,
					"browser": {
						"browser_name": true,
						"browser_version": true,
						"browser_eng": true,
						"browser_eng_ver": true,
						"canvas": true,
						"ja3": true,
						"fp_usr_agent_str": true,
						"is_cookie": true,
						"is_local_storage": true,
						"is_sess_storage": true
					},
					"location": {
						"lat": true,
						"lon": true
					},
					"system": {
						"os_name": true,
						"os_version": true,
						"cpu_arch": true,
						"avail_screen_resolution": true,
						"curr_screen_resolution": true,
						"tz": true,
						"sys_lang": true
					}
				}
			},
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var bp BrowserProtection
		err := json.Unmarshal([]byte(apiResponse), &bp)
		require.NoError(t, err)

		assert.Equal(t, "bp-456", bp.ID)
		assert.Equal(t, "Enhanced Protection", bp.Name)
		assert.False(t, bp.DefaultCSP)
		assert.True(t, bp.Criteria.FingerPrintCriteria.CollectLocation)
		assert.True(t, bp.Criteria.FingerPrintCriteria.Browser.Ja3)
		assert.True(t, bp.Criteria.FingerPrintCriteria.Location.Lat)
	})
}

// TestBrowserProtection_MockServerOperations tests operations
func TestBrowserProtection_MockServerOperations(t *testing.T) {
	t.Run("GET active browser protection profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/activeBrowserProtectionProfile")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Active Profile", "defaultCSP": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/activeBrowserProtectionProfile")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all browser protection profiles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/browserProtectionProfile")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Profile A", "defaultCSP": true},
					{"id": "2", "name": "Profile B", "defaultCSP": false}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/browserProtectionProfile")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT set active browser protection profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/setActive/")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/browserProtectionProfile/setActive/bp-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestBrowserProtection_SpecialCases tests edge cases
func TestBrowserProtection_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("All browser criteria enabled", func(t *testing.T) {
		criteria := BrowserCriteria{
			BrowserEng:     true,
			BrowserEngVer:  true,
			BrowserName:    true,
			BrowserVersion: true,
			Canvas:         true,
			FlashVer:       true,
			FpUsrAgentStr:  true,
			IsCookie:       true,
			IsLocalStorage: true,
			IsSessStorage:  true,
			Ja3:            true,
			Mime:           true,
			Plugin:         true,
			SilverlightVer: true,
		}

		data, err := json.Marshal(criteria)
		require.NoError(t, err)

		var unmarshaled BrowserCriteria
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.Ja3)
		assert.True(t, unmarshaled.Canvas)
	})

	t.Run("All system criteria enabled", func(t *testing.T) {
		criteria := SystemCriteria{
			AvailScreenResolution: true,
			CPUArch:               true,
			CurrScreenResolution:  true,
			Font:                  true,
			JavaVer:               true,
			MobileDevType:         true,
			MonitorMobile:         true,
			OSName:                true,
			OSVersion:             true,
			SysLang:               true,
			Tz:                    true,
			UsrLang:               true,
		}

		data, err := json.Marshal(criteria)
		require.NoError(t, err)

		var unmarshaled SystemCriteria
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.CPUArch)
		assert.True(t, unmarshaled.Tz)
	})
}

