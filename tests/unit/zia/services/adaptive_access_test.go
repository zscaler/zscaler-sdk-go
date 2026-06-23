// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adaptive_access"
)

const (
	adaptiveAccessPath      = "/zia/api/v1/adaptiveAccessProfiles"
	adaptiveAccessRulesPath = "/zia/api/v1/adaptiveAccessProfiles/profiles/rules"
)

// =====================================================
// GetProfileRules
// =====================================================

func TestAdaptiveAccess_GetProfileRules_NoOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", adaptiveAccessRulesPath, func(r *http.Request, _ []byte) common.MockResponse {
		// Without options there should be no query string.
		assert.Empty(t, r.URL.RawQuery)
		return common.SuccessResponse([]adaptive_access.AdaptiveAccess{
			{ID: 1, Name: "Profile A", Type: "USER", AapIndex: 0, IamAapID: "aap-1"},
			{ID: 2, Name: "Profile B", Type: "DEVICE", AapIndex: 1, IamAapID: "aap-2"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetProfileRules(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, "Profile A", result[0].Name)
	assert.Equal(t, "aap-1", result[0].IamAapID)
}

func TestAdaptiveAccess_GetProfileRules_EmptyOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", adaptiveAccessRulesPath, func(r *http.Request, _ []byte) common.MockResponse {
		// An empty (but non-nil) options struct should not add any query params.
		assert.Empty(t, r.URL.RawQuery)
		return common.SuccessResponse([]adaptive_access.AdaptiveAccess{
			{ID: 10, Name: "Profile C"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetProfileRules(context.Background(), service, &adaptive_access.GetFilterOptions{})

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 10, result[0].ID)
}

func TestAdaptiveAccess_GetProfileRules_WithIAMAapIDs_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", adaptiveAccessRulesPath, func(r *http.Request, _ []byte) common.MockResponse {
		// Each ID is sent as a repeated query parameter.
		ids := r.URL.Query()["iamAapIds"]
		assert.ElementsMatch(t, []string{"aap-1", "aap-2"}, ids)
		assert.Empty(t, r.URL.Query().Get("orgId"))
		return common.SuccessResponse([]adaptive_access.AdaptiveAccess{
			{ID: 1, Name: "Profile A", IamAapID: "aap-1"},
			{ID: 2, Name: "Profile B", IamAapID: "aap-2"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &adaptive_access.GetFilterOptions{
		IAMAapIDs: []string{"aap-1", "aap-2"},
	}
	result, err := adaptive_access.GetProfileRules(context.Background(), service, opts)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestAdaptiveAccess_GetProfileRules_WithOrgID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", adaptiveAccessRulesPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Equal(t, "98765", r.URL.Query().Get("orgId"))
		assert.Empty(t, r.URL.Query()["iamAapIds"])
		return common.SuccessResponse([]adaptive_access.AdaptiveAccess{
			{ID: 5, Name: "Org Scoped Profile"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	orgID := 98765
	opts := &adaptive_access.GetFilterOptions{OrgID: &orgID}
	result, err := adaptive_access.GetProfileRules(context.Background(), service, opts)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 5, result[0].ID)
}

func TestAdaptiveAccess_GetProfileRules_WithAllFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", adaptiveAccessRulesPath, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.ElementsMatch(t, []string{"aap-1", "aap-2"}, q["iamAapIds"])
		assert.Equal(t, "42", q.Get("orgId"))
		return common.SuccessResponse([]adaptive_access.AdaptiveAccess{
			{ID: 1, Name: "Profile A", IamAapID: "aap-1"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	orgID := 42
	opts := &adaptive_access.GetFilterOptions{
		IAMAapIDs: []string{"aap-1", "aap-2"},
		OrgID:     &orgID,
	}
	result, err := adaptive_access.GetProfileRules(context.Background(), service, opts)

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestAdaptiveAccess_GetProfileRules_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", adaptiveAccessRulesPath, common.SuccessResponse([]adaptive_access.AdaptiveAccess{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetProfileRules(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestAdaptiveAccess_GetProfileRules_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", adaptiveAccessRulesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetProfileRules(context.Background(), service, nil)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// GetAll
// =====================================================

func TestAdaptiveAccess_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", adaptiveAccessPath, common.SuccessResponse([]adaptive_access.AdaptiveAccess{
		{ID: 1, Name: "Profile A"},
		{ID: 2, Name: "Profile B"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// GetByName
// =====================================================

func TestAdaptiveAccess_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", adaptiveAccessPath, common.SuccessResponse([]adaptive_access.AdaptiveAccess{
		{ID: 1, Name: "Profile A"},
		{ID: 2, Name: "Profile B"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetByName(context.Background(), service, "Profile B")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
	assert.Equal(t, "Profile B", result.Name)
}

func TestAdaptiveAccess_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", adaptiveAccessPath, common.SuccessResponse([]adaptive_access.AdaptiveAccess{
		{ID: 1, Name: "Profile A"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetByName(context.Background(), service, "profile a")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
}

func TestAdaptiveAccess_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", adaptiveAccessPath, common.SuccessResponse([]adaptive_access.AdaptiveAccess{
		{ID: 1, Name: "Profile A"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetByName(context.Background(), service, "NonExistent")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestAdaptiveAccess_GetByName_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", adaptiveAccessPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adaptive_access.GetByName(context.Background(), service, "Profile A")

	require.Error(t, err)
	assert.Nil(t, result)
}
