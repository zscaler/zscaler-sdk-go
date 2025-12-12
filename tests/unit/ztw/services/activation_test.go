// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/activation"
)

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

