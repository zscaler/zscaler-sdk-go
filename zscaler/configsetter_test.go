package zscaler

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigSetters(t *testing.T) {
	t.Run("WithClientID", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithClientID("test-client-id")
		setter(cfg)
		assert.Equal(t, "test-client-id", cfg.Zscaler.Client.ClientID)
	})

	t.Run("WithClientSecret", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithClientSecret("test-client-secret")
		setter(cfg)
		assert.Equal(t, "test-client-secret", cfg.Zscaler.Client.ClientSecret)
	})

	t.Run("WithVanityDomain", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithVanityDomain("mycompany.zscaler.com")
		setter(cfg)
		assert.Equal(t, "mycompany.zscaler.com", cfg.Zscaler.Client.VanityDomain)
	})

	t.Run("WithZscalerCloud", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithZscalerCloud("zscaler.net")
		setter(cfg)
		assert.Equal(t, "zscaler.net", cfg.Zscaler.Client.Cloud)
	})

	t.Run("WithSandboxToken", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithSandboxToken("sandbox-token-123")
		setter(cfg)
		assert.Equal(t, "sandbox-token-123", cfg.Zscaler.Client.SandboxToken)
	})

	t.Run("WithSandboxCloud", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithSandboxCloud("sandbox.zscaler.net")
		setter(cfg)
		assert.Equal(t, "sandbox.zscaler.net", cfg.Zscaler.Client.SandboxCloud)
	})

	t.Run("WithZPACustomerID", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithZPACustomerID("customer-12345")
		setter(cfg)
		assert.Equal(t, "customer-12345", cfg.Zscaler.Client.CustomerID)
	})

	t.Run("WithZPAMicrotenantID", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithZPAMicrotenantID("microtenant-67890")
		setter(cfg)
		assert.Equal(t, "microtenant-67890", cfg.Zscaler.Client.MicrotenantID)
	})

	t.Run("WithPartnerID", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithPartnerID("partner-abc")
		setter(cfg)
		assert.Equal(t, "partner-abc", cfg.Zscaler.Client.PartnerID)
	})

	t.Run("WithCache", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithCache(true)
		setter(cfg)
		assert.True(t, cfg.Zscaler.Client.Cache.Enabled)
	})

	t.Run("WithCacheTtl", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithCacheTtl(10 * time.Minute)
		setter(cfg)
		assert.Equal(t, 10*time.Minute, cfg.Zscaler.Client.Cache.DefaultTtl)
	})

	t.Run("WithCacheMaxSizeMB", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithCacheMaxSizeMB(500)
		setter(cfg)
		assert.Equal(t, int64(500), cfg.Zscaler.Client.Cache.DefaultCacheMaxSizeMB)
	})

	t.Run("WithCacheTti", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithCacheTti(5 * time.Minute)
		setter(cfg)
		assert.Equal(t, 5*time.Minute, cfg.Zscaler.Client.Cache.DefaultTti)
	})

	t.Run("WithHttpClientPtr", func(t *testing.T) {
		cfg := &Configuration{}
		httpClient := &http.Client{Timeout: 30 * time.Second}
		setter := WithHttpClientPtr(httpClient)
		setter(cfg)
		assert.NotNil(t, cfg.HTTPClient)
		assert.Equal(t, 30*time.Second, cfg.HTTPClient.Timeout)
	})

	t.Run("WithProxyPort", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithProxyPort(8080)
		setter(cfg)
		assert.Equal(t, int32(8080), cfg.Zscaler.Client.Proxy.Port)
	})

	t.Run("WithProxyHost", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithProxyHost("proxy.example.com")
		setter(cfg)
		assert.Equal(t, "proxy.example.com", cfg.Zscaler.Client.Proxy.Host)
	})

	t.Run("WithProxyUsername", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithProxyUsername("proxyuser")
		setter(cfg)
		assert.Equal(t, "proxyuser", cfg.Zscaler.Client.Proxy.Username)
	})

	t.Run("WithProxyPassword", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithProxyPassword("proxypass")
		setter(cfg)
		assert.Equal(t, "proxypass", cfg.Zscaler.Client.Proxy.Password)
	})

	t.Run("WithTestingDisableHttpsCheck", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithTestingDisableHttpsCheck(true)
		setter(cfg)
		assert.True(t, cfg.Zscaler.Testing.DisableHttpsCheck)
	})

	t.Run("WithRequestTimeout", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithRequestTimeout(2 * time.Minute)
		setter(cfg)
		assert.Equal(t, 2*time.Minute, cfg.Zscaler.Client.RequestTimeout)
	})

	t.Run("WithRateLimitMaxRetries", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithRateLimitMaxRetries(10)
		setter(cfg)
		assert.Equal(t, int32(10), cfg.Zscaler.Client.RateLimit.MaxRetries)
	})

	t.Run("WithRateLimitMaxWait", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithRateLimitMaxWait(30 * time.Second)
		setter(cfg)
		assert.Equal(t, 30*time.Second, cfg.Zscaler.Client.RateLimit.RetryWaitMax)
	})

	t.Run("WithRateLimitMinWait", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithRateLimitMinWait(5 * time.Second)
		setter(cfg)
		assert.Equal(t, 5*time.Second, cfg.Zscaler.Client.RateLimit.RetryWaitMin)
	})

	t.Run("WithRateLimitRemainingThreshold", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithRateLimitRemainingThreshold(3)
		setter(cfg)
		assert.Equal(t, int32(3), cfg.Zscaler.Client.RateLimit.RetryRemainingThreshold)
	})

	t.Run("WithRateLimitMaxSessionNotValidRetries", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithRateLimitMaxSessionNotValidRetries(5)
		setter(cfg)
		assert.Equal(t, int32(5), cfg.Zscaler.Client.RateLimit.MaxSessionNotValidRetries)
	})

	t.Run("WithUserAgentExtra", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithUserAgentExtra("terraform-provider/1.0.0")
		setter(cfg)
		assert.Equal(t, "terraform-provider/1.0.0", cfg.UserAgentExtra)
	})

	t.Run("WithDebug", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithDebug(true)
		setter(cfg)
		assert.True(t, cfg.Debug)
	})

	t.Run("WithLegacyClient", func(t *testing.T) {
		cfg := &Configuration{}
		setter := WithLegacyClient(true)
		setter(cfg)
		assert.True(t, cfg.UseLegacyClient)
	})
}

func TestAddDefaultHeader(t *testing.T) {
	t.Run("Add single header", func(t *testing.T) {
		cfg := &Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Custom", "value")
		assert.Equal(t, "value", cfg.DefaultHeader["X-Custom"])
	})

	t.Run("Override header", func(t *testing.T) {
		cfg := &Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Custom", "old")
		cfg.AddDefaultHeader("X-Custom", "new")
		assert.Equal(t, "new", cfg.DefaultHeader["X-Custom"])
	})
}

