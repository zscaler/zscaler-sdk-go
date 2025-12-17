// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SecurityPolicySettings represents security policy configuration
type SecurityPolicySettings struct {
	WhitelistUrls   []string `json:"whitelistUrls,omitempty"`
	BlacklistUrls   []string `json:"blacklistUrls,omitempty"`
	EnableDynamicContentAnalysis bool `json:"enableDynamicContentAnalysis,omitempty"`
	EnablePatientZeroAlert       bool `json:"enablePatientZeroAlert,omitempty"`
	EnableSSLScanning            bool `json:"enableSSLScanning,omitempty"`
}

// AdvancedThreatSettings represents advanced threat protection settings
type AdvancedThreatSettings struct {
	EnableInfectedContentDeliveryProtection bool   `json:"enableInfectedContentDeliveryProtection,omitempty"`
	EnableInfectedContentProtection         bool   `json:"enableInfectedContentProtection,omitempty"`
	EnableMalwareProtection                 bool   `json:"enableMalwareProtection,omitempty"`
	EnableAdvancedThreatProtection          bool   `json:"enableAdvancedThreatProtection,omitempty"`
	MalwareThreatAction                     string `json:"malwareThreatAction,omitempty"`
	VirusProtectionAction                   string `json:"virusProtectionAction,omitempty"`
}

// SSLInspectionSettings represents SSL inspection configuration
type SSLInspectionSettings struct {
	InterceptSSL                bool     `json:"interceptSsl,omitempty"`
	EnableSSLForCloudApps       bool     `json:"enableSslForCloudApps,omitempty"`
	BlockSSLTraficWithUntrustedCert bool `json:"blockSslTraficWithUntrustedCert,omitempty"`
	SSLCertValidation           string   `json:"sslCertValidation,omitempty"`
	BypassUrls                  []string `json:"bypassUrls,omitempty"`
	BypassCategories            []string `json:"bypassCategories,omitempty"`
}

// SandboxSettings represents sandbox configuration
type SandboxSettings struct {
	FileHashCount       int      `json:"fileHashCount,omitempty"`
	BehavioralAnalysis  bool     `json:"behavioralAnalysis,omitempty"`
	FileTypesForAnalysis []string `json:"fileTypesForAnalysis,omitempty"`
	UrlCategoriesToScan []string `json:"urlCategoriesToScan,omitempty"`
}

func TestSecurityPolicySettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("SecurityPolicySettings JSON marshaling", func(t *testing.T) {
		settings := SecurityPolicySettings{
			WhitelistUrls: []string{"trusted.com", "*.internal.com"},
			BlacklistUrls: []string{"malware.com", "phishing.com"},
			EnableDynamicContentAnalysis: true,
			EnablePatientZeroAlert:       true,
			EnableSSLScanning:            true,
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"whitelistUrls"`)
		assert.Contains(t, string(data), `"enableDynamicContentAnalysis":true`)
	})

	t.Run("SecurityPolicySettings JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"whitelistUrls": ["safe1.com", "safe2.com"],
			"blacklistUrls": ["bad1.com"],
			"enableDynamicContentAnalysis": true,
			"enablePatientZeroAlert": false,
			"enableSSLScanning": true
		}`

		var settings SecurityPolicySettings
		err := json.Unmarshal([]byte(jsonData), &settings)
		require.NoError(t, err)

		assert.Len(t, settings.WhitelistUrls, 2)
		assert.Len(t, settings.BlacklistUrls, 1)
		assert.True(t, settings.EnableDynamicContentAnalysis)
	})
}

func TestAdvancedThreatSettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdvancedThreatSettings JSON marshaling", func(t *testing.T) {
		settings := AdvancedThreatSettings{
			EnableInfectedContentDeliveryProtection: true,
			EnableInfectedContentProtection:         true,
			EnableMalwareProtection:                 true,
			EnableAdvancedThreatProtection:          true,
			MalwareThreatAction:                     "BLOCK",
			VirusProtectionAction:                   "QUARANTINE",
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enableMalwareProtection":true`)
		assert.Contains(t, string(data), `"malwareThreatAction":"BLOCK"`)
	})

	t.Run("AdvancedThreatSettings JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"enableInfectedContentDeliveryProtection": true,
			"enableInfectedContentProtection": true,
			"enableMalwareProtection": false,
			"enableAdvancedThreatProtection": true,
			"malwareThreatAction": "ALLOW",
			"virusProtectionAction": "BLOCK"
		}`

		var settings AdvancedThreatSettings
		err := json.Unmarshal([]byte(jsonData), &settings)
		require.NoError(t, err)

		assert.False(t, settings.EnableMalwareProtection)
		assert.Equal(t, "ALLOW", settings.MalwareThreatAction)
	})
}

func TestSSLInspectionSettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("SSLInspectionSettings JSON marshaling", func(t *testing.T) {
		settings := SSLInspectionSettings{
			InterceptSSL:                    true,
			EnableSSLForCloudApps:           true,
			BlockSSLTraficWithUntrustedCert: true,
			SSLCertValidation:               "STRICT",
			BypassUrls:                      []string{"banking.com", "healthcare.org"},
			BypassCategories:                []string{"FINANCIAL_SERVICES", "HEALTHCARE"},
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"interceptSsl":true`)
		assert.Contains(t, string(data), `"sslCertValidation":"STRICT"`)
	})

	t.Run("SSLInspectionSettings JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"interceptSsl": true,
			"enableSslForCloudApps": false,
			"blockSslTraficWithUntrustedCert": true,
			"sslCertValidation": "RELAXED",
			"bypassUrls": ["bypass1.com"],
			"bypassCategories": ["EDUCATION"]
		}`

		var settings SSLInspectionSettings
		err := json.Unmarshal([]byte(jsonData), &settings)
		require.NoError(t, err)

		assert.True(t, settings.InterceptSSL)
		assert.False(t, settings.EnableSSLForCloudApps)
		assert.Equal(t, "RELAXED", settings.SSLCertValidation)
	})
}

func TestSandboxSettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("SandboxSettings JSON marshaling", func(t *testing.T) {
		settings := SandboxSettings{
			FileHashCount:       1000,
			BehavioralAnalysis:  true,
			FileTypesForAnalysis: []string{"EXE", "DLL", "PDF", "DOC"},
			UrlCategoriesToScan: []string{"MALWARE", "PHISHING"},
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"fileHashCount":1000`)
		assert.Contains(t, string(data), `"behavioralAnalysis":true`)
	})

	t.Run("SandboxSettings JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"fileHashCount": 500,
			"behavioralAnalysis": false,
			"fileTypesForAnalysis": ["EXE", "MSI"],
			"urlCategoriesToScan": ["UNCATEGORIZED"]
		}`

		var settings SandboxSettings
		err := json.Unmarshal([]byte(jsonData), &settings)
		require.NoError(t, err)

		assert.Equal(t, 500, settings.FileHashCount)
		assert.False(t, settings.BehavioralAnalysis)
		assert.Len(t, settings.FileTypesForAnalysis, 2)
	})
}

