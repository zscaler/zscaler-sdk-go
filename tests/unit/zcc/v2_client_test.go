// Package zcc provides unit tests for ZCC v2_client and v2_config
package zcc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc"
)

// =====================================================
// Configuration Structure Tests
// =====================================================

func TestConfiguration_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AuthToken JSON marshaling", func(t *testing.T) {
		token := zcc.AuthToken{
			TokenType:   "Bearer",
			AccessToken: "jwt-token-123",
			ExpiresIn:   3600,
			Expiry:      time.Now().Add(time.Hour),
		}

		data, err := json.Marshal(token)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"token_type":"Bearer"`)
		assert.Contains(t, string(data), `"jwtToken":"jwt-token-123"`)
	})

	t.Run("AuthToken JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"token_type": "Bearer",
			"jwtToken": "test-jwt-token-abc",
			"expires_in": "3600"
		}`

		var token zcc.AuthToken
		err := json.Unmarshal([]byte(jsonData), &token)
		require.NoError(t, err)

		assert.Equal(t, "Bearer", token.TokenType)
		assert.Equal(t, "test-jwt-token-abc", token.AccessToken)
		assert.Equal(t, "3600", token.RawExpiresIn)
	})

	t.Run("AuthRequest JSON marshaling", func(t *testing.T) {
		request := zcc.AuthRequest{
			APIKey:    "my-api-key",
			SecretKey: "my-secret-key",
		}

		data, err := json.Marshal(request)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"apiKey":"my-api-key"`)
		assert.Contains(t, string(data), `"secretKey":"my-secret-key"`)
	})

	t.Run("AuthRequest JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"apiKey": "api-key-xyz",
			"secretKey": "secret-key-xyz"
		}`

		var request zcc.AuthRequest
		err := json.Unmarshal([]byte(jsonData), &request)
		require.NoError(t, err)

		assert.Equal(t, "api-key-xyz", request.APIKey)
		assert.Equal(t, "secret-key-xyz", request.SecretKey)
	})

	t.Run("BackoffConfig structure", func(t *testing.T) {
		config := zcc.BackoffConfig{
			Enabled:             true,
			RetryWaitMinSeconds: 5,
			RetryWaitMaxSeconds: 60,
			MaxNumOfRetries:     3,
		}

		assert.True(t, config.Enabled)
		assert.Equal(t, 5, config.RetryWaitMinSeconds)
		assert.Equal(t, 60, config.RetryWaitMaxSeconds)
		assert.Equal(t, 3, config.MaxNumOfRetries)
	})
}

// =====================================================
// Configuration Options Tests - Direct option application
// =====================================================

func TestConfiguration_Options(t *testing.T) {
	t.Parallel()

	t.Run("WithZCCClientID sets client ID", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithZCCClientID("test-client-id")(cfg)
		assert.Equal(t, "test-client-id", cfg.ZCC.Client.ZCCClientID)
	})

	t.Run("WithZCCClientSecret sets client secret", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithZCCClientSecret("test-client-secret")(cfg)
		assert.Equal(t, "test-client-secret", cfg.ZCC.Client.ZCCClientSecret)
	})

	t.Run("WithZCCCloud sets cloud", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithZCCCloud("zscalertwo")(cfg)
		assert.Equal(t, "zscalertwo", cfg.ZCC.Client.ZCCCloud)
	})

	t.Run("WithDebug sets debug mode", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithDebug(true)(cfg)
		assert.True(t, cfg.Debug)
	})

	t.Run("WithUserAgentExtra sets user agent", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithUserAgentExtra("my-custom-agent/1.0")(cfg)
		assert.Contains(t, cfg.UserAgentExtra, "my-custom-agent/1.0")
	})

	t.Run("WithProxyHost sets proxy host", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithProxyHost("proxy.example.com")(cfg)
		assert.Equal(t, "proxy.example.com", cfg.ZCC.Client.Proxy.Host)
	})

	t.Run("WithProxyPort sets proxy port", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithProxyPort(8080)(cfg)
		assert.Equal(t, int32(8080), cfg.ZCC.Client.Proxy.Port)
	})

	t.Run("WithProxyUsername sets proxy username", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithProxyUsername("proxyuser")(cfg)
		assert.Equal(t, "proxyuser", cfg.ZCC.Client.Proxy.Username)
	})

	t.Run("WithProxyPassword sets proxy password", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithProxyPassword("proxypass")(cfg)
		assert.Equal(t, "proxypass", cfg.ZCC.Client.Proxy.Password)
	})

	t.Run("WithPartnerID sets partner ID", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithPartnerID("partner-123")(cfg)
		assert.Equal(t, "partner-123", cfg.ZCC.Client.PartnerID)
	})

	t.Run("WithRequestTimeout sets request timeout", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		timeout := 30 * time.Second
		zcc.WithRequestTimeout(timeout)(cfg)
		assert.Equal(t, timeout, cfg.ZCC.Client.RequestTimeout)
	})

	t.Run("WithRateLimitMaxRetries sets max retries", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithRateLimitMaxRetries(10)(cfg)
		assert.Equal(t, int32(10), cfg.ZCC.Client.RateLimit.MaxRetries)
	})

	t.Run("WithRateLimitMinWait sets min wait", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		minWait := 5 * time.Second
		zcc.WithRateLimitMinWait(minWait)(cfg)
		assert.Equal(t, minWait, cfg.ZCC.Client.RateLimit.RetryWaitMin)
	})

	t.Run("WithRateLimitMaxWait sets max wait", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		maxWait := 60 * time.Second
		zcc.WithRateLimitMaxWait(maxWait)(cfg)
		assert.Equal(t, maxWait, cfg.ZCC.Client.RateLimit.RetryWaitMax)
	})

	t.Run("WithTestingDisableHttpsCheck sets https check", func(t *testing.T) {
		cfg := &zcc.Configuration{}
		zcc.WithTestingDisableHttpsCheck(true)(cfg)
		assert.True(t, cfg.ZCC.Testing.DisableHttpsCheck)
	})
}

// =====================================================
// SetBackoffConfig Test
// =====================================================

func TestConfiguration_SetBackoffConfig(t *testing.T) {
	t.Parallel()

	t.Run("Set backoff configuration", func(t *testing.T) {
		cfg := &zcc.Configuration{}

		backoff := &zcc.BackoffConfig{
			Enabled:             true,
			RetryWaitMinSeconds: 10,
			RetryWaitMaxSeconds: 120,
			MaxNumOfRetries:     5,
		}

		cfg.SetBackoffConfig(backoff)

		assert.NotNil(t, cfg.ZCC.Client.RateLimit.BackoffConf)
		assert.True(t, cfg.ZCC.Client.RateLimit.BackoffConf.Enabled)
		assert.Equal(t, 5, cfg.ZCC.Client.RateLimit.BackoffConf.MaxNumOfRetries)
		assert.Equal(t, 10, cfg.ZCC.Client.RateLimit.BackoffConf.RetryWaitMinSeconds)
		assert.Equal(t, 120, cfg.ZCC.Client.RateLimit.BackoffConf.RetryWaitMaxSeconds)
	})
}

// =====================================================
// Client Tests
// =====================================================

func TestClient_NilConfiguration(t *testing.T) {
	t.Parallel()

	_, err := zcc.NewClient(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration cannot be nil")
}

// =====================================================
// Response Parsing Tests
// =====================================================

func TestConfiguration_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse auth response", func(t *testing.T) {
		jsonResponse := `{
			"token_type": "Bearer",
			"jwtToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			"expires_in": "7200"
		}`

		var token zcc.AuthToken
		err := json.Unmarshal([]byte(jsonResponse), &token)
		require.NoError(t, err)

		assert.Equal(t, "Bearer", token.TokenType)
		assert.NotEmpty(t, token.AccessToken)
		assert.Equal(t, "7200", token.RawExpiresIn)
	})

	t.Run("Parse auth response with empty expires_in", func(t *testing.T) {
		jsonResponse := `{
			"token_type": "Bearer",
			"jwtToken": "test-token"
		}`

		var token zcc.AuthToken
		err := json.Unmarshal([]byte(jsonResponse), &token)
		require.NoError(t, err)

		assert.Equal(t, "Bearer", token.TokenType)
		assert.Equal(t, "test-token", token.AccessToken)
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
		cloud := cloud // capture range variable
		t.Run("Set cloud to "+cloud, func(t *testing.T) {
			cfg := &zcc.Configuration{}
			zcc.WithZCCCloud(cloud)(cfg)
			assert.Equal(t, cloud, cfg.ZCC.Client.ZCCCloud)
		})
	}
}

// =====================================================
// Combined Configuration Test
// =====================================================

func TestConfiguration_MultipleOptions(t *testing.T) {
	t.Parallel()

	t.Run("Apply multiple configuration options", func(t *testing.T) {
		cfg := &zcc.Configuration{}

		// Apply multiple options
		zcc.WithZCCClientID("client-123")(cfg)
		zcc.WithZCCClientSecret("secret-456")(cfg)
		zcc.WithZCCCloud("zscalertwo")(cfg)
		zcc.WithDebug(true)(cfg)
		zcc.WithProxyHost("proxy.corp.com")(cfg)
		zcc.WithProxyPort(3128)(cfg)
		zcc.WithUserAgentExtra("terraform-provider-zscaler/1.0")(cfg)

		// Verify all options were applied
		assert.Equal(t, "client-123", cfg.ZCC.Client.ZCCClientID)
		assert.Equal(t, "secret-456", cfg.ZCC.Client.ZCCClientSecret)
		assert.Equal(t, "zscalertwo", cfg.ZCC.Client.ZCCCloud)
		assert.True(t, cfg.Debug)
		assert.Equal(t, "proxy.corp.com", cfg.ZCC.Client.Proxy.Host)
		assert.Equal(t, int32(3128), cfg.ZCC.Client.Proxy.Port)
		assert.Contains(t, cfg.UserAgentExtra, "terraform-provider-zscaler/1.0")
	})
}
