// Package services provides unit tests for ZTW services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/activation"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestActivation_GetActivationStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ecAdminActivateStatus"

	server.On("GET", path, common.SuccessResponse(activation.ECAdminActivation{
		OrgEditStatus:         "ACTIVE",
		OrgLastActivateStatus: "SUCCESS",
		AdminActivateStatus:   "READY",
		AdminStatusMap: map[string]interface{}{
			"admin@company.com": "active",
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := activation.GetActivationStatus(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "ACTIVE", result.OrgEditStatus)
	assert.Equal(t, "SUCCESS", result.OrgLastActivateStatus)
}

func TestActivation_UpdateActivationStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ecAdminActivateStatus/activate"

	server.On("PUT", path, common.SuccessResponse(activation.ECAdminActivation{
		OrgEditStatus:         "PENDING",
		OrgLastActivateStatus: "ACTIVATION_IN_PROGRESS",
		AdminActivateStatus:   "ACTIVATING",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateReq := activation.ECAdminActivation{
		AdminActivateStatus: "ACTIVATING",
	}

	result, err := activation.UpdateActivationStatus(context.Background(), service, updateReq)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "ACTIVATION_IN_PROGRESS", result.OrgLastActivateStatus)
}

func TestActivation_ForceActivationStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ecAdminActivateStatus/forcedActivate"

	server.On("PUT", path, common.SuccessResponse(activation.ECAdminActivation{
		OrgEditStatus:         "FORCE_ACTIVATED",
		OrgLastActivateStatus: "FORCED_SUCCESS",
		AdminActivateStatus:   "FORCE_ACTIVATED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	forceReq := activation.ECAdminActivation{
		AdminActivateStatus: "FORCE_ACTIVATE",
	}

	result, err := activation.ForceActivationStatus(context.Background(), service, forceReq)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "FORCE_ACTIVATED", result.OrgEditStatus)
}

// =====================================================
// Structure Tests
// =====================================================

func TestActivation_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ECAdminActivation JSON marshaling", func(t *testing.T) {
		act := activation.ECAdminActivation{
			OrgEditStatus:         "ACTIVE",
			OrgLastActivateStatus: "SUCCESS",
			AdminActivateStatus:   "PENDING",
			AdminStatusMap: map[string]interface{}{
				"admin1": "active",
				"admin2": "pending",
			},
		}

		data, err := json.Marshal(act)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"orgEditStatus":"ACTIVE"`)
		assert.Contains(t, string(data), `"orgLastActivateStatus":"SUCCESS"`)
		assert.Contains(t, string(data), `"adminActivateStatus":"PENDING"`)
		assert.Contains(t, string(data), `"adminStatusMap"`)
	})

	t.Run("ECAdminActivation JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"orgEditStatus": "MODIFIED",
			"orgLastActivateStatus": "FAILED",
			"adminActivateStatus": "INACTIVE",
			"adminStatusMap": {
				"admin1": "inactive",
				"admin2": "active"
			}
		}`

		var act activation.ECAdminActivation
		err := json.Unmarshal([]byte(jsonData), &act)
		require.NoError(t, err)

		assert.Equal(t, "MODIFIED", act.OrgEditStatus)
		assert.Equal(t, "FAILED", act.OrgLastActivateStatus)
		assert.Equal(t, "INACTIVE", act.AdminActivateStatus)
		assert.NotNil(t, act.AdminStatusMap)
		assert.Equal(t, "inactive", act.AdminStatusMap["admin1"])
	})

	t.Run("ECAdminActivation empty status map", func(t *testing.T) {
		act := activation.ECAdminActivation{
			OrgEditStatus:       "NONE",
			AdminActivateStatus: "NONE",
		}

		data, err := json.Marshal(act)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"orgEditStatus":"NONE"`)
	})
}

func TestActivation_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse activation status response", func(t *testing.T) {
		jsonResponse := `{
			"orgEditStatus": "PENDING_ACTIVATION",
			"orgLastActivateStatus": "ACTIVATION_IN_PROGRESS",
			"adminActivateStatus": "READY",
			"adminStatusMap": {
				"super_admin@company.com": "ready",
				"admin@company.com": "pending"
			}
		}`

		var act activation.ECAdminActivation
		err := json.Unmarshal([]byte(jsonResponse), &act)
		require.NoError(t, err)

		assert.Equal(t, "PENDING_ACTIVATION", act.OrgEditStatus)
		assert.Equal(t, "ACTIVATION_IN_PROGRESS", act.OrgLastActivateStatus)
		assert.Equal(t, "READY", act.AdminActivateStatus)
		assert.Len(t, act.AdminStatusMap, 2)
	})
}

