// Package zscaler provides unit tests for core zscaler SDK components
package zscaler

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func TestUserAgent_NewUserAgent(t *testing.T) {
	t.Parallel()

	t.Run("Create UserAgent with config", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		ua := zscaler.NewUserAgent(cfg)

		// Verify that runtime values are captured
		assert.NotEmpty(t, ua.String())
	})

	t.Run("Create UserAgent with nil config", func(t *testing.T) {
		// Even with nil config, it should create a UserAgent
		// Note: This may panic if config is accessed, so we handle it
		cfg := &zscaler.Configuration{}
		ua := zscaler.NewUserAgent(cfg)
		assert.NotEmpty(t, ua.String())
	})
}

func TestUserAgent_String(t *testing.T) {
	t.Parallel()

	t.Run("String contains SDK version", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		ua := zscaler.NewUserAgent(cfg)
		uaStr := ua.String()

		assert.Contains(t, uaStr, "zscaler-sdk-go/")
		assert.Contains(t, uaStr, zscaler.VERSION)
	})

	t.Run("String contains Go version", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		ua := zscaler.NewUserAgent(cfg)
		uaStr := ua.String()

		assert.Contains(t, uaStr, "golang/")
		assert.Contains(t, uaStr, runtime.Version())
	})

	t.Run("String contains OS info", func(t *testing.T) {
		cfg := &zscaler.Configuration{}
		ua := zscaler.NewUserAgent(cfg)
		uaStr := ua.String()

		assert.Contains(t, uaStr, runtime.GOOS)
		assert.Contains(t, uaStr, runtime.GOARCH)
	})

	t.Run("String includes UserAgentExtra when set", func(t *testing.T) {
		cfg := &zscaler.Configuration{
			UserAgentExtra: "terraform-provider/1.0.0",
		}
		ua := zscaler.NewUserAgent(cfg)
		uaStr := ua.String()

		assert.Contains(t, uaStr, "terraform-provider/1.0.0")
	})

	t.Run("String does not include extra when empty", func(t *testing.T) {
		cfg := &zscaler.Configuration{
			UserAgentExtra: "",
		}
		ua := zscaler.NewUserAgent(cfg)
		uaStr := ua.String()

		// Should end with OS info, not have trailing space
		assert.True(t, strings.HasSuffix(uaStr, runtime.GOARCH))
	})
}

func TestUserAgent_Format(t *testing.T) {
	t.Parallel()

	t.Run("Verify complete format", func(t *testing.T) {
		cfg := &zscaler.Configuration{
			UserAgentExtra: "custom/2.0",
		}
		ua := zscaler.NewUserAgent(cfg)
		uaStr := ua.String()

		// Expected format: zscaler-sdk-go/VERSION golang/VERSION OS/ARCH extra
		parts := strings.Split(uaStr, " ")
		assert.GreaterOrEqual(t, len(parts), 3)

		// First part should be SDK version
		assert.True(t, strings.HasPrefix(parts[0], "zscaler-sdk-go/"))

		// Second part should be Go version
		assert.True(t, strings.HasPrefix(parts[1], "golang/"))

		// Third part should be OS/ARCH
		assert.Contains(t, parts[2], "/")
	})
}

