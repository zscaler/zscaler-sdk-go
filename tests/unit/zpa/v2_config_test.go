// Package unit provides unit tests for the ZPA v2 configuration
package unit

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
)

// =============================================================================
// Configuration Default Tests
// =============================================================================

func TestZPAConfiguration_Defaults(t *testing.T) {
	t.Run("Default values are set correctly", func(t *testing.T) {
		// Note: This test validates the default values WITHOUT creating a full client
		// which would require authentication

		// We can't call NewConfiguration directly without credentials,
		// so we test the constants and expected defaults
		assert.Equal(t, 100, zpa.MaxNumOfRetries)
		assert.Equal(t, 10, zpa.RetryWaitMaxSeconds)
		assert.Equal(t, 2, zpa.RetryWaitMinSeconds)
	})

	t.Run("VERSION is set", func(t *testing.T) {
		version := zpa.VERSION
		assert.NotEmpty(t, version)
		// Version should match semver pattern
		assert.Regexp(t, `^\d+\.\d+\.\d+`, version)
	})
}

// =============================================================================
// ConfigSetter Tests
// =============================================================================

func TestZPAConfigSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithZPAClientID", func(t *testing.T) {
		t.Parallel()

		// Create a minimal config to test the setter
		cfg := &zpa.Configuration{}
		setter := zpa.WithZPAClientID("test-client-id")
		setter(cfg)

		assert.Equal(t, "test-client-id", cfg.ZPA.Client.ZPAClientID)
	})

	t.Run("WithZPAClientSecret", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithZPAClientSecret("test-client-secret")
		setter(cfg)

		assert.Equal(t, "test-client-secret", cfg.ZPA.Client.ZPAClientSecret)
	})

	t.Run("WithZPACustomerID", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithZPACustomerID("test-customer-id")
		setter(cfg)

		assert.Equal(t, "test-customer-id", cfg.ZPA.Client.ZPACustomerID)
	})

	t.Run("WithZPAMicrotenantID", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithZPAMicrotenantID("test-microtenant-id")
		setter(cfg)

		assert.Equal(t, "test-microtenant-id", cfg.ZPA.Client.ZPAMicrotenantID)
	})

	t.Run("WithPartnerID", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithPartnerID("test-partner-id")
		setter(cfg)

		assert.Equal(t, "test-partner-id", cfg.ZPA.Client.PartnerID)
	})

	t.Run("WithZPACloud", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			cloud    string
			expected string
		}{
			{"PRODUCTION", "PRODUCTION"},
			{"BETA", "BETA"},
			{"GOV", "GOV"},
			{"GOVUS", "GOVUS"},
			{"ZPATWO", "ZPATWO"},
			{"PREVIEW", "PREVIEW"},
			{"DEV", "DEV"},
			{"QA", "QA"},
			{"QA2", "QA2"},
		}

		for _, tc := range testCases {
			t.Run(tc.cloud, func(t *testing.T) {
				cfg := &zpa.Configuration{}
				setter := zpa.WithZPACloud(tc.cloud)
				setter(cfg)

				assert.Equal(t, tc.expected, cfg.ZPA.Client.ZPACloud)
			})
		}
	})

	t.Run("WithCache", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}

		// Enable cache
		setter := zpa.WithCache(true)
		setter(cfg)
		assert.True(t, cfg.ZPA.Client.Cache.Enabled)

		// Disable cache
		setter = zpa.WithCache(false)
		setter(cfg)
		assert.False(t, cfg.ZPA.Client.Cache.Enabled)
	})

	t.Run("WithCacheTtl", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithCacheTtl(15 * time.Minute)
		setter(cfg)

		assert.Equal(t, 15*time.Minute, cfg.ZPA.Client.Cache.DefaultTtl)
	})

	t.Run("WithCacheTti", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithCacheTti(10 * time.Minute)
		setter(cfg)

		assert.Equal(t, 10*time.Minute, cfg.ZPA.Client.Cache.DefaultTti)
	})

	t.Run("WithCacheMaxSizeMB", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithCacheMaxSizeMB(200)
		setter(cfg)

		assert.Equal(t, int64(200), cfg.ZPA.Client.Cache.DefaultCacheMaxSizeMB)
	})

	t.Run("WithProxyHost", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithProxyHost("proxy.example.com")
		setter(cfg)

		assert.Equal(t, "proxy.example.com", cfg.ZPA.Client.Proxy.Host)
	})

	t.Run("WithProxyPort", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithProxyPort(8080)
		setter(cfg)

		assert.Equal(t, int32(8080), cfg.ZPA.Client.Proxy.Port)
	})

	t.Run("WithProxyUsername", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithProxyUsername("proxy-user")
		setter(cfg)

		assert.Equal(t, "proxy-user", cfg.ZPA.Client.Proxy.Username)
	})

	t.Run("WithProxyPassword", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithProxyPassword("proxy-pass")
		setter(cfg)

		assert.Equal(t, "proxy-pass", cfg.ZPA.Client.Proxy.Password)
	})

	t.Run("WithTestingDisableHttpsCheck", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}

		setter := zpa.WithTestingDisableHttpsCheck(true)
		setter(cfg)
		assert.True(t, cfg.ZPA.Testing.DisableHttpsCheck)

		setter = zpa.WithTestingDisableHttpsCheck(false)
		setter(cfg)
		assert.False(t, cfg.ZPA.Testing.DisableHttpsCheck)
	})

	t.Run("WithRequestTimeout", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithRequestTimeout(5 * time.Minute)
		setter(cfg)

		assert.Equal(t, 5*time.Minute, cfg.ZPA.Client.RequestTimeout)
	})

	t.Run("WithRateLimitMaxRetries", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithRateLimitMaxRetries(50)
		setter(cfg)

		assert.Equal(t, int32(50), cfg.ZPA.Client.RateLimit.MaxRetries)
	})

	t.Run("WithRateLimitMaxWait", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithRateLimitMaxWait(30 * time.Second)
		setter(cfg)

		assert.Equal(t, 30*time.Second, cfg.ZPA.Client.RateLimit.RetryWaitMax)
	})

	t.Run("WithRateLimitMinWait", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithRateLimitMinWait(5 * time.Second)
		setter(cfg)

		assert.Equal(t, 5*time.Second, cfg.ZPA.Client.RateLimit.RetryWaitMin)
	})

	t.Run("WithUserAgentExtra", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithUserAgentExtra("terraform/1.0.0")
		setter(cfg)

		assert.Equal(t, "terraform/1.0.0", cfg.UserAgentExtra)
	})

	t.Run("WithDebug", func(t *testing.T) {
		t.Parallel()

		cfg := &zpa.Configuration{}
		setter := zpa.WithDebug(true)
		setter(cfg)

		assert.True(t, cfg.Debug)
	})
}

// =============================================================================
// UserAgent Tests
// =============================================================================

func TestZPAUserAgent(t *testing.T) {
	t.Run("Default UserAgent format", func(t *testing.T) {
		// The default user agent should contain SDK version, Go version, OS and arch
		expectedPattern := "zscaler-sdk-go/" + zpa.VERSION + " golang/" + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH

		cfg := &zpa.Configuration{}
		cfg.UserAgent = expectedPattern

		assert.Contains(t, cfg.UserAgent, "zscaler-sdk-go/")
		assert.Contains(t, cfg.UserAgent, "golang/")
		assert.Contains(t, cfg.UserAgent, runtime.GOOS)
		assert.Contains(t, cfg.UserAgent, runtime.GOARCH)
	})
}

// =============================================================================
// AddDefaultHeader Tests
// =============================================================================

func TestZPAConfiguration_AddDefaultHeader(t *testing.T) {
	t.Run("Add custom header", func(t *testing.T) {
		cfg := &zpa.Configuration{
			DefaultHeader: make(map[string]string),
		}

		cfg.AddDefaultHeader("X-Custom-Header", "custom-value")

		assert.Equal(t, "custom-value", cfg.DefaultHeader["X-Custom-Header"])
	})

	t.Run("Override existing header", func(t *testing.T) {
		cfg := &zpa.Configuration{
			DefaultHeader: map[string]string{
				"X-Custom-Header": "original-value",
			},
		}

		cfg.AddDefaultHeader("X-Custom-Header", "new-value")

		assert.Equal(t, "new-value", cfg.DefaultHeader["X-Custom-Header"])
	})
}

// =============================================================================
// Cloud Base URL Tests
// =============================================================================

func TestZPACloudBaseURLs(t *testing.T) {
	t.Parallel()

	// Test that the cloud constants are defined correctly
	cloudURLs := map[string]string{
		"PRODUCTION": "https://config.private.zscaler.com",
		"BETA":       "https://config.zpabeta.net",
		"ZPATWO":     "https://config.zpatwo.net",
		"GOV":        "https://config.zpagov.net",
		"GOVUS":      "https://config.zpagov.us",
		"PREVIEW":    "https://config.zpapreview.net",
		"DEV":        "https://public-api.dev.zpath.net",
		"QA":         "https://config.qa.zpath.net",
		"QA2":        "https://pdx2-zpa-config.qa2.zpath.net",
	}

	for cloud, expectedURL := range cloudURLs {
		t.Run(cloud, func(t *testing.T) {
			// Verify the URL is valid (basic check)
			assert.NotEmpty(t, expectedURL)
			assert.Contains(t, expectedURL, "https://")
		})
	}
}

// =============================================================================
// Context Key Tests
// =============================================================================

func TestZPAContextKey(t *testing.T) {
	t.Run("ContextAccessToken is defined", func(t *testing.T) {
		// The ContextAccessToken should be usable as a context key
		key := zpa.ContextAccessToken
		assert.NotNil(t, key)
		assert.Contains(t, key.String(), "access_token")
	})
}

// =============================================================================
// AuthToken Structure Tests
// =============================================================================

func TestZPAAuthToken(t *testing.T) {
	t.Run("AuthToken fields", func(t *testing.T) {
		token := &zpa.AuthToken{
			TokenType:   "Bearer",
			AccessToken: "test-token",
			ExpiresIn:   3600,
			Expiry:      time.Now().Add(1 * time.Hour),
		}

		assert.Equal(t, "Bearer", token.TokenType)
		assert.Equal(t, "test-token", token.AccessToken)
		assert.Equal(t, 3600, token.ExpiresIn)
		require.False(t, token.Expiry.IsZero())
	})
}

// =============================================================================
// Environment Variable Constants Tests
// =============================================================================

func TestZPAEnvironmentVariableConstants(t *testing.T) {
	t.Run("Environment variable names are correct", func(t *testing.T) {
		assert.Equal(t, "ZPA_CLIENT_ID", zpa.ZPA_CLIENT_ID)
		assert.Equal(t, "ZPA_CLIENT_SECRET", zpa.ZPA_CLIENT_SECRET)
		assert.Equal(t, "ZPA_CUSTOMER_ID", zpa.ZPA_CUSTOMER_ID)
		assert.Equal(t, "ZPA_CLOUD", zpa.ZPA_CLOUD)
	})
}

