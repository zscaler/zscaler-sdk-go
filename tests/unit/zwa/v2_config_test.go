// Package zwa provides unit tests for ZWA client configuration
package zwa

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa"
)

func TestZWA_ConfigSetter_Functions(t *testing.T) {
	t.Parallel()

	t.Run("WithZWAAPIKeyID", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithZWAAPIKeyID("test-key-id")
		setter(cfg)
		assert.Equal(t, "test-key-id", cfg.ZWA.Client.ZWAAPIKeyID)
	})

	t.Run("WithZWAAPISecret", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithZWAAPISecret("test-secret")
		setter(cfg)
		assert.Equal(t, "test-secret", cfg.ZWA.Client.ZWAAPISecret)
	})

	t.Run("WithZWACloud", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithZWACloud("us1")
		setter(cfg)
		assert.Equal(t, "us1", cfg.ZWA.Client.ZWACloud)
	})

	t.Run("WithPartnerID", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithPartnerID("partner-123")
		setter(cfg)
		assert.Equal(t, "partner-123", cfg.ZWA.Client.PartnerID)
	})

	t.Run("WithHttpClientPtr", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		httpClient := &http.Client{Timeout: 30 * time.Second}
		setter := zwa.WithHttpClientPtr(httpClient)
		setter(cfg)
		assert.Equal(t, httpClient, cfg.HTTPClient)
	})

	t.Run("WithProxyPort", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithProxyPort(8080)
		setter(cfg)
		assert.Equal(t, int32(8080), cfg.ZWA.Client.Proxy.Port)
	})

	t.Run("WithProxyHost", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithProxyHost("proxy.example.com")
		setter(cfg)
		assert.Equal(t, "proxy.example.com", cfg.ZWA.Client.Proxy.Host)
	})

	t.Run("WithProxyUsername", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithProxyUsername("proxyuser")
		setter(cfg)
		assert.Equal(t, "proxyuser", cfg.ZWA.Client.Proxy.Username)
	})

	t.Run("WithProxyPassword", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithProxyPassword("proxypass")
		setter(cfg)
		assert.Equal(t, "proxypass", cfg.ZWA.Client.Proxy.Password)
	})

	t.Run("WithTestingDisableHttpsCheck", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithTestingDisableHttpsCheck(true)
		setter(cfg)
		assert.True(t, cfg.ZWA.Testing.DisableHttpsCheck)
	})

	t.Run("WithRequestTimeout", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithRequestTimeout(45 * time.Second)
		setter(cfg)
		assert.Equal(t, 45*time.Second, cfg.ZWA.Client.RequestTimeout)
	})

	t.Run("WithRateLimitMaxRetries", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithRateLimitMaxRetries(10)
		setter(cfg)
		assert.Equal(t, int32(10), cfg.ZWA.Client.RateLimit.MaxRetries)
	})

	t.Run("WithRateLimitMaxWait", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithRateLimitMaxWait(30 * time.Second)
		setter(cfg)
		assert.Equal(t, 30*time.Second, cfg.ZWA.Client.RateLimit.RetryWaitMax)
	})

	t.Run("WithRateLimitMinWait", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithRateLimitMinWait(5 * time.Second)
		setter(cfg)
		assert.Equal(t, 5*time.Second, cfg.ZWA.Client.RateLimit.RetryWaitMin)
	})

	t.Run("WithUserAgentExtra", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithUserAgentExtra("CustomAgent/1.0")
		setter(cfg)
		assert.Equal(t, "CustomAgent/1.0", cfg.UserAgentExtra)
	})

	t.Run("WithDebug", func(t *testing.T) {
		cfg := &zwa.Configuration{}
		setter := zwa.WithDebug(true)
		setter(cfg)
		assert.True(t, cfg.Debug)
	})
}

func TestZWA_Configuration_Methods(t *testing.T) {
	t.Parallel()

	t.Run("AddDefaultHeader", func(t *testing.T) {
		cfg := &zwa.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Custom-Header", "custom-value")
		assert.Equal(t, "custom-value", cfg.DefaultHeader["X-Custom-Header"])
	})

	t.Run("ContextKey String method", func(t *testing.T) {
		// Access through context token
		assert.NotNil(t, zwa.ContextAccessToken)
	})
}

func TestZWA_Configuration_SetBackoffConfig(t *testing.T) {
	t.Parallel()

	t.Run("SetBackoffConfig on Configuration", func(t *testing.T) {
		baseURL, _ := url.Parse("https://api.test.zsworkflow.net")
		cfg := &zwa.Configuration{
			BaseURL: baseURL,
		}

		backoffCfg := &zwa.BackoffConfig{
			Enabled:             true,
			RetryWaitMinSeconds: 2,
			RetryWaitMaxSeconds: 10,
			MaxNumOfRetries:     5,
		}

		cfg.SetBackoffConfig(backoffCfg)

		assert.True(t, cfg.ZWA.Client.RateLimit.BackoffConf.Enabled)
		assert.Equal(t, 2, cfg.ZWA.Client.RateLimit.BackoffConf.RetryWaitMinSeconds)
		assert.Equal(t, 10, cfg.ZWA.Client.RateLimit.BackoffConf.RetryWaitMaxSeconds)
		assert.Equal(t, 5, cfg.ZWA.Client.RateLimit.BackoffConf.MaxNumOfRetries)
	})
}

func TestZWA_AuthToken_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AuthToken JSON structure", func(t *testing.T) {
		token := zwa.AuthToken{
			TokenType:   "Bearer",
			AccessToken: "mock-access-token",
			ExpiresIn:   3600,
		}

		assert.Equal(t, "Bearer", token.TokenType)
		assert.Equal(t, "mock-access-token", token.AccessToken)
		assert.Equal(t, 3600, token.ExpiresIn)
	})

	t.Run("AuthRequest structure", func(t *testing.T) {
		req := zwa.AuthRequest{
			APIKeyID:     "key-id",
			APIKeySecret: "key-secret",
			Timestamp:    1234567890,
		}

		assert.Equal(t, "key-id", req.APIKeyID)
		assert.Equal(t, "key-secret", req.APIKeySecret)
		assert.Equal(t, int64(1234567890), req.Timestamp)
	})

	t.Run("BackoffConfig structure", func(t *testing.T) {
		cfg := zwa.BackoffConfig{
			Enabled:             true,
			RetryWaitMinSeconds: 5,
			RetryWaitMaxSeconds: 20,
			MaxNumOfRetries:     3,
		}

		assert.True(t, cfg.Enabled)
		assert.Equal(t, 5, cfg.RetryWaitMinSeconds)
		assert.Equal(t, 20, cfg.RetryWaitMaxSeconds)
		assert.Equal(t, 3, cfg.MaxNumOfRetries)
	})
}

func TestZWA_Constants(t *testing.T) {
	t.Parallel()

	t.Run("VERSION constant", func(t *testing.T) {
		assert.NotEmpty(t, zwa.VERSION)
	})

	t.Run("Environment variable constants", func(t *testing.T) {
		assert.Equal(t, "ZWA_API_KEY_ID", zwa.ZWA_API_KEY_ID)
		assert.Equal(t, "ZWA_API_SECRET", zwa.ZWA_API_SECRET)
	})

	t.Run("Retry constants", func(t *testing.T) {
		assert.Equal(t, 100, zwa.MaxNumOfRetries)
		assert.Equal(t, 20, zwa.RetryWaitMaxSeconds)
		assert.Equal(t, 5, zwa.RetryWaitMinSeconds)
	})
}

func TestZWA_Configuration_Combined(t *testing.T) {
	t.Parallel()

	t.Run("Apply multiple configuration options", func(t *testing.T) {
		cfg := &zwa.Configuration{}

		zwa.WithZWAAPIKeyID("test-key-id")(cfg)
		zwa.WithZWAAPISecret("test-secret")(cfg)
		zwa.WithZWACloud("us1")(cfg)
		zwa.WithDebug(true)(cfg)

		assert.Equal(t, "test-key-id", cfg.ZWA.Client.ZWAAPIKeyID)
		assert.Equal(t, "test-secret", cfg.ZWA.Client.ZWAAPISecret)
		assert.Equal(t, "us1", cfg.ZWA.Client.ZWACloud)
		assert.True(t, cfg.Debug)
	})

	t.Run("Configuration with proxy settings", func(t *testing.T) {
		cfg := &zwa.Configuration{}

		zwa.WithProxyHost("proxy.example.com")(cfg)
		zwa.WithProxyPort(8080)(cfg)
		zwa.WithProxyUsername("proxyuser")(cfg)
		zwa.WithProxyPassword("proxypass")(cfg)

		assert.Equal(t, "proxy.example.com", cfg.ZWA.Client.Proxy.Host)
		assert.Equal(t, int32(8080), cfg.ZWA.Client.Proxy.Port)
		assert.Equal(t, "proxyuser", cfg.ZWA.Client.Proxy.Username)
		assert.Equal(t, "proxypass", cfg.ZWA.Client.Proxy.Password)
	})

	t.Run("Configuration with rate limit settings", func(t *testing.T) {
		cfg := &zwa.Configuration{}

		zwa.WithRateLimitMaxRetries(5)(cfg)
		zwa.WithRateLimitMinWait(2 * time.Second)(cfg)
		zwa.WithRateLimitMaxWait(30 * time.Second)(cfg)

		assert.Equal(t, int32(5), cfg.ZWA.Client.RateLimit.MaxRetries)
		assert.Equal(t, 2*time.Second, cfg.ZWA.Client.RateLimit.RetryWaitMin)
		assert.Equal(t, 30*time.Second, cfg.ZWA.Client.RateLimit.RetryWaitMax)
	})
}

func TestZWA_ContextKey_String(t *testing.T) {
	t.Parallel()

	t.Run("ContextAccessToken has value", func(t *testing.T) {
		// ContextAccessToken is a contextKey with string value "access_token"
		assert.NotNil(t, zwa.ContextAccessToken)
	})
}

func TestZWA_Cloud_Settings(t *testing.T) {
	t.Parallel()

	testClouds := []string{
		"beta",
		"us1",
		"us2",
		"eu1",
		"eu2",
		"au1",
	}

	for _, cloud := range testClouds {
		cloud := cloud // capture range variable
		t.Run("Set cloud to "+cloud, func(t *testing.T) {
			cfg := &zwa.Configuration{}
			setter := zwa.WithZWACloud(cloud)
			setter(cfg)
			assert.Equal(t, cloud, cfg.ZWA.Client.ZWACloud)
		})
	}
}

func TestZWA_Configuration_RateLimit(t *testing.T) {
	t.Parallel()

	t.Run("Full rate limit configuration", func(t *testing.T) {
		cfg := &zwa.Configuration{}

		zwa.WithRateLimitMaxRetries(10)(cfg)
		zwa.WithRateLimitMinWait(5 * time.Second)(cfg)
		zwa.WithRateLimitMaxWait(60 * time.Second)(cfg)

		assert.Equal(t, int32(10), cfg.ZWA.Client.RateLimit.MaxRetries)
		assert.Equal(t, 5*time.Second, cfg.ZWA.Client.RateLimit.RetryWaitMin)
		assert.Equal(t, 60*time.Second, cfg.ZWA.Client.RateLimit.RetryWaitMax)
	})

	t.Run("SetBackoffConfig full configuration", func(t *testing.T) {
		cfg := &zwa.Configuration{}

		backoff := &zwa.BackoffConfig{
			Enabled:             true,
			MaxNumOfRetries:     5,
			RetryWaitMinSeconds: 10,
			RetryWaitMaxSeconds: 120,
		}
		cfg.SetBackoffConfig(backoff)

		assert.NotNil(t, cfg.ZWA.Client.RateLimit.BackoffConf)
		assert.True(t, cfg.ZWA.Client.RateLimit.BackoffConf.Enabled)
		assert.Equal(t, 5, cfg.ZWA.Client.RateLimit.BackoffConf.MaxNumOfRetries)
		assert.Equal(t, 10, cfg.ZWA.Client.RateLimit.BackoffConf.RetryWaitMinSeconds)
		assert.Equal(t, 120, cfg.ZWA.Client.RateLimit.BackoffConf.RetryWaitMaxSeconds)
	})
}

func TestZWA_DefaultHeader(t *testing.T) {
	t.Parallel()

	t.Run("Configuration with default headers", func(t *testing.T) {
		cfg := &zwa.Configuration{
			DefaultHeader: map[string]string{
				"X-Custom-Header": "custom-value",
				"Accept":          "application/json",
			},
		}

		assert.Equal(t, "custom-value", cfg.DefaultHeader["X-Custom-Header"])
		assert.Equal(t, "application/json", cfg.DefaultHeader["Accept"])
	})

	t.Run("AddDefaultHeader method", func(t *testing.T) {
		cfg := &zwa.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Test-Header", "test-value")

		assert.Equal(t, "test-value", cfg.DefaultHeader["X-Test-Header"])
	})

	t.Run("AddDefaultHeader multiple headers", func(t *testing.T) {
		cfg := &zwa.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-API-Version", "v1")
		cfg.AddDefaultHeader("X-Request-ID", "req-12345")
		cfg.AddDefaultHeader("X-Correlation-ID", "corr-67890")

		assert.Equal(t, "v1", cfg.DefaultHeader["X-API-Version"])
		assert.Equal(t, "req-12345", cfg.DefaultHeader["X-Request-ID"])
		assert.Equal(t, "corr-67890", cfg.DefaultHeader["X-Correlation-ID"])
	})

	t.Run("AddDefaultHeader override existing", func(t *testing.T) {
		cfg := &zwa.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Custom", "old-value")
		assert.Equal(t, "old-value", cfg.DefaultHeader["X-Custom"])

		cfg.AddDefaultHeader("X-Custom", "new-value")
		assert.Equal(t, "new-value", cfg.DefaultHeader["X-Custom"])
	})
}

func TestZWA_NewConfiguration_MissingCredentials(t *testing.T) {
	t.Run("NewConfiguration missing credentials returns error", func(t *testing.T) {
		_, err := zwa.NewConfiguration()
		assert.Error(t, err)
	})
}

func TestZWA_Client_NilConfiguration(t *testing.T) {
	t.Run("NewClient with nil configuration returns error", func(t *testing.T) {
		client, err := zwa.NewClient(nil)
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "configuration cannot be nil")
	})
}

