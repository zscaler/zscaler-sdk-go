// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/activation"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestActivation_GetActivationStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/status"

	server.On("GET", path, common.SuccessResponse(activation.Activation{
		Status: "ACTIVE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := activation.GetActivationStatus(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "ACTIVE", result.Status)
}

func TestActivation_CreateActivation_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/status/activate"

	server.On("POST", path, common.SuccessResponse(activation.Activation{
		Status: "PENDING",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	activationReq := activation.Activation{
		Status: "ACTIVE",
	}

	result, err := activation.CreateActivation(context.Background(), service, activationReq)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "PENDING", result.Status)
}

func TestActivation_GetEusaStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/eusaStatus/latest"

	server.On("GET", path, common.SuccessResponse(activation.ZiaEusaStatus{
		ID:             12345,
		AcceptedStatus: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := activation.GetEusaStatus(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 12345, result.ID)
	assert.True(t, result.AcceptedStatus)
}

func TestActivation_UpdateEusaStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	statusID := 12345
	path := "/zia/api/v1/eusaStatus/12345"

	server.On("PUT", path, common.SuccessResponse(activation.ZiaEusaStatus{
		ID:             statusID,
		AcceptedStatus: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	eusaStatus := &activation.ZiaEusaStatus{
		ID:             statusID,
		AcceptedStatus: true,
	}

	result, _, err := activation.UpdateEusaStatus(context.Background(), service, statusID, eusaStatus)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.AcceptedStatus)
}

// =====================================================
// Structure Tests
// =====================================================

func TestActivation_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Activation JSON marshaling", func(t *testing.T) {
		act := activation.Activation{
			Status: "ACTIVE",
		}

		data, err := json.Marshal(act)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"status":"ACTIVE"`)
	})

	t.Run("Activation JSON unmarshaling", func(t *testing.T) {
		jsonData := `{"status":"PENDING"}`

		var act activation.Activation
		err := json.Unmarshal([]byte(jsonData), &act)
		require.NoError(t, err)

		assert.Equal(t, "PENDING", act.Status)
	})

	t.Run("ZiaEusaStatus JSON marshaling", func(t *testing.T) {
		eusa := activation.ZiaEusaStatus{
			ID:             12345,
			AcceptedStatus: true,
		}

		data, err := json.Marshal(eusa)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"acceptedStatus":true`)
	})

	t.Run("ZiaEusaStatus JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"version": {
				"id": 1,
				"name": "v1.0"
			},
			"acceptedStatus": false
		}`

		var eusa activation.ZiaEusaStatus
		err := json.Unmarshal([]byte(jsonData), &eusa)
		require.NoError(t, err)

		assert.Equal(t, 54321, eusa.ID)
		assert.False(t, eusa.AcceptedStatus)
		assert.NotNil(t, eusa.Version)
	})
}

func TestActivation_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse activation status response", func(t *testing.T) {
		jsonResponse := `{"status":"ACTIVE"}`

		var act activation.Activation
		err := json.Unmarshal([]byte(jsonResponse), &act)
		require.NoError(t, err)

		assert.Equal(t, "ACTIVE", act.Status)
	})

	t.Run("Parse various activation statuses", func(t *testing.T) {
		statuses := []string{"ACTIVE", "PENDING", "INPROGRESS", "NONE"}

		for _, status := range statuses {
			jsonData := `{"status":"` + status + `"}`
			var act activation.Activation
			err := json.Unmarshal([]byte(jsonData), &act)
			require.NoError(t, err)
			assert.Equal(t, status, act.Status)
		}
	})
}

