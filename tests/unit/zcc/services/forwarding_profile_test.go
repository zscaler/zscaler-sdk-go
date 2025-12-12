// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/forwarding_profile"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestForwardingProfile_GetByCompanyID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webForwardingProfile"

	server.On("GET", path, common.SuccessResponse([]forwarding_profile.ForwardingProfile{
		{ID: 1, Name: "Default Profile", Active: "true"},
		{ID: 2, Name: "Custom Profile", Active: "false"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_profile.GetForwardingProfileByCompanyID(context.Background(), service, "", nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Default Profile", result[0].Name)
}

func TestForwardingProfile_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webForwardingProfile"

	server.On("POST", path, common.SuccessResponse(forwarding_profile.ForwardingProfile{
		ID:     99,
		Name:   "New Profile",
		Active: "true",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newProfile := &forwarding_profile.ForwardingProfile{
		Name:   "New Profile",
		Active: "true",
	}

	result, err := forwarding_profile.CreateForwardingProfile(context.Background(), service, newProfile)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "New Profile", result.Name)
}

func TestForwardingProfile_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webForwardingProfile/99"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = forwarding_profile.DeleteForwardingProfile(context.Background(), service, 99)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestForwardingProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ForwardingProfile JSON marshaling", func(t *testing.T) {
		profile := forwarding_profile.ForwardingProfile{
			ID:                   123,
			Name:                 "Enterprise Profile",
			Active:               "true",
			ConditionType:        1,
			DnsServers:           "8.8.8.8",
			DnsSearchDomains:     "corp.example.com",
			EnableLWFDriver:      "true",
			EvaluateTrustedNetwork: 1,
			TrustedGateways:      "192.168.1.1",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":123`)
		assert.Contains(t, string(data), `"name":"Enterprise Profile"`)
		assert.Contains(t, string(data), `"active":"true"`)
	})

	t.Run("ForwardingProfile JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 456,
			"name": "Branch Profile",
			"active": "false",
			"conditionType": 2,
			"dnsServers": "1.1.1.1",
			"enableLWFDriver": "false",
			"evaluateTrustedNetwork": 0
		}`

		var profile forwarding_profile.ForwardingProfile
		err := json.Unmarshal([]byte(jsonData), &profile)
		require.NoError(t, err)

		assert.Equal(t, 456, int(profile.ID))
		assert.Equal(t, "Branch Profile", profile.Name)
		assert.Equal(t, "false", profile.Active)
	})
}

func TestForwardingProfile_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse forwarding profiles list", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "Default",
				"active": "true"
			},
			{
				"id": 2,
				"name": "Custom",
				"active": "false"
			}
		]`

		var profiles []forwarding_profile.ForwardingProfile
		err := json.Unmarshal([]byte(jsonResponse), &profiles)
		require.NoError(t, err)

		assert.Len(t, profiles, 2)
		assert.Equal(t, "Default", profiles[0].Name)
		assert.Equal(t, "true", profiles[0].Active)
		assert.Equal(t, "Custom", profiles[1].Name)
		assert.Equal(t, "false", profiles[1].Active)
	})
}
