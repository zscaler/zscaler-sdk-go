// Package zscaler provides unit tests for core zscaler SDK service
package zscaler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func TestService_Constants(t *testing.T) {
	t.Parallel()

	t.Run("SortOrder constants", func(t *testing.T) {
		assert.Equal(t, "asc", string(zscaler.ASCSortOrder))
		assert.Equal(t, "desc", string(zscaler.DESCSortOrder))
	})

	t.Run("SortField constants", func(t *testing.T) {
		assert.Equal(t, "id", string(zscaler.IDSortField))
		assert.Equal(t, "name", string(zscaler.NameSortField))
		assert.Equal(t, "creationTime", string(zscaler.CreationTimeSortField))
		assert.Equal(t, "modifiedTime", string(zscaler.ModifiedTimeSortField))
	})
}

func TestService_NewService(t *testing.T) {
	t.Parallel()

	t.Run("Create service with nil clients", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)

		require.NotNil(t, svc)
		assert.Nil(t, svc.Client)
		assert.Nil(t, svc.LegacyClient)
		// Verify default values using string comparison
		assert.Equal(t, "asc", string(svc.SortOrder))
		assert.Equal(t, "name", string(svc.SortBy))
	})

	t.Run("Create service with default sort values", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)

		assert.Equal(t, "asc", string(svc.SortOrder))
		assert.Equal(t, "name", string(svc.SortBy))
	})
}

func TestService_WithMicroTenant(t *testing.T) {
	t.Parallel()

	t.Run("Set micro tenant ID", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		microTenantSvc := svc.WithMicroTenant("tenant-12345")

		require.NotNil(t, microTenantSvc)
		require.NotNil(t, microTenantSvc.MicroTenantID())
		assert.Equal(t, "tenant-12345", *microTenantSvc.MicroTenantID())
	})

	t.Run("Set empty micro tenant ID", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		microTenantSvc := svc.WithMicroTenant("")

		require.NotNil(t, microTenantSvc)
		assert.Nil(t, microTenantSvc.MicroTenantID())
	})

	t.Run("Original service unchanged", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		_ = svc.WithMicroTenant("tenant-xyz")

		// Original service should still have nil micro tenant
		assert.Nil(t, svc.MicroTenantID())
	})
}

func TestService_MicroTenantID(t *testing.T) {
	t.Parallel()

	t.Run("Get nil micro tenant ID", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)

		assert.Nil(t, svc.MicroTenantID())
	})

	t.Run("Get set micro tenant ID", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		microTenantSvc := svc.WithMicroTenant("mt-123")

		mid := microTenantSvc.MicroTenantID()
		require.NotNil(t, mid)
		assert.Equal(t, "mt-123", *mid)
	})
}

func TestService_WithSort(t *testing.T) {
	t.Parallel()

	t.Run("Set sort by ID ascending", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		sortedSvc := svc.WithSort(zscaler.IDSortField, zscaler.ASCSortOrder)

		require.NotNil(t, sortedSvc)
		assert.Equal(t, "id", string(sortedSvc.SortBy))
		assert.Equal(t, "asc", string(sortedSvc.SortOrder))
	})

	t.Run("Set sort by name descending", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		sortedSvc := svc.WithSort(zscaler.NameSortField, zscaler.DESCSortOrder)

		assert.Equal(t, "name", string(sortedSvc.SortBy))
		assert.Equal(t, "desc", string(sortedSvc.SortOrder))
	})

	t.Run("Set sort by creationTime", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		sortedSvc := svc.WithSort(zscaler.CreationTimeSortField, zscaler.ASCSortOrder)

		assert.Equal(t, "creationTime", string(sortedSvc.SortBy))
	})

	t.Run("Set sort by modifiedTime", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		sortedSvc := svc.WithSort(zscaler.ModifiedTimeSortField, zscaler.DESCSortOrder)

		assert.Equal(t, "modifiedTime", string(sortedSvc.SortBy))
	})

	t.Run("Invalid sort field keeps default", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		sortedSvc := svc.WithSort(zscaler.SortField("invalid"), zscaler.ASCSortOrder)

		// Should keep the default since "invalid" is not a valid field
		assert.Equal(t, "name", string(sortedSvc.SortBy))
	})

	t.Run("Invalid sort order keeps default", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		sortedSvc := svc.WithSort(zscaler.IDSortField, zscaler.SortOrder("invalid"))

		// Should keep the default since "invalid" is not a valid order
		assert.Equal(t, "asc", string(sortedSvc.SortOrder))
		// But sort field should be set
		assert.Equal(t, "id", string(sortedSvc.SortBy))
	})

	t.Run("Original service unchanged", func(t *testing.T) {
		svc := zscaler.NewService(nil, nil)
		_ = svc.WithSort(zscaler.IDSortField, zscaler.DESCSortOrder)

		// Original should still have defaults
		assert.Equal(t, "name", string(svc.SortBy))
		assert.Equal(t, "asc", string(svc.SortOrder))
	})
}

func TestService_NewZPAScimService(t *testing.T) {
	t.Parallel()

	t.Run("Nil config returns nil", func(t *testing.T) {
		svc := zscaler.NewZPAScimService(nil)
		assert.Nil(t, svc)
	})
}

func TestService_NewZIAScimService(t *testing.T) {
	t.Parallel()

	t.Run("Nil config returns nil", func(t *testing.T) {
		svc := zscaler.NewZIAScimService(nil)
		assert.Nil(t, svc)
	})
}

