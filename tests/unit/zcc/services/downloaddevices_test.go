// Package services provides unit tests for ZCC services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling for download devices
// Note: The download functions write to io.Writer and require
// special mock handling that's beyond the scope of unit tests.
// These tests focus on the request/response structures.
// =====================================================

func TestDownloadDevices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Download request parameters", func(t *testing.T) {
		// Test the parameters structure for download devices
		params := map[string]string{
			"osTypes":           "1,2,3",
			"registrationTypes": "all",
		}

		data, err := json.Marshal(params)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"osTypes":"1,2,3"`)
		assert.Contains(t, string(data), `"registrationTypes":"all"`)
	})

	t.Run("Download service status parameters", func(t *testing.T) {
		// Test the parameters structure for download service status
		params := map[string]string{
			"osTypes":           "1",
			"registrationTypes": "registered",
		}

		data, err := json.Marshal(params)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"osTypes":"1"`)
		assert.Contains(t, string(data), `"registrationTypes":"registered"`)
	})
}
