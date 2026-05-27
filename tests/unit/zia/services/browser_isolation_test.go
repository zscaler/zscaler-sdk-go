// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/browser_isolation"
)

const cbiProfilesPath = "/zia/api/v1/browserIsolation/profiles"

// =====================================================
// SDK Function Tests
// =====================================================

func TestBrowserIsolation_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileName := "Corporate Isolation"
	server.On("GET", cbiProfilesPath, common.SuccessResponse([]browser_isolation.CBIProfile{
		{ID: "1", Name: "Other Profile", URL: "https://other.example.com"},
		{ID: "2", Name: profileName, URL: "https://corp.example.com", DefaultProfile: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetByName(context.Background(), service, profileName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileName, result.Name)
	assert.Equal(t, "2", result.ID)
	assert.True(t, result.DefaultProfile)
}

func TestBrowserIsolation_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cbiProfilesPath, common.SuccessResponse([]browser_isolation.CBIProfile{
		{ID: "abc", Name: "My Isolation Profile", URL: "https://iso.example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetByName(context.Background(), service, "MY ISOLATION PROFILE")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "My Isolation Profile", result.Name)
}

func TestBrowserIsolation_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cbiProfilesPath, common.SuccessResponse([]browser_isolation.CBIProfile{
		{ID: "1", Name: "Existing Profile"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetByName(context.Background(), service, "Missing Profile")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no cloud browser isolation profile found with name: Missing Profile")
}

func TestBrowserIsolation_GetByName_APIError_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cbiProfilesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetByName(context.Background(), service, "Any")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestBrowserIsolation_GetByName_NotSubscribed_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cbiProfilesPath, common.MockResponse{
		StatusCode: http.StatusForbidden,
		Body:       `{"message": "Cloud Browser Isolation subscription is required for this feature"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetByName(context.Background(), service, "Any")

	require.Error(t, err)
	assert.Nil(t, result)

	var notSubscribed *browser_isolation.NotSubscribedError
	require.True(t, errors.As(err, &notSubscribed))
	assert.Contains(t, err.Error(), "NOT_SUBSCRIBED")
	assert.Contains(t, err.Error(), "Cloud Browser Isolation subscription is required")
}

func TestBrowserIsolation_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cbiProfilesPath, common.SuccessResponse([]browser_isolation.CBIProfile{
		{ID: "1", Name: "Profile A", URL: "https://a.example.com"},
		{ID: "2", Name: "Profile B", URL: "https://b.example.com", DefaultProfile: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Profile A", result[0].Name)
}

func TestBrowserIsolation_GetAll_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cbiProfilesPath, common.SuccessResponse([]browser_isolation.CBIProfile{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestBrowserIsolation_GetAll_APIError_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cbiProfilesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetAll(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestBrowserIsolation_GetAll_NotSubscribed_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cbiProfilesPath, common.MockResponse{
		StatusCode: http.StatusForbidden,
		Body:       `{"message": "Cloud Browser Isolation subscription is required"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_isolation.GetAll(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)

	var notSubscribed *browser_isolation.NotSubscribedError
	require.True(t, errors.As(err, &notSubscribed))
}

// =====================================================
// Structure Tests
// =====================================================

func TestBrowserIsolation_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CBIProfile JSON marshaling", func(t *testing.T) {
		profile := browser_isolation.CBIProfile{
			ID:             "uuid-123",
			Name:           "Isolation Profile",
			URL:            "https://isolation.zscaler.com",
			DefaultProfile: true,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"uuid-123"`)
		assert.Contains(t, string(data), `"name":"Isolation Profile"`)
		assert.Contains(t, string(data), `"url":"https://isolation.zscaler.com"`)
		assert.Contains(t, string(data), `"defaultProfile":true`)
	})

	t.Run("CBIProfile JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "profile-id",
			"name": "Default CBI",
			"url": "https://cbi.example.com",
			"defaultProfile": false
		}`

		var profile browser_isolation.CBIProfile
		err := json.Unmarshal([]byte(jsonData), &profile)
		require.NoError(t, err)

		assert.Equal(t, "profile-id", profile.ID)
		assert.Equal(t, "Default CBI", profile.Name)
		assert.Equal(t, "https://cbi.example.com", profile.URL)
		assert.False(t, profile.DefaultProfile)
	})

	t.Run("CBIProfile minimal fields", func(t *testing.T) {
		jsonData := `{"name": "Name Only"}`

		var profile browser_isolation.CBIProfile
		err := json.Unmarshal([]byte(jsonData), &profile)
		require.NoError(t, err)

		assert.Equal(t, "Name Only", profile.Name)
		assert.Empty(t, profile.ID)
	})
}

func TestBrowserIsolation_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse CBI profiles list", func(t *testing.T) {
		jsonResponse := `[
			{"id": "1", "name": "Profile 1", "url": "https://1.example.com"},
			{"id": "2", "name": "Profile 2", "url": "https://2.example.com", "defaultProfile": true}
		]`

		var profiles []browser_isolation.CBIProfile
		err := json.Unmarshal([]byte(jsonResponse), &profiles)
		require.NoError(t, err)

		assert.Len(t, profiles, 2)
		assert.True(t, profiles[1].DefaultProfile)
	})
}
