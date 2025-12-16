// Package unit provides unit tests for ZPA services common utilities
package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestCommon_RemoveCloudSuffix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes cloud suffix with parentheses",
			input:    "CrowdStrike_ZPA_Pre-ZTA (zscalerthree.net)",
			expected: "CrowdStrike_ZPA_Pre-ZTA",
		},
		{
			name:     "removes cloud suffix with different domain",
			input:    "My App (zscaler.net)",
			expected: "My App",
		},
		{
			name:     "removes cloud suffix with underscores and hyphens",
			input:    "Test_App-Name (test-cloud_123.net)",
			expected: "Test_App-Name",
		},
		{
			name:     "no change when no suffix",
			input:    "Simple Name",
			expected: "Simple Name",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "handles string with only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "removes trailing spaces after suffix removal",
			input:    "App Name   (cloud.net)  ",
			expected: "App Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.RemoveCloudSuffix(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCommon_InList(t *testing.T) {
	tests := []struct {
		name     string
		list     []string
		item     string
		expected bool
	}{
		{
			name:     "item exists in list",
			list:     []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "item not in list",
			list:     []string{"apple", "banana", "cherry"},
			item:     "grape",
			expected: false,
		},
		{
			name:     "empty list",
			list:     []string{},
			item:     "apple",
			expected: false,
		},
		{
			name:     "nil list",
			list:     nil,
			item:     "apple",
			expected: false,
		},
		{
			name:     "empty string item exists",
			list:     []string{"", "test"},
			item:     "",
			expected: true,
		},
		{
			name:     "case sensitive match",
			list:     []string{"Apple", "Banana"},
			item:     "apple",
			expected: false,
		},
		{
			name:     "exact case match",
			list:     []string{"Apple", "Banana"},
			item:     "Apple",
			expected: true,
		},
		{
			name:     "single item list - match",
			list:     []string{"only"},
			item:     "only",
			expected: true,
		},
		{
			name:     "single item list - no match",
			list:     []string{"only"},
			item:     "other",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.InList(tt.list, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCommon_Filter_Structure(t *testing.T) {
	t.Run("Filter with all fields", func(t *testing.T) {
		tenantID := "tenant-123"
		filter := common.Filter{
			Search:        "test-search",
			MicroTenantID: &tenantID,
			SortBy:        "name",
			SortOrder:     "ASC",
		}

		assert.Equal(t, "test-search", filter.Search)
		assert.Equal(t, "tenant-123", *filter.MicroTenantID)
		assert.Equal(t, "name", filter.SortBy)
		assert.Equal(t, "ASC", filter.SortOrder)
	})

	t.Run("Empty filter", func(t *testing.T) {
		filter := common.Filter{}

		assert.Empty(t, filter.Search)
		assert.Nil(t, filter.MicroTenantID)
		assert.Empty(t, filter.SortBy)
		assert.Empty(t, filter.SortOrder)
	})
}

func TestCommon_Pagination_Structure(t *testing.T) {
	t.Run("Pagination with all fields", func(t *testing.T) {
		tenantID := "tenant-456"
		pagination := common.Pagination{
			Page:          1,
			PageSize:      100,
			Search:        "test",
			Search2:       "search2",
			MicroTenantID: &tenantID,
			SortBy:        "creationTime",
			SortOrder:     "DESC",
		}

		assert.Equal(t, 1, pagination.Page)
		assert.Equal(t, 100, pagination.PageSize)
		assert.Equal(t, "test", pagination.Search)
		assert.Equal(t, "search2", pagination.Search2)
		assert.Equal(t, "tenant-456", *pagination.MicroTenantID)
		assert.Equal(t, "creationTime", pagination.SortBy)
		assert.Equal(t, "DESC", pagination.SortOrder)
	})
}

