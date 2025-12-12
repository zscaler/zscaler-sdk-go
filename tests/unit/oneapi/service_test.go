// Package oneapi provides unit tests for the OneAPI client
package oneapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func TestService_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NewService creates service with nil client", func(t *testing.T) {
		service := zscaler.NewService(nil, nil)
		assert.NotNil(t, service)
		assert.Nil(t, service.Client)
	})

	t.Run("Service MicroTenantID", func(t *testing.T) {
		service := zscaler.NewService(nil, nil)
		assert.NotNil(t, service)
		// MicroTenantID should be accessible
		_ = service.MicroTenantID()
	})
}

func TestLegacyClient_Structure(t *testing.T) {
	t.Parallel()

	t.Run("LegacyClient can be nil-initialized", func(t *testing.T) {
		legacy := &zscaler.LegacyClient{}
		assert.Nil(t, legacy.ZiaClient)
		assert.Nil(t, legacy.ZpaClient)
		assert.Nil(t, legacy.ZccClient)
		assert.Nil(t, legacy.ZdxClient)
		assert.Nil(t, legacy.ZtwClient)
	})
}

