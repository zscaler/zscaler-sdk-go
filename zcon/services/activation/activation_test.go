package activation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestZCONActivation(t *testing.T) {
	// Assuming client is already set up correctly
	client, err := tests.NewZConClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	// Test GetActivationStatus
	t.Run("GetActivationStatus", func(t *testing.T) {
		status, err := service.GetActivationStatus()
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Contains(t, []string{"EDITS_CLEARED", "EDITS_PRESENT", "EDITS_ACTIVATED_ON_RESTART"}, status.OrgEditStatus)
		assert.Contains(t, []string{"CAC_ACTV_UNKNOWN", "CAC_ACTV_UI", "CAC_ACTV_OLD_UI", "CAC_ACTV_SUPERADMIN", "CAC_ACTV_AUTOSYNC", "CAC_ACTV_TIMER"}, status.OrgLastActivateStatus)
	})

	// Test UpdateActivationStatus
	t.Run("UpdateActivationStatus", func(t *testing.T) {
		updateActivation := ECAdminActivation{}
		updatedStatus, err := service.UpdateActivationStatus(updateActivation)
		assert.NoError(t, err)
		assert.NotNil(t, updatedStatus)
		assert.Contains(t, []string{"ADM_LOGGED_IN", "ADM_EDITING", "ADM_ACTV_QUEUED", "ADM_ACTIVATING", "ADM_ACTV_DONE", "ADM_ACTV_FAIL", "ADM_EXPIRED"}, updatedStatus.AdminActivateStatus)
	})

	// Test ForceActivationStatus
	t.Run("ForceActivationStatus", func(t *testing.T) {
		forceActivation := ECAdminActivation{}
		forcedStatus, err := service.ForceActivationStatus(forceActivation)
		assert.NoError(t, err)
		assert.NotNil(t, forcedStatus)
		assert.Contains(t, []string{"ADM_LOGGED_IN", "ADM_EDITING", "ADM_ACTV_QUEUED", "ADM_ACTIVATING", "ADM_ACTV_DONE", "ADM_ACTV_FAIL", "ADM_EXPIRED"}, forcedStatus.AdminActivateStatus)
	})
}
