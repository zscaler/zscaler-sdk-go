package activation

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestZTWActivation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Test GetActivationStatus
	t.Run("GetActivationStatus", func(t *testing.T) {
		status, err := GetActivationStatus(context.Background(), service)
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Contains(t, []string{"EDITS_CLEARED", "EDITS_PRESENT", "EDITS_ACTIVATED_ON_RESTART"}, status.OrgEditStatus)
		assert.Contains(t, []string{"CAC_ACTV_UNKNOWN", "CAC_ACTV_UI", "CAC_ACTV_OLD_UI", "CAC_ACTV_SUPERADMIN", "CAC_ACTV_AUTOSYNC", "CAC_ACTV_TIMER"}, status.OrgLastActivateStatus)
	})

	// Test UpdateActivationStatus
	t.Run("UpdateActivationStatus", func(t *testing.T) {
		updateActivation := ECAdminActivation{}
		updatedStatus, err := UpdateActivationStatus(context.Background(), service, updateActivation)
		assert.NoError(t, err)
		assert.NotNil(t, updatedStatus)
		assert.Contains(t, []string{"ADM_LOGGED_IN", "ADM_EDITING", "ADM_ACTV_QUEUED", "ADM_ACTIVATING", "ADM_ACTV_DONE", "ADM_ACTV_FAIL", "ADM_EXPIRED"}, updatedStatus.AdminActivateStatus)
	})

	// Test ForceActivationStatus
	t.Run("ForceActivationStatus", func(t *testing.T) {
		forceActivation := ECAdminActivation{}
		forcedStatus, err := ForceActivationStatus(context.Background(), service, forceActivation)
		assert.NoError(t, err)
		assert.NotNil(t, forcedStatus)
		assert.Contains(t, []string{"ADM_LOGGED_IN", "ADM_EDITING", "ADM_ACTV_QUEUED", "ADM_ACTIVATING", "ADM_ACTV_DONE", "ADM_ACTV_FAIL", "ADM_EXPIRED"}, forcedStatus.AdminActivateStatus)
	})
}
