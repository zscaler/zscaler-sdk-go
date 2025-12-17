// Package zdx provides unit tests for ZDX client configuration
package zdx

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx"
)

// =====================================================
// Configuration Tests
// =====================================================

func TestConfiguration_Structs(t *testing.T) {
	t.Parallel()

	t.Run("AuthRequest JSON marshaling", func(t *testing.T) {
		req := zdx.AuthRequest{
			APIKeyID:     "test-key-id",
			APIKeySecret: "test-secret",
			Timestamp:    1699900000,
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"key_id":"test-key-id"`)
		assert.Contains(t, string(data), `"key_secret":"test-secret"`)
		assert.Contains(t, string(data), `"timestamp":1699900000`)
	})

	t.Run("AuthRequest JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"key_id": "my-api-key",
			"key_secret": "my-api-secret",
			"timestamp": 1700000000
		}`

		var req zdx.AuthRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, "my-api-key", req.APIKeyID)
		assert.Equal(t, "my-api-secret", req.APIKeySecret)
		assert.Equal(t, int64(1700000000), req.Timestamp)
	})

	t.Run("AuthToken JSON marshaling", func(t *testing.T) {
		token := zdx.AuthToken{
			TokenType:   "Bearer",
			AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			ExpiresIn:   3600,
		}

		data, err := json.Marshal(token)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"token_type":"Bearer"`)
		assert.Contains(t, string(data), `"expires_in":3600`)
	})

	t.Run("AuthToken JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"token_type": "Bearer",
			"token": "test-access-token",
			"expires_in": 7200
		}`

		var token zdx.AuthToken
		err := json.Unmarshal([]byte(jsonData), &token)
		require.NoError(t, err)

		assert.Equal(t, "Bearer", token.TokenType)
		assert.Equal(t, "test-access-token", token.AccessToken)
		assert.Equal(t, 7200, token.ExpiresIn)
	})

	t.Run("BackoffConfig struct initialization", func(t *testing.T) {
		config := zdx.BackoffConfig{
			Enabled:             true,
			RetryWaitMinSeconds: 5,
			RetryWaitMaxSeconds: 20,
			MaxNumOfRetries:     10,
		}

		assert.True(t, config.Enabled)
		assert.Equal(t, 5, config.RetryWaitMinSeconds)
		assert.Equal(t, 20, config.RetryWaitMaxSeconds)
		assert.Equal(t, 10, config.MaxNumOfRetries)
	})
}

func TestConfigSetter_Functions(t *testing.T) {
	t.Parallel()

	t.Run("WithZDXAPIKeyID", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithZDXAPIKeyID("test-api-key")
		setter(cfg)
		assert.Equal(t, "test-api-key", cfg.ZDX.Client.ZDXAPIKeyID)
	})

	t.Run("WithZDXAPISecret", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithZDXAPISecret("test-secret")
		setter(cfg)
		assert.Equal(t, "test-secret", cfg.ZDX.Client.ZDXAPISecret)
	})

	t.Run("WithZDXCloud", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithZDXCloud("zscloud.net")
		setter(cfg)
		assert.Equal(t, "zscloud.net", cfg.ZDX.Client.ZDXCloud)
	})

	t.Run("WithDebug", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithDebug(true)
		setter(cfg)
		assert.True(t, cfg.Debug)
	})

	t.Run("WithUserAgentExtra", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithUserAgentExtra("custom-agent/1.0")
		setter(cfg)
		assert.Equal(t, "custom-agent/1.0", cfg.UserAgentExtra)
	})

	t.Run("WithTestingDisableHttpsCheck", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithTestingDisableHttpsCheck(true)
		setter(cfg)
		assert.True(t, cfg.ZDX.Testing.DisableHttpsCheck)
	})

	t.Run("WithRateLimitMaxRetries", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithRateLimitMaxRetries(5)
		setter(cfg)
		assert.Equal(t, int32(5), cfg.ZDX.Client.RateLimit.MaxRetries)
	})

	t.Run("WithProxyHost", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithProxyHost("proxy.example.com")
		setter(cfg)
		assert.Equal(t, "proxy.example.com", cfg.ZDX.Client.Proxy.Host)
	})

	t.Run("WithProxyPort", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithProxyPort(8080)
		setter(cfg)
		assert.Equal(t, int32(8080), cfg.ZDX.Client.Proxy.Port)
	})

	t.Run("WithPartnerID", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithPartnerID("partner-123")
		setter(cfg)
		assert.Equal(t, "partner-123", cfg.ZDX.Client.PartnerID)
	})

	t.Run("WithProxyUsername", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithProxyUsername("proxy-user")
		setter(cfg)
		assert.Equal(t, "proxy-user", cfg.ZDX.Client.Proxy.Username)
	})

	t.Run("WithProxyPassword", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithProxyPassword("proxy-pass")
		setter(cfg)
		assert.Equal(t, "proxy-pass", cfg.ZDX.Client.Proxy.Password)
	})

	t.Run("WithRequestTimeout", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithRequestTimeout(30 * time.Second)
		setter(cfg)
		assert.Equal(t, 30*time.Second, cfg.ZDX.Client.RequestTimeout)
	})

	t.Run("WithRateLimitMaxWait", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithRateLimitMaxWait(10 * time.Second)
		setter(cfg)
		assert.Equal(t, 10*time.Second, cfg.ZDX.Client.RateLimit.RetryWaitMax)
	})

	t.Run("WithRateLimitMinWait", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		setter := zdx.WithRateLimitMinWait(2 * time.Second)
		setter(cfg)
		assert.Equal(t, 2*time.Second, cfg.ZDX.Client.RateLimit.RetryWaitMin)
	})

	t.Run("WithHttpClientPtr", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		httpClient := &http.Client{Timeout: 5 * time.Second}
		setter := zdx.WithHttpClientPtr(httpClient)
		setter(cfg)
		assert.NotNil(t, cfg.HTTPClient)
		assert.Equal(t, 5*time.Second, cfg.HTTPClient.Timeout)
	})
}

func TestConstants(t *testing.T) {
	t.Parallel()

	t.Run("VERSION constant is set", func(t *testing.T) {
		assert.NotEmpty(t, zdx.VERSION)
	})

	t.Run("Environment variable names", func(t *testing.T) {
		assert.Equal(t, "ZDX_API_KEY_ID", zdx.ZDX_API_KEY_ID)
		assert.Equal(t, "ZDX_API_SECRET", zdx.ZDX_API_SECRET)
	})
}

func TestCloudURLs(t *testing.T) {
	t.Parallel()

	// Test known cloud names
	clouds := []string{
		"zscloud.net",
		"zscalerone.net",
		"zscalertwo.net",
		"zscalerthree.net",
		"zspreview.net",
		"zsdemo.net",
		"zscalerbeta.net",
		"zsclouddev.net",
	}

	for _, cloud := range clouds {
		t.Run("Cloud: "+cloud, func(t *testing.T) {
			cfg := &zdx.Configuration{}
			setter := zdx.WithZDXCloud(cloud)
			setter(cfg)
			assert.Equal(t, cloud, cfg.ZDX.Client.ZDXCloud)
		})
	}
}

func TestConfiguration_DefaultValues(t *testing.T) {
	t.Parallel()

	t.Run("Configuration default state", func(t *testing.T) {
		cfg := &zdx.Configuration{}

		// Default values should be empty/zero
		assert.Empty(t, cfg.ZDX.Client.ZDXAPIKeyID)
		assert.Empty(t, cfg.ZDX.Client.ZDXAPISecret)
		assert.Empty(t, cfg.UserAgent)
		assert.False(t, cfg.Debug)
	})

	t.Run("BackoffConfig disabled by default", func(t *testing.T) {
		config := zdx.BackoffConfig{}

		assert.False(t, config.Enabled)
		assert.Equal(t, 0, config.RetryWaitMinSeconds)
		assert.Equal(t, 0, config.RetryWaitMaxSeconds)
		assert.Equal(t, 0, config.MaxNumOfRetries)
	})
}

func TestConfiguration_AddDefaultHeader(t *testing.T) {
	t.Parallel()

	t.Run("Add single header", func(t *testing.T) {
		cfg := &zdx.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Custom-Header", "custom-value")
		assert.Equal(t, "custom-value", cfg.DefaultHeader["X-Custom-Header"])
	})

	t.Run("Add multiple headers", func(t *testing.T) {
		cfg := &zdx.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("Header1", "Value1")
		cfg.AddDefaultHeader("Header2", "Value2")
		assert.Equal(t, "Value1", cfg.DefaultHeader["Header1"])
		assert.Equal(t, "Value2", cfg.DefaultHeader["Header2"])
	})

	t.Run("Override existing header", func(t *testing.T) {
		cfg := &zdx.Configuration{
			DefaultHeader: make(map[string]string),
		}
		cfg.AddDefaultHeader("X-Header", "old-value")
		cfg.AddDefaultHeader("X-Header", "new-value")
		assert.Equal(t, "new-value", cfg.DefaultHeader["X-Header"])
	})
}

func TestConfiguration_SetBackoffConfig(t *testing.T) {
	t.Parallel()

	t.Run("Set backoff config", func(t *testing.T) {
		cfg := &zdx.Configuration{}
		backoffCfg := &zdx.BackoffConfig{
			Enabled:             true,
			RetryWaitMinSeconds: 5,
			RetryWaitMaxSeconds: 30,
			MaxNumOfRetries:     10,
		}
		cfg.SetBackoffConfig(backoffCfg)
		assert.NotNil(t, cfg.ZDX.Client.RateLimit.BackoffConf)
		assert.True(t, cfg.ZDX.Client.RateLimit.BackoffConf.Enabled)
		assert.Equal(t, 10, cfg.ZDX.Client.RateLimit.BackoffConf.MaxNumOfRetries)
	})
}

// =====================================================
// Client Creation Tests
// =====================================================

func TestNewConfiguration_MissingCredentials(t *testing.T) {
	t.Run("NewConfiguration missing credentials returns error", func(t *testing.T) {
		_, err := zdx.NewConfiguration()
		assert.Error(t, err)
	})
}

func TestAuthToken_Structure(t *testing.T) {
	t.Run("AuthToken fields", func(t *testing.T) {
		token := &zdx.AuthToken{
			TokenType:   "Bearer",
			AccessToken: "test-token",
			ExpiresIn:   3600,
		}

		assert.Equal(t, "Bearer", token.TokenType)
		assert.Equal(t, "test-token", token.AccessToken)
		assert.Equal(t, 3600, token.ExpiresIn)
	})

	t.Run("AuthToken JSON round-trip", func(t *testing.T) {
		token := zdx.AuthToken{
			TokenType:   "Bearer",
			AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			ExpiresIn:   7200,
		}

		data, err := json.Marshal(token)
		require.NoError(t, err)

		var parsed zdx.AuthToken
		err = json.Unmarshal(data, &parsed)
		require.NoError(t, err)

		assert.Equal(t, token.TokenType, parsed.TokenType)
		assert.Equal(t, token.ExpiresIn, parsed.ExpiresIn)
	})
}

