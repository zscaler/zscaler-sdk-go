// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/manage_pass"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestManagePass_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/managePass"

	server.On("POST", path, common.SuccessResponse(manage_pass.ManagePassResponseContract{
		ErrorMessage: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updatePass := &manage_pass.ManagePass{
		CompanyID:     123456,
		DeviceType:    1,
		ExitPass:      "password123",
		LogoutPass:    "logout123",
		UninstallPass: "uninstall123",
		PolicyName:    "Default Policy",
	}

	result, err := manage_pass.UpdateManagePass(context.Background(), service, updatePass)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.ErrorMessage)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestManagePass_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ManagePass JSON marshaling", func(t *testing.T) {
		pass := manage_pass.ManagePass{
			CompanyID:      123456,
			DeviceType:     1,
			ExitPass:       "exit-password",
			LogoutPass:     "logout-password",
			UninstallPass:  "uninstall-password",
			PolicyName:     "Enterprise Policy",
			ZiaDisablePass: "zia-disable",
			ZpaDisablePass: "zpa-disable",
			ZdxDisablePass: "zdx-disable",
			ZdpDisablePass: "zdp-disable",
			ZadDisablePass: "zad-disable",
		}

		data, err := json.Marshal(pass)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"companyId":123456`)
		assert.Contains(t, string(data), `"deviceType":1`)
		assert.Contains(t, string(data), `"exitPass":"exit-password"`)
		assert.Contains(t, string(data), `"policyName":"Enterprise Policy"`)
	})

	t.Run("ManagePass JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"companyId": 789012,
			"deviceType": 2,
			"exitPass": "exit123",
			"logoutPass": "logout123",
			"uninstallPass": "uninstall123",
			"policyName": "Custom Policy",
			"ziaDisablePass": "zia123",
			"zpaDisablePass": "zpa123"
		}`

		var pass manage_pass.ManagePass
		err := json.Unmarshal([]byte(jsonData), &pass)
		require.NoError(t, err)

		assert.Equal(t, 789012, pass.CompanyID)
		assert.Equal(t, 2, pass.DeviceType)
		assert.Equal(t, "exit123", pass.ExitPass)
		assert.Equal(t, "Custom Policy", pass.PolicyName)
	})

	t.Run("ManagePassResponseContract JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"errorMessage": ""
		}`

		var response manage_pass.ManagePassResponseContract
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Empty(t, response.ErrorMessage)
	})

	t.Run("ManagePassResponseContract with error", func(t *testing.T) {
		jsonData := `{
			"errorMessage": "Invalid password format"
		}`

		var response manage_pass.ManagePassResponseContract
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, "Invalid password format", response.ErrorMessage)
	})
}
