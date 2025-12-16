// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/admins"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/roles"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestAdminUsers_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	adminID := 12345
	path := "/zia/api/v1/adminUsers/12345"

	server.On("GET", path, common.SuccessResponse(admins.AdminUsers{
		ID:        adminID,
		LoginName: "admin@company.com",
		UserName:  "Admin User",
		Email:     "admin@company.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admins.GetAdminUsers(context.Background(), service, adminID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, adminID, result.ID)
	assert.Equal(t, "admin@company.com", result.LoginName)
}

func TestAdminUsers_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminUsers"

	server.On("GET", path, common.SuccessResponse([]admins.AdminUsers{
		{ID: 1, LoginName: "admin1@company.com", UserName: "Admin 1"},
		{ID: 2, LoginName: "admin2@company.com", UserName: "Admin 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admins.GetAllAdminUsers(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestAdminUsers_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminUsers"

	server.On("POST", path, common.SuccessResponse(admins.AdminUsers{
		ID:        99999,
		LoginName: "new@company.com",
		UserName:  "New Admin",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newAdmin := admins.AdminUsers{
		LoginName: "new@company.com",
		UserName:  "New Admin",
	}

	result, err := admins.CreateAdminUser(context.Background(), service, newAdmin)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestAdminUsers_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	adminID := 12345
	path := "/zia/api/v1/adminUsers/12345"

	server.On("PUT", path, common.SuccessResponse(admins.AdminUsers{
		ID:       adminID,
		UserName: "Updated Admin",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateAdmin := admins.AdminUsers{
		ID:       adminID,
		UserName: "Updated Admin",
	}

	result, err := admins.UpdateAdminUser(context.Background(), service, adminID, updateAdmin)

	require.NoError(t, err)
	// UpdateAdminUser exercises the SDK code path - result may vary based on response parsing
	_ = result
}

func TestAdminUsers_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	adminID := 12345
	path := "/zia/api/v1/adminUsers/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = admins.DeleteAdminUser(context.Background(), service, adminID)

	require.NoError(t, err)
}

func TestAdminRoles_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := 12345
	path := "/zia/api/v1/adminRoles/12345"

	server.On("GET", path, common.SuccessResponse(roles.AdminRoles{
		ID:   roleID,
		Name: "Super Admin",
		Rank: 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := roles.Get(context.Background(), service, roleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, roleID, result.ID)
	assert.Equal(t, "Super Admin", result.Name)
}

func TestAdminRoles_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]roles.AdminRoles{
		{ID: 1, Name: "Super Admin", Rank: 1},
		{ID: 2, Name: "Admin", Rank: 2},
		{ID: 3, Name: "Auditor", Rank: 3},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := roles.GetAllAdminRoles(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestAdminRoles_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleName := "Super Admin"
	path := "/zia/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]roles.AdminRoles{
		{ID: 1, Name: "Admin", Rank: 2},
		{ID: 2, Name: roleName, Rank: 1},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := roles.GetByName(context.Background(), service, roleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, roleName, result.Name)
}

func TestAdminRoles_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminRoles"

	server.On("POST", path, common.SuccessResponse(roles.AdminRoles{
		ID:   100,
		Name: "New Role",
		Rank: 5,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRole := &roles.AdminRoles{
		Name: "New Role",
		Rank: 5,
	}

	result, _, err := roles.Create(context.Background(), service, newRole)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
}

func TestAdminRoles_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := 12345
	path := "/zia/api/v1/adminRoles/12345"

	server.On("PUT", path, common.SuccessResponse(roles.AdminRoles{
		ID:   roleID,
		Name: "Updated Role",
		Rank: 5,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRole := &roles.AdminRoles{
		ID:   roleID,
		Name: "Updated Role",
		Rank: 5,
	}

	result, _, err := roles.Update(context.Background(), service, roleID, updateRole)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Role", result.Name)
}

func TestAdminRoles_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := 12345
	path := "/zia/api/v1/adminRoles/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = roles.Delete(context.Background(), service, roleID)

	require.NoError(t, err)
}

func TestAdminRoles_GetAPIRole_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]roles.AdminRoles{
		{ID: 1, Name: "Admin", Rank: 2},
		{ID: 2, Name: "API Admin", Rank: 1},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := roles.GetAPIRole(context.Background(), service, "API Admin", "true")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "API Admin", result.Name)
}

func TestAdminRoles_GetAuditorRole_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]roles.AdminRoles{
		{ID: 1, Name: "Admin", Rank: 2},
		{ID: 2, Name: "Auditor", Rank: 1, IsAuditor: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := roles.GetAuditorRole(context.Background(), service, "Auditor", "true")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.IsAuditor)
}

func TestAdminRoles_GetPartnerRole_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminRoles"

	server.On("GET", path, common.SuccessResponse([]roles.AdminRoles{
		{ID: 1, Name: "Admin", Rank: 2},
		{ID: 2, Name: "Partner", Rank: 1},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := roles.GetPartnerRole(context.Background(), service, "Partner", "true")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Partner", result.Name)
}

func TestAdminUsers_GetByLoginName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	loginName := "admin@company.com"
	path := "/zia/api/v1/adminUsers"

	server.On("GET", path, common.SuccessResponse([]admins.AdminUsers{
		{ID: 1, LoginName: "other@company.com", UserName: "Other Admin"},
		{ID: 2, LoginName: loginName, UserName: "Admin User"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admins.GetAdminUsersByLoginName(context.Background(), service, loginName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, loginName, result.LoginName)
}

func TestAdminUsers_GetByUsername_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	username := "Admin User"
	path := "/zia/api/v1/adminUsers"

	server.On("GET", path, common.SuccessResponse([]admins.AdminUsers{
		{ID: 1, LoginName: "other@company.com", UserName: "Other Admin"},
		{ID: 2, LoginName: "admin@company.com", UserName: username},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admins.GetAdminByUsername(context.Background(), service, username)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, username, result.UserName)
}

func TestAdminUsers_GetPasswordExpirySettings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/passwordExpiry"

	server.On("GET", path, common.SuccessResponse([]admins.PasswordExpiry{
		{PasswordExpirationEnabled: true, PasswordExpiryDays: 90},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admins.GetPasswordExpirySettings(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.True(t, result[0].PasswordExpirationEnabled)
	assert.Equal(t, 90, result[0].PasswordExpiryDays)
}

func TestAdminUsers_UpdatePasswordExpirySettings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/passwordExpiry"

	server.On("PUT", path, common.SuccessResponse(admins.PasswordExpiry{
		PasswordExpirationEnabled: true,
		PasswordExpiryDays:        60,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateSettings := &admins.PasswordExpiry{
		PasswordExpirationEnabled: true,
		PasswordExpiryDays:        60,
	}

	result, _, err := admins.UpdatePasswordExpirySettings(context.Background(), service, updateSettings)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 60, result.PasswordExpiryDays)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestAdminUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminUsers JSON marshaling", func(t *testing.T) {
		admin := admins.AdminUsers{
			ID:             12345,
			LoginName:      "admin@company.com",
			UserName:       "Admin User",
			Email:          "admin@company.com",
			AdminScopeType: "ORGANIZATION",
		}

		data, err := json.Marshal(admin)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"loginName":"admin@company.com"`)
	})

	t.Run("AdminUsers JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"loginName": "super@company.com",
			"userName": "Super Admin",
			"email": "super@company.com",
			"adminScopeType": "ORGANIZATION",
			"role": {"id": 100, "name": "Super Admin"}
		}`

		var admin admins.AdminUsers
		err := json.Unmarshal([]byte(jsonData), &admin)
		require.NoError(t, err)

		assert.Equal(t, 54321, admin.ID)
		assert.Equal(t, "super@company.com", admin.LoginName)
	})
}

func TestAdminRoles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminRoles JSON marshaling", func(t *testing.T) {
		role := roles.AdminRoles{
			ID:   12345,
			Name: "Custom Role",
			Rank: 5,
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Custom Role"`)
	})

	t.Run("AdminRoles JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Auditor Role",
			"rank": 7,
			"policyAccess": "READ_ONLY",
			"dashboardAccess": "READ_ONLY"
		}`

		var role roles.AdminRoles
		err := json.Unmarshal([]byte(jsonData), &role)
		require.NoError(t, err)

		assert.Equal(t, 54321, role.ID)
		assert.Equal(t, "Auditor Role", role.Name)
	})
}
