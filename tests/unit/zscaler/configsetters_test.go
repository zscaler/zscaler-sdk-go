// Package zscaler provides unit tests for OneAPI ConfigSetter functions
package zscaler

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

// =====================================================
// ConfigSetter Tests - Test all With* functions
// =====================================================

func TestConfigSetters_Client(t *testing.T) {
	t.Parallel()

	t.Run("WithClientID", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithClientID("test-client-id")
		setter(cfg)
		assert.Equal(t, "test-client-id", cfg.Zscaler.Client.ClientID)
	})

	t.Run("WithClientSecret", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithClientSecret("test-secret")
		setter(cfg)
		assert.Equal(t, "test-secret", cfg.Zscaler.Client.ClientSecret)
	})

	t.Run("WithVanityDomain", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithVanityDomain("mycompany.zscaler.com")
		setter(cfg)
		assert.Equal(t, "mycompany.zscaler.com", cfg.Zscaler.Client.VanityDomain)
	})

	t.Run("WithZscalerCloud", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithZscalerCloud("zscaler.net")
		setter(cfg)
		assert.Equal(t, "zscaler.net", cfg.Zscaler.Client.Cloud)
	})

	t.Run("WithSandboxToken", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithSandboxToken("sandbox-token-123")
		setter(cfg)
		assert.Equal(t, "sandbox-token-123", cfg.Zscaler.Client.SandboxToken)
	})

	t.Run("WithSandboxCloud", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithSandboxCloud("sandbox.zscaler.net")
		setter(cfg)
		assert.Equal(t, "sandbox.zscaler.net", cfg.Zscaler.Client.SandboxCloud)
	})

	t.Run("WithZPACustomerID", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithZPACustomerID("customer-12345")
		setter(cfg)
		assert.Equal(t, "customer-12345", cfg.Zscaler.Client.CustomerID)
	})

	t.Run("WithZPAMicrotenantID", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithZPAMicrotenantID("microtenant-67890")
		setter(cfg)
		assert.Equal(t, "microtenant-67890", cfg.Zscaler.Client.MicrotenantID)
	})

	t.Run("WithPartnerID", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithPartnerID("partner-abc")
		setter(cfg)
		assert.Equal(t, "partner-abc", cfg.Zscaler.Client.PartnerID)
	})
}

func TestConfigSetters_Cache(t *testing.T) {
	t.Parallel()

	t.Run("WithCache enabled", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithCache(true)
		setter(cfg)
		assert.True(t, cfg.Zscaler.Client.Cache.Enabled)
	})

	t.Run("WithCache disabled", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithCache(false)
		setter(cfg)
		assert.False(t, cfg.Zscaler.Client.Cache.Enabled)
	})

	t.Run("WithCacheTtl", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithCacheTtl(10 * time.Minute)
		setter(cfg)
		assert.Equal(t, 10*time.Minute, cfg.Zscaler.Client.Cache.DefaultTtl)
	})

	t.Run("WithCacheMaxSizeMB", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithCacheMaxSizeMB(500)
		setter(cfg)
		assert.Equal(t, int64(500), cfg.Zscaler.Client.Cache.DefaultCacheMaxSizeMB)
	})

	t.Run("WithCacheTti", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithCacheTti(5 * time.Minute)
		setter(cfg)
		assert.Equal(t, 5*time.Minute, cfg.Zscaler.Client.Cache.DefaultTti)
	})
}

func TestConfigSetters_HTTP(t *testing.T) {
	t.Parallel()

	t.Run("WithHttpClientPtr", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		httpClient := &http.Client{Timeout: 30 * time.Second}
		setter := zscaler.WithHttpClientPtr(httpClient)
		setter(cfg)
		require.NotNil(t, cfg.HTTPClient)
		assert.Equal(t, 30*time.Second, cfg.HTTPClient.Timeout)
	})

	t.Run("WithRequestTimeout", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithRequestTimeout(2 * time.Minute)
		setter(cfg)
		assert.Equal(t, 2*time.Minute, cfg.Zscaler.Client.RequestTimeout)
	})

	t.Run("WithTestingDisableHttpsCheck", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithTestingDisableHttpsCheck(true)
		setter(cfg)
		assert.True(t, cfg.Zscaler.Testing.DisableHttpsCheck)
	})
}

func TestConfigSetters_Proxy(t *testing.T) {
	t.Parallel()

	t.Run("WithProxyHost", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithProxyHost("proxy.example.com")
		setter(cfg)
		assert.Equal(t, "proxy.example.com", cfg.Zscaler.Client.Proxy.Host)
	})

	t.Run("WithProxyPort", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithProxyPort(8080)
		setter(cfg)
		assert.Equal(t, int32(8080), cfg.Zscaler.Client.Proxy.Port)
	})

	t.Run("WithProxyUsername", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithProxyUsername("proxyuser")
		setter(cfg)
		assert.Equal(t, "proxyuser", cfg.Zscaler.Client.Proxy.Username)
	})

	t.Run("WithProxyPassword", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithProxyPassword("proxypass")
		setter(cfg)
		assert.Equal(t, "proxypass", cfg.Zscaler.Client.Proxy.Password)
	})
}

func TestConfigSetters_RateLimit(t *testing.T) {
	t.Parallel()

	t.Run("WithRateLimitMaxRetries", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithRateLimitMaxRetries(10)
		setter(cfg)
		assert.Equal(t, int32(10), cfg.Zscaler.Client.RateLimit.MaxRetries)
	})

	t.Run("WithRateLimitMaxWait", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithRateLimitMaxWait(30 * time.Second)
		setter(cfg)
		assert.Equal(t, 30*time.Second, cfg.Zscaler.Client.RateLimit.RetryWaitMax)
	})

	t.Run("WithRateLimitMinWait", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithRateLimitMinWait(5 * time.Second)
		setter(cfg)
		assert.Equal(t, 5*time.Second, cfg.Zscaler.Client.RateLimit.RetryWaitMin)
	})

	t.Run("WithRateLimitRemainingThreshold", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithRateLimitRemainingThreshold(3)
		setter(cfg)
		assert.Equal(t, int32(3), cfg.Zscaler.Client.RateLimit.RetryRemainingThreshold)
	})

	t.Run("WithRateLimitMaxSessionNotValidRetries", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithRateLimitMaxSessionNotValidRetries(5)
		setter(cfg)
		assert.Equal(t, int32(5), cfg.Zscaler.Client.RateLimit.MaxSessionNotValidRetries)
	})
}

func TestConfigSetters_Misc(t *testing.T) {
	t.Parallel()

	t.Run("WithUserAgentExtra", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithUserAgentExtra("terraform-provider/1.0.0")
		setter(cfg)
		assert.Equal(t, "terraform-provider/1.0.0", cfg.UserAgentExtra)
	})

	t.Run("WithDebug enabled", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithDebug(true)
		setter(cfg)
		assert.True(t, cfg.Debug)
	})

	t.Run("WithDebug disabled", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithDebug(false)
		setter(cfg)
		assert.False(t, cfg.Debug)
	})

	t.Run("WithLegacyClient", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithLegacyClient(true)
		setter(cfg)
		assert.True(t, cfg.UseLegacyClient)
	})
}

func TestConfiguration_AddDefaultHeader(t *testing.T) {
	t.Parallel()

	t.Run("Add single header", func(t *testing.T) {
		cfg := &zscaler.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Custom-Header", "custom-value")
		assert.Equal(t, "custom-value", cfg.DefaultHeader["X-Custom-Header"])
	})

	t.Run("Add multiple headers", func(t *testing.T) {
		cfg := &zscaler.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("Header1", "Value1")
		cfg.AddDefaultHeader("Header2", "Value2")
		assert.Equal(t, "Value1", cfg.DefaultHeader["Header1"])
		assert.Equal(t, "Value2", cfg.DefaultHeader["Header2"])
	})

	t.Run("Override existing header", func(t *testing.T) {
		cfg := &zscaler.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Header", "old-value")
		cfg.AddDefaultHeader("X-Header", "new-value")
		assert.Equal(t, "new-value", cfg.DefaultHeader["X-Header"])
	})
}

