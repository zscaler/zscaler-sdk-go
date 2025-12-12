// Package oneapi provides unit tests for the OneAPI client
package oneapi

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func TestAuthToken_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AuthToken JSON marshaling", func(t *testing.T) {
		token := zscaler.AuthToken{
			TokenType:   "Bearer",
			AccessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
			ExpiresIn:   json.Number("3600"),
			Expiry:      time.Now().Add(time.Hour),
		}

		data, err := json.Marshal(token)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"token_type":"Bearer"`)
		assert.Contains(t, string(data), `"access_token"`)
		// expires_in is marshaled as a number
		assert.Contains(t, string(data), `"expires_in":3600`)
	})

	t.Run("AuthToken JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"token_type": "Bearer",
			"access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.test.signature",
			"expires_in": "7200"
		}`

		var token zscaler.AuthToken
		err := json.Unmarshal([]byte(jsonData), &token)
		require.NoError(t, err)

		assert.Equal(t, "Bearer", token.TokenType)
		assert.Equal(t, json.Number("7200"), token.ExpiresIn)
	})

	t.Run("AuthToken with numeric expires_in", func(t *testing.T) {
		jsonData := `{
			"token_type": "Bearer",
			"access_token": "test_token",
			"expires_in": 3600
		}`

		var token zscaler.AuthToken
		err := json.Unmarshal([]byte(jsonData), &token)
		require.NoError(t, err)

		expiresIn, err := token.ExpiresIn.Int64()
		require.NoError(t, err)
		assert.Equal(t, int64(3600), expiresIn)
	})
}

func TestConfiguration_Defaults(t *testing.T) {
	t.Parallel()

	t.Run("Configuration has expected defaults", func(t *testing.T) {
		// Test that VERSION constant is set
		assert.NotEmpty(t, zscaler.VERSION)
	})

	t.Run("Context key string conversion", func(t *testing.T) {
		key := zscaler.ContextAccessToken
		assert.Contains(t, key.String(), "zscaler")
	})
}

func TestGetAPIBaseURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		cloud    string
		expected string
	}{
		{
			name:     "Empty cloud returns default",
			cloud:    "",
			expected: "https://api.zsapi.net",
		},
		{
			name:     "PRODUCTION cloud returns default",
			cloud:    "PRODUCTION",
			expected: "https://api.zsapi.net",
		},
		{
			name:     "production lowercase returns default",
			cloud:    "production",
			expected: "https://api.zsapi.net",
		},
		{
			name:     "BETA cloud returns beta URL",
			cloud:    "BETA",
			expected: "https://api.beta.zsapi.net",
		},
		{
			name:     "ZSCALERTHREE cloud returns zscalerthree URL",
			cloud:    "ZSCALERTHREE",
			expected: "https://api.zscalerthree.zsapi.net",
		},
		{
			name:     "ZSCLOUD cloud returns zscloud URL",
			cloud:    "ZSCLOUD",
			expected: "https://api.zscloud.zsapi.net",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := zscaler.GetAPIBaseURL(tc.cloud)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConfigSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithClientID sets client ID", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithClientID("test-client-id")
		setter(cfg)
		assert.Equal(t, "test-client-id", cfg.Zscaler.Client.ClientID)
	})

	t.Run("WithClientSecret sets client secret", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithClientSecret("test-secret")
		setter(cfg)
		assert.Equal(t, "test-secret", cfg.Zscaler.Client.ClientSecret)
	})

	t.Run("WithVanityDomain sets vanity domain", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithVanityDomain("company")
		setter(cfg)
		assert.Equal(t, "company", cfg.Zscaler.Client.VanityDomain)
	})

	t.Run("WithZscalerCloud sets cloud", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithZscalerCloud("BETA")
		setter(cfg)
		assert.Equal(t, "BETA", cfg.Zscaler.Client.Cloud)
	})

	t.Run("WithZPACustomerID sets customer ID", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithZPACustomerID("123456")
		setter(cfg)
		assert.Equal(t, "123456", cfg.Zscaler.Client.CustomerID)
	})

	t.Run("WithZPAMicrotenantID sets microtenant ID", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithZPAMicrotenantID("mt-123")
		setter(cfg)
		assert.Equal(t, "mt-123", cfg.Zscaler.Client.MicrotenantID)
	})

	t.Run("WithPartnerID sets partner ID", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithPartnerID("partner-123")
		setter(cfg)
		assert.Equal(t, "partner-123", cfg.Zscaler.Client.PartnerID)
	})

	t.Run("WithCache enables cache", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithCache(true)
		setter(cfg)
		assert.True(t, cfg.Zscaler.Client.Cache.Enabled)
	})

	t.Run("WithCacheTtl sets cache TTL", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		cfg.Zscaler.Client.Cache.DefaultTti = 5 * time.Minute
		setter := zscaler.WithCacheTtl(15 * time.Minute)
		setter(cfg)
		assert.Equal(t, 15*time.Minute, cfg.Zscaler.Client.Cache.DefaultTtl)
	})

	t.Run("WithProxyHost sets proxy host", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithProxyHost("proxy.company.com")
		setter(cfg)
		assert.Equal(t, "proxy.company.com", cfg.Zscaler.Client.Proxy.Host)
	})

	t.Run("WithProxyPort sets proxy port", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithProxyPort(8080)
		setter(cfg)
		assert.Equal(t, int32(8080), cfg.Zscaler.Client.Proxy.Port)
	})

	t.Run("WithSandboxToken sets sandbox token", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithSandboxToken("sandbox-token-123")
		setter(cfg)
		assert.Equal(t, "sandbox-token-123", cfg.Zscaler.Client.SandboxToken)
	})

	t.Run("WithSandboxCloud sets sandbox cloud", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithSandboxCloud("zspreview")
		setter(cfg)
		assert.Equal(t, "zspreview", cfg.Zscaler.Client.SandboxCloud)
	})

	t.Run("WithTestingDisableHttpsCheck disables HTTPS check", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithTestingDisableHttpsCheck(true)
		setter(cfg)
		assert.True(t, cfg.Zscaler.Testing.DisableHttpsCheck)
	})

	t.Run("WithRateLimitMaxRetries sets max retries", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithRateLimitMaxRetries(10)
		setter(cfg)
		assert.Equal(t, int32(10), cfg.Zscaler.Client.RateLimit.MaxRetries)
	})

	t.Run("WithRateLimitMaxSessionNotValidRetries sets max session retries", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithRateLimitMaxSessionNotValidRetries(5)
		setter(cfg)
		assert.Equal(t, int32(5), cfg.Zscaler.Client.RateLimit.MaxSessionNotValidRetries)
	})

	t.Run("WithUserAgentExtra sets extra user agent", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithUserAgentExtra("terraform/1.0.0")
		setter(cfg)
		assert.Equal(t, "terraform/1.0.0", cfg.UserAgentExtra)
	})

	t.Run("WithDebug enables debug mode", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithDebug(true)
		setter(cfg)
		assert.True(t, cfg.Debug)
	})

	t.Run("WithLegacyClient enables legacy client", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		setter := zscaler.WithLegacyClient(true)
		setter(cfg)
		assert.True(t, cfg.UseLegacyClient)
	})
}

