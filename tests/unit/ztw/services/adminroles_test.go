// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/adminuserrolemgmt/adminroles"
)

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

