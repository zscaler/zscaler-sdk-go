// Package zia provides unit tests for ZIA configuration
package zia_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
)

// =====================================================
// Configuration Options Tests
// =====================================================

func TestConfiguration_Options(t *testing.T) {
	t.Parallel()

	t.Run("WithZiaUsername sets username", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithZiaUsername("test-user")(cfg)
		assert.Equal(t, "test-user", cfg.ZIA.Client.ZIAUsername)
	})

	t.Run("WithZiaPassword sets password", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithZiaPassword("test-password")(cfg)
		assert.Equal(t, "test-password", cfg.ZIA.Client.ZIAPassword)
	})

	t.Run("WithZiaAPIKey sets API key", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithZiaAPIKey("test-api-key")(cfg)
		assert.Equal(t, "test-api-key", cfg.ZIA.Client.ZIAApiKey)
	})

	t.Run("WithZiaCloud sets cloud", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithZiaCloud("zscalertwo")(cfg)
		assert.Equal(t, "zscalertwo", cfg.ZIA.Client.ZIACloud)
	})

	t.Run("WithPartnerID sets partner ID", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithPartnerID("partner-123")(cfg)
		assert.Equal(t, "partner-123", cfg.ZIA.Client.PartnerID)
	})

	t.Run("WithDebug sets debug mode", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithDebug(true)(cfg)
		assert.True(t, cfg.Debug)
	})

	t.Run("WithUserAgentExtra sets user agent extra", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithUserAgentExtra("my-custom-agent/1.0")(cfg)
		assert.Contains(t, cfg.UserAgentExtra, "my-custom-agent/1.0")
	})

	t.Run("WithProxyHost sets proxy host", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithProxyHost("proxy.example.com")(cfg)
		assert.Equal(t, "proxy.example.com", cfg.ZIA.Client.Proxy.Host)
	})

	t.Run("WithProxyPort sets proxy port", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithProxyPort(8080)(cfg)
		assert.Equal(t, int32(8080), cfg.ZIA.Client.Proxy.Port)
	})

	t.Run("WithProxyUsername sets proxy username", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithProxyUsername("proxyuser")(cfg)
		assert.Equal(t, "proxyuser", cfg.ZIA.Client.Proxy.Username)
	})

	t.Run("WithProxyPassword sets proxy password", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithProxyPassword("proxypass")(cfg)
		assert.Equal(t, "proxypass", cfg.ZIA.Client.Proxy.Password)
	})

	t.Run("WithTestingDisableHttpsCheck sets HTTPS check", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithTestingDisableHttpsCheck(true)(cfg)
		assert.True(t, cfg.ZIA.Testing.DisableHttpsCheck)
	})

	t.Run("WithRequestTimeout sets request timeout", func(t *testing.T) {
		cfg := &zia.Configuration{}
		timeout := 30 * time.Second
		zia.WithRequestTimeout(timeout)(cfg)
		assert.Equal(t, timeout, cfg.ZIA.Client.RequestTimeout)
	})

	t.Run("WithRateLimitMaxRetries sets max retries", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithRateLimitMaxRetries(10)(cfg)
		assert.Equal(t, int32(10), cfg.ZIA.Client.RateLimit.MaxRetries)
	})

	t.Run("WithRateLimitMinWait sets min wait", func(t *testing.T) {
		cfg := &zia.Configuration{}
		minWait := 5 * time.Second
		zia.WithRateLimitMinWait(minWait)(cfg)
		assert.Equal(t, minWait, cfg.ZIA.Client.RateLimit.RetryWaitMin)
	})

	t.Run("WithRateLimitMaxWait sets max wait", func(t *testing.T) {
		cfg := &zia.Configuration{}
		maxWait := 60 * time.Second
		zia.WithRateLimitMaxWait(maxWait)(cfg)
		assert.Equal(t, maxWait, cfg.ZIA.Client.RateLimit.RetryWaitMax)
	})

	t.Run("WithHttpClientPtr sets custom HTTP client", func(t *testing.T) {
		customClient := &http.Client{Timeout: 120 * time.Second}
		cfg := &zia.Configuration{}
		zia.WithHttpClientPtr(customClient)(cfg)
		assert.NotNil(t, cfg)
	})

	t.Run("WithCache sets cache enabled", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithCache(true)(cfg)
		assert.True(t, cfg.ZIA.Client.Cache.Enabled)
	})

	t.Run("WithCacheTtl sets cache TTL", func(t *testing.T) {
		cfg := &zia.Configuration{}
		ttl := 10 * time.Minute
		zia.WithCacheTtl(ttl)(cfg)
		assert.Equal(t, ttl, cfg.ZIA.Client.Cache.DefaultTtl)
	})

	t.Run("WithCacheMaxSizeMB sets cache max size", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithCacheMaxSizeMB(100)(cfg)
		assert.Equal(t, int64(100), cfg.ZIA.Client.Cache.DefaultCacheMaxSizeMB)
	})

	t.Run("WithCacheTti sets cache TTI", func(t *testing.T) {
		cfg := &zia.Configuration{}
		tti := 5 * time.Minute
		zia.WithCacheTti(tti)(cfg)
		assert.Equal(t, tti, cfg.ZIA.Client.Cache.DefaultTti)
	})
}

// =====================================================
// NewConfiguration Test
// =====================================================

func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("Create configuration with defaults", func(t *testing.T) {
		cfg, err := zia.NewConfiguration(
			zia.WithZiaUsername("test@example.com"),
			zia.WithZiaPassword("password123"),
			zia.WithZiaAPIKey("api-key-123"),
			zia.WithZiaCloud("zscalerone"),
		)
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "test@example.com", cfg.ZIA.Client.ZIAUsername)
		assert.Equal(t, "password123", cfg.ZIA.Client.ZIAPassword)
		assert.Equal(t, "api-key-123", cfg.ZIA.Client.ZIAApiKey)
		assert.Equal(t, "zscalerone", cfg.ZIA.Client.ZIACloud)
	})

	t.Run("Create configuration with no options", func(t *testing.T) {
		cfg, err := zia.NewConfiguration()
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		// Default values should be applied
		assert.NotNil(t, cfg.ZIA)
	})
}

// =====================================================
// Cloud Configuration Tests
// =====================================================

func TestConfiguration_CloudSettings(t *testing.T) {
	t.Parallel()

	testClouds := []string{
		"zscaler",
		"zscalerone",
		"zscalertwo",
		"zscalerthree",
		"zscloud",
		"zscalerbeta",
		"zscalergov",
		"zscalerten",
	}

	for _, cloud := range testClouds {
		cloud := cloud
		t.Run("Set cloud to "+cloud, func(t *testing.T) {
			t.Parallel()
			cfg := &zia.Configuration{}
			zia.WithZiaCloud(cloud)(cfg)
			assert.Equal(t, cloud, cfg.ZIA.Client.ZIACloud)
		})
	}
}

// =====================================================
// Combined Configuration Test
// =====================================================

func TestConfiguration_MultipleOptions(t *testing.T) {
	t.Parallel()

	t.Run("Apply multiple configuration options", func(t *testing.T) {
		cfg := &zia.Configuration{}

		zia.WithZiaUsername("admin@corp.com")(cfg)
		zia.WithZiaPassword("securepass")(cfg)
		zia.WithZiaAPIKey("api-key-456")(cfg)
		zia.WithZiaCloud("zscalertwo")(cfg)
		zia.WithDebug(true)(cfg)
		zia.WithProxyHost("proxy.corp.com")(cfg)
		zia.WithProxyPort(3128)(cfg)
		zia.WithUserAgentExtra("terraform-provider-zscaler/1.0")(cfg)

		assert.Equal(t, "admin@corp.com", cfg.ZIA.Client.ZIAUsername)
		assert.Equal(t, "securepass", cfg.ZIA.Client.ZIAPassword)
		assert.Equal(t, "api-key-456", cfg.ZIA.Client.ZIAApiKey)
		assert.Equal(t, "zscalertwo", cfg.ZIA.Client.ZIACloud)
		assert.True(t, cfg.Debug)
		assert.Equal(t, "proxy.corp.com", cfg.ZIA.Client.Proxy.Host)
		assert.Equal(t, int32(3128), cfg.ZIA.Client.Proxy.Port)
		assert.Contains(t, cfg.UserAgentExtra, "terraform-provider-zscaler/1.0")
	})
}

// =====================================================
// Context Tests
// =====================================================

func TestConfiguration_Context(t *testing.T) {
	t.Parallel()

	t.Run("Configuration with context", func(t *testing.T) {
		cfg := &zia.Configuration{
			Context: context.Background(),
		}
		assert.NotNil(t, cfg.Context)
	})

	t.Run("Configuration with timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cfg := &zia.Configuration{
			Context: ctx,
		}
		assert.NotNil(t, cfg.Context)
	})
}

// =====================================================
// Default Header Tests
// =====================================================

func TestConfiguration_DefaultHeaders(t *testing.T) {
	t.Parallel()

	t.Run("Configuration with default headers", func(t *testing.T) {
		cfg := &zia.Configuration{
			DefaultHeader: map[string]string{
				"X-Custom-Header": "custom-value",
				"Accept":          "application/json",
			},
		}
		assert.Equal(t, "custom-value", cfg.DefaultHeader["X-Custom-Header"])
		assert.Equal(t, "application/json", cfg.DefaultHeader["Accept"])
	})
}

// =====================================================
// Testing Mode Tests
// =====================================================

func TestConfiguration_TestingMode(t *testing.T) {
	t.Parallel()

	t.Run("Enable HTTPS check disable for testing", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithTestingDisableHttpsCheck(true)(cfg)
		assert.True(t, cfg.ZIA.Testing.DisableHttpsCheck)
	})

	t.Run("HTTPS check enabled by default", func(t *testing.T) {
		cfg := &zia.Configuration{}
		assert.False(t, cfg.ZIA.Testing.DisableHttpsCheck)
	})
}

// =====================================================
// User Agent Tests
// =====================================================

func TestConfiguration_UserAgent(t *testing.T) {
	t.Parallel()

	t.Run("Set user agent extra", func(t *testing.T) {
		cfg := &zia.Configuration{}
		zia.WithUserAgentExtra("terraform-provider-zscaler/2.0")(cfg)
		assert.Contains(t, cfg.UserAgentExtra, "terraform-provider-zscaler/2.0")
	})

	t.Run("Set custom user agent", func(t *testing.T) {
		cfg := &zia.Configuration{
			UserAgent: "custom-sdk/1.0",
		}
		assert.Equal(t, "custom-sdk/1.0", cfg.UserAgent)
	})
}
