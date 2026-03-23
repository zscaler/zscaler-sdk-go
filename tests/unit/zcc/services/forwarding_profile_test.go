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

func TestForwardingProfile_GetByCompanyID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webForwardingProfile/listByCompany"

	server.On("GET", path, common.SuccessResponse([]forwarding_profile.ForwardingProfile{
		{ID: 1, Name: "Default Profile", Active: "1"},
		{ID: 2, Name: "Custom Profile", Active: "0"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_profile.GetForwardingProfileByCompanyID(context.Background(), service, "", nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Default Profile", result[0].Name)
	assert.Equal(t, "1", result[0].Active)
}

func TestForwardingProfile_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webForwardingProfile/edit"

	server.On("POST", path, common.SuccessResponse(forwarding_profile.CreateUpdateResponse{
		Success: "true",
		ID:      99,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	req := &forwarding_profile.ForwardingProfileRequest{
		ID:     "-1",
		Name:   "New Profile",
		Active: 1,
	}

	result, err := forwarding_profile.CreateForwardingProfile(context.Background(), service, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "true", result.Success)
	assert.Equal(t, 99, result.ID)
}

func TestForwardingProfile_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webForwardingProfile/99/delete"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = forwarding_profile.DeleteForwardingProfile(context.Background(), service, 99)

	require.NoError(t, err)
}

func TestForwardingProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ForwardingProfile JSON marshaling", func(t *testing.T) {
		profile := forwarding_profile.ForwardingProfile{
			ID:                     123,
			Name:                   "Enterprise Profile",
			Active:                 "1",
			ConditionType:          1,
			DnsServers:             "8.8.8.8",
			DnsSearchDomains:       "corp.example.com",
			EnableLWFDriver:        "1",
			EvaluateTrustedNetwork: 1,
			TrustedGateways:        "192.168.1.1",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":123`)
		assert.Contains(t, string(data), `"name":"Enterprise Profile"`)
		assert.Contains(t, string(data), `"active":"1"`)
		assert.Contains(t, string(data), `"evaluateTrustedNetwork":1`)
	})

	t.Run("ForwardingProfile JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 456,
			"name": "Branch Profile",
			"active": "0",
			"conditionType": 2,
			"dnsServers": "1.1.1.1",
			"enableLWFDriver": "0",
			"evaluateTrustedNetwork": 0
		}`

		var profile forwarding_profile.ForwardingProfile
		err := json.Unmarshal([]byte(jsonData), &profile)
		require.NoError(t, err)

		assert.Equal(t, 456, int(profile.ID))
		assert.Equal(t, "Branch Profile", profile.Name)
		assert.Equal(t, "0", profile.Active)
		assert.Equal(t, 0, profile.EvaluateTrustedNetwork)
	})

	t.Run("CreateUpdateResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{"success": "true", "id": 42065}`

		var resp forwarding_profile.CreateUpdateResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "true", resp.Success)
		assert.Equal(t, 42065, resp.ID)
	})

	t.Run("ForwardingProfileRequest JSON marshaling", func(t *testing.T) {
		req := forwarding_profile.ForwardingProfileRequest{
			ID:                  "-1",
			Active:              1,
			Name:                "Test Profile",
			ConditionType:       1,
			EnableLWFDriver:     1,
			EnableSplitVpnTN:    0,
			PredefinedTnAll:     true,
			TrustedNetworkIds:   []int{},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"-1"`)
		assert.Contains(t, string(data), `"active":1`)
		assert.Contains(t, string(data), `"name":"Test Profile"`)
		assert.Contains(t, string(data), `"predefinedTnAll":true`)
	})
}

func TestForwardingProfile_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse forwarding profiles list", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": "1",
				"name": "Default",
				"active": "1",
				"evaluateTrustedNetwork": 0
			},
			{
				"id": "42065",
				"name": "Custom",
				"active": "0",
				"evaluateTrustedNetwork": 1
			}
		]`

		var profiles []forwarding_profile.ForwardingProfile
		err := json.Unmarshal([]byte(jsonResponse), &profiles)
		require.NoError(t, err)

		assert.Len(t, profiles, 2)
		assert.Equal(t, "Default", profiles[0].Name)
		assert.Equal(t, "1", profiles[0].Active)
		assert.Equal(t, 0, profiles[0].EvaluateTrustedNetwork)
		assert.Equal(t, "Custom", profiles[1].Name)
		assert.Equal(t, "0", profiles[1].Active)
		assert.Equal(t, 1, profiles[1].EvaluateTrustedNetwork)
		assert.Equal(t, 42065, int(profiles[1].ID))
	})
}
