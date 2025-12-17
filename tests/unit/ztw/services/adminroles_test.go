// Package services provides unit tests for ZTW services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/adminuserrolemgmt/adminroles"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestAdminRoles_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := 12345
	path := "/ztw/api/v1/adminRoles/12345"

	server.On("GET", path, common.SuccessResponse(adminroles.AdminRoles{
		ID:           roleID,
		Name:         "Super Admin",
		Rank:         1,
		PolicyAccess: "READ_WRITE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adminroles.Get(context.Background(), service, roleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, roleID, result.ID)
	assert.Equal(t, "Super Admin", result.Name)
}

func TestAdminRoles_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleName := "Super Admin"
	path := "/ztw/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]adminroles.AdminRoles{
		{ID: 1, Name: "Other Role", Rank: 2},
		{ID: 2, Name: roleName, Rank: 1},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adminroles.GetByName(context.Background(), service, roleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, roleName, result.Name)
}

func TestAdminRoles_GetAllAdminRoles_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]adminroles.AdminRoles{
		{ID: 1, Name: "Role 1", Rank: 1},
		{ID: 2, Name: "Role 2", Rank: 2},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adminroles.GetAllAdminRoles(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestAdminRoles_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/adminRoles"

	server.On("POST", path, common.SuccessResponse(adminroles.AdminRoles{
		ID:           99999,
		Name:         "New Role",
		Rank:         3,
		PolicyAccess: "READ_ONLY",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRole := &adminroles.AdminRoles{
		Name:         "New Role",
		Rank:         3,
		PolicyAccess: "READ_ONLY",
	}

	result, err := adminroles.Create(context.Background(), service, newRole)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestAdminRoles_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := 12345
	path := "/ztw/api/v1/adminRoles/12345"

	server.On("PUT", path, common.SuccessResponse(adminroles.AdminRoles{
		ID:   roleID,
		Name: "Updated Role",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRole := &adminroles.AdminRoles{
		ID:   roleID,
		Name: "Updated Role",
	}

	result, err := adminroles.Update(context.Background(), service, roleID, updateRole)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Role", result.Name)
}

func TestAdminRoles_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := 12345
	path := "/ztw/api/v1/adminRoles/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = adminroles.Delete(context.Background(), service, roleID)

	require.NoError(t, err)
}

func TestAdminRoles_GetAPIRole_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	apiRoleName := "API Admin"
	path := "/ztw/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]adminroles.AdminRoles{
		{ID: 1, Name: "Other Role", Rank: 2},
		{ID: 2, Name: apiRoleName, Rank: 1},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adminroles.GetAPIRole(context.Background(), service, apiRoleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, apiRoleName, result.Name)
}

func TestAdminRoles_GetAuditorRole_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	auditorRoleName := "Auditor"
	path := "/ztw/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]adminroles.AdminRoles{
		{ID: 1, Name: "Other Role", Rank: 2, IsAuditor: false},
		{ID: 2, Name: auditorRoleName, Rank: 1, IsAuditor: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adminroles.GetAuditorRole(context.Background(), service, auditorRoleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, auditorRoleName, result.Name)
}

func TestAdminRoles_GetPartnerRole_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	partnerRoleName := "Partner Admin"
	path := "/ztw/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]adminroles.AdminRoles{
		{ID: 1, Name: "Other Role", Rank: 2},
		{ID: 2, Name: partnerRoleName, Rank: 1},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adminroles.GetPartnerRole(context.Background(), service, partnerRoleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, partnerRoleName, result.Name)
}

// =====================================================
// Structure Tests
// =====================================================

func TestAdminRoles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminRoles JSON marshaling", func(t *testing.T) {
		role := adminroles.AdminRoles{
			ID:               1234,
			Rank:             1,
			Name:             "Super Admin",
			PolicyAccess:     "READ_WRITE",
			AlertingAccess:   "READ_WRITE",
			DashboardAccess:  "READ_ONLY",
			ReportAccess:     "READ_WRITE",
			AnalysisAccess:   "READ_ONLY",
			UsernameAccess:   "READ_ONLY",
			AdminAcctAccess:  "READ_WRITE",
			DeviceInfoAccess: "READ_ONLY",
			IsAuditor:        false,
			Permissions:      []string{"POLICY_MANAGE", "USER_MANAGE", "REPORT_VIEW"},
			IsNonEditable:    true,
			LogsLimit:        "UNLIMITED",
			RoleType:         "EXEC_INSIGHT_AND_ORG_ADMIN",
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":1234`)
		assert.Contains(t, string(data), `"name":"Super Admin"`)
		assert.Contains(t, string(data), `"policyAccess":"READ_WRITE"`)
		assert.Contains(t, string(data), `"isNonEditable":true`)
	})

	t.Run("AdminRoles JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 5678,
			"rank": 3,
			"name": "Read Only Admin",
			"policyAccess": "READ_ONLY",
			"alertingAccess": "READ_ONLY",
			"dashboardAccess": "READ_ONLY",
			"reportAccess": "READ_ONLY",
			"analysisAccess": "READ_ONLY",
			"usernameAccess": "NONE",
			"adminAcctAccess": "NONE",
			"deviceInfoAccess": "READ_ONLY",
			"isAuditor": true,
			"permissions": ["REPORT_VIEW", "DASHBOARD_VIEW"],
			"isNonEditable": false,
			"logsLimit": "1_WEEK",
			"roleType": "AUDITOR"
		}`

		var role adminroles.AdminRoles
		err := json.Unmarshal([]byte(jsonData), &role)
		require.NoError(t, err)

		assert.Equal(t, 5678, role.ID)
		assert.Equal(t, 3, role.Rank)
		assert.Equal(t, "Read Only Admin", role.Name)
		assert.Equal(t, "READ_ONLY", role.PolicyAccess)
		assert.True(t, role.IsAuditor)
		assert.Len(t, role.Permissions, 2)
		assert.Equal(t, "AUDITOR", role.RoleType)
	})

	t.Run("AdminRoles with feature permissions", func(t *testing.T) {
		jsonData := `{
			"id": 9999,
			"name": "Custom Admin",
			"featurePermissions": {
				"firewall": "READ_WRITE",
				"dlp": "READ_ONLY",
				"sandbox": "NONE"
			}
		}`

		var role adminroles.AdminRoles
		err := json.Unmarshal([]byte(jsonData), &role)
		require.NoError(t, err)

		assert.NotNil(t, role.FeaturePermissions)
		assert.Equal(t, "READ_WRITE", role.FeaturePermissions["firewall"])
		assert.Equal(t, "READ_ONLY", role.FeaturePermissions["dlp"])
	})
}

func TestAdminRoles_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse admin roles list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "Super Admin",
				"rank": 1,
				"policyAccess": "READ_WRITE",
				"isNonEditable": true
			},
			{
				"id": 2,
				"name": "Security Admin",
				"rank": 2,
				"policyAccess": "READ_WRITE",
				"isNonEditable": false
			},
			{
				"id": 3,
				"name": "Auditor",
				"rank": 7,
				"isAuditor": true,
				"isNonEditable": true
			}
		]`

		var roles []adminroles.AdminRoles
		err := json.Unmarshal([]byte(jsonResponse), &roles)
		require.NoError(t, err)

		assert.Len(t, roles, 3)
		assert.Equal(t, "Super Admin", roles[0].Name)
		assert.Equal(t, 1, roles[0].Rank)
		assert.True(t, roles[2].IsAuditor)
	})
}

