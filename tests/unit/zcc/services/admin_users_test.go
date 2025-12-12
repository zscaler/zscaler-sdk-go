// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/admin_users"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestAdminUsers_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getAdminUsers"

	server.On("GET", path, common.SuccessResponse([]admin_users.AdminUser{
		{ID: 1, UserName: "admin@company.com", AccountEnabled: "true", ServiceType: "ZCC"},
		{ID: 2, UserName: "user@company.com", AccountEnabled: "true", ServiceType: "ZIA"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admin_users.GetAdminUsers(context.Background(), service, "")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "admin@company.com", result[0].UserName)
}

func TestAdminUsers_GetByUserType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getAdminUsers"

	server.On("GET", path, common.SuccessResponse([]admin_users.AdminUser{
		{ID: 1, UserName: "zcc-admin@company.com", ServiceType: "ZCC"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admin_users.GetAdminUsers(context.Background(), service, "ZCC")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "ZCC", result[0].ServiceType)
}

func TestAdminUsers_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/editAdminUser"

	server.On("PUT", path, common.SuccessResponse(admin_users.AdminUser{
		ID:             1,
		UserName:       "admin@company.com",
		AccountEnabled: "true",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateUser := &admin_users.AdminUser{
		ID:             1,
		UserName:       "admin@company.com",
		AccountEnabled: "true",
	}

	result, err := admin_users.UpdateAdminUser(context.Background(), service, updateUser)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
}

func TestAdminUsers_GetSyncZiaZdx_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/syncZiaZdxAdminUsers"

	server.On("POST", path, common.SuccessResponse(admin_users.SyncZiaZdxZpaAdminUsers{
		CompanyIDs:   []int{123456},
		ErrorCode:    "",
		ErrorMessage: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admin_users.GetSyncZiaZdxAdminUsers(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.CompanyIDs, 1)
}

func TestAdminUsers_GetSyncZpa_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/syncZpaAdminUsers"

	server.On("POST", path, common.SuccessResponse(admin_users.SyncZiaZdxZpaAdminUsers{
		CompanyIDs:   []int{123456},
		ErrorCode:    "",
		ErrorMessage: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admin_users.GetSyncZpaAdminUsers(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestAdminUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminUser JSON marshaling", func(t *testing.T) {
		user := admin_users.AdminUser{
			ID:             123,
			UserName:       "admin@company.com",
			AccountEnabled: "true",
			CompanyID:      "company-456",
			EditEnabled:    "true",
			IsDefaultAdmin: "false",
			ServiceType:    "ZCC",
			CompanyRole: admin_users.Role{
				ID:              "role-789",
				RoleName:        "Administrator",
				AdminManagement: "FULL",
				Dashboard:       "FULL",
			},
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":123`)
		assert.Contains(t, string(data), `"userName":"admin@company.com"`)
		assert.Contains(t, string(data), `"accountEnabled":"true"`)
	})

	t.Run("AdminUser JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 456,
			"userName": "user@company.com",
			"accountEnabled": "true",
			"companyId": "company-123",
			"editEnabled": "false",
			"isDefaultAdmin": "true",
			"serviceType": "ZIA",
			"companyRole": {
				"id": "role-001",
				"roleName": "Read Only",
				"adminManagement": "NONE",
				"dashboard": "READ"
			}
		}`

		var user admin_users.AdminUser
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, 456, user.ID)
		assert.Equal(t, "user@company.com", user.UserName)
		assert.Equal(t, "ZIA", user.ServiceType)
		assert.Equal(t, "Read Only", user.CompanyRole.RoleName)
	})

	t.Run("Role JSON marshaling", func(t *testing.T) {
		role := admin_users.Role{
			ID:                           "role-123",
			RoleName:                     "Custom Role",
			AdminManagement:              "FULL",
			Dashboard:                    "FULL",
			DeviceOverview:               "FULL",
			AuditLogs:                    "READ",
			TrustedNetwork:               "FULL",
			ForwardingProfile:            "FULL",
			ClientConnectorAppStore:      "NONE",
			ClientConnectorIDP:           "FULL",
			ClientConnectorSupport:       "FULL",
			ClientConnectorNotifications: "FULL",
			IsEditable:                   true,
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"role-123"`)
		assert.Contains(t, string(data), `"roleName":"Custom Role"`)
		assert.Contains(t, string(data), `"isEditable":true`)
	})
}

func TestAdminUsers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse admin users list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"userName": "admin1@company.com",
				"accountEnabled": "true",
				"serviceType": "ZCC"
			},
			{
				"id": 2,
				"userName": "admin2@company.com",
				"accountEnabled": "false",
				"serviceType": "ZIA"
			}
		]`

		var users []admin_users.AdminUser
		err := json.Unmarshal([]byte(jsonResponse), &users)
		require.NoError(t, err)

		assert.Len(t, users, 2)
		assert.Equal(t, "admin1@company.com", users[0].UserName)
		assert.Equal(t, "true", users[0].AccountEnabled)
		assert.Equal(t, "ZIA", users[1].ServiceType)
	})
}
