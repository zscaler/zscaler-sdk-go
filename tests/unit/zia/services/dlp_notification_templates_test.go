package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_notification_templates"
)

const dlpNotificationTemplatesPath = "/zia/api/v1/dlpNotificationTemplates"

func TestDLPNotificationTemplates_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	templateID := 100
	server.On("GET", dlpNotificationTemplatesPath+"/100", common.SuccessResponse(dlp_notification_templates.DlpNotificationTemplates{
		ID:            templateID,
		Name:          "DLP Template Test",
		Subject:       "DLP Violation: ${TRANSACTION_ID} ${ENGINES}",
		AttachContent: true,
		TLSEnabled:    true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_notification_templates.Get(context.Background(), service, templateID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, templateID, result.ID)
	assert.Equal(t, "DLP Template Test", result.Name)
	assert.True(t, result.AttachContent)
	assert.True(t, result.TLSEnabled)
}

func TestDLPNotificationTemplates_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	templateName := "DLP Template Test"
	server.On("GET", dlpNotificationTemplatesPath, common.SuccessResponse([]dlp_notification_templates.DlpNotificationTemplates{
		{ID: 1, Name: "Other Template"},
		{ID: 100, Name: templateName, Subject: "DLP Violation: ${TRANSACTION_ID} ${ENGINES}"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_notification_templates.GetByName(context.Background(), service, templateName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, templateName, result.Name)
}

func TestDLPNotificationTemplates_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", dlpNotificationTemplatesPath, common.SuccessResponse(dlp_notification_templates.DlpNotificationTemplates{
		ID:            99999,
		Name:          "DLP Template Test",
		Subject:       "DLP Violation: ${TRANSACTION_ID} ${ENGINES}",
		AttachContent: true,
		TLSEnabled:    true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newTemplate := &dlp_notification_templates.DlpNotificationTemplates{
		Name:          "DLP Template Test",
		Subject:       "DLP Violation: ${TRANSACTION_ID} ${ENGINES}",
		AttachContent: true,
		TLSEnabled:    true,
		PlainTextMessage: "DLP violation detected.",
		HtmlMessage:      "<p>DLP violation detected.</p>",
	}

	result, _, err := dlp_notification_templates.Create(context.Background(), service, newTemplate)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestDLPNotificationTemplates_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	templateID := 100
	server.On("PUT", dlpNotificationTemplatesPath+"/100", common.SuccessResponse(dlp_notification_templates.DlpNotificationTemplates{
		ID:      templateID,
		Name:    "Updated DLP Template",
		Subject: "Updated Subject",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateTemplate := &dlp_notification_templates.DlpNotificationTemplates{
		ID:      templateID,
		Name:    "Updated DLP Template",
		Subject: "Updated Subject",
	}

	result, _, err := dlp_notification_templates.Update(context.Background(), service, templateID, updateTemplate)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated DLP Template", result.Name)
}

func TestDLPNotificationTemplates_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	templateID := 100
	server.On("DELETE", dlpNotificationTemplatesPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = dlp_notification_templates.Delete(context.Background(), service, templateID)

	require.NoError(t, err)
}

func TestDLPNotificationTemplates_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", dlpNotificationTemplatesPath, common.SuccessResponse([]dlp_notification_templates.DlpNotificationTemplates{
		{ID: 1, Name: "Template 1", AttachContent: true},
		{ID: 2, Name: "Template 2", TLSEnabled: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_notification_templates.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDLPNotificationTemplates_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		template := dlp_notification_templates.DlpNotificationTemplates{
			ID:            100,
			Name:          "DLP Template Test",
			Subject:       "DLP Violation: ${TRANSACTION_ID} ${ENGINES}",
			AttachContent: true,
			TLSEnabled:    true,
		}

		data, err := json.Marshal(template)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"attachContent":true`)
		assert.Contains(t, string(data), `"tlsEnabled":true`)
	})
}
