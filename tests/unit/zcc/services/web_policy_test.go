// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/web_policy"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestWebPolicy_GetListByCompanyID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/web/policy/listByCompany"

	server.On("GET", path, common.SuccessResponse([]map[string]interface{}{
		{"id": 1, "name": "Default Policy", "active": "true"},
		{"id": 2, "name": "Custom Policy", "active": "false"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_policy.GetPolicyListByCompanyID(context.Background(), service, nil, nil, nil, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestWebPolicy_Activate_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/web/policy/activate"

	server.On("PUT", path, common.SuccessResponse(web_policy.WebPolicyActivation{
		PolicyId:   123,
		DeviceType: 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	activation := &web_policy.WebPolicyActivation{
		PolicyId:   123,
		DeviceType: 1,
	}

	result, err := web_policy.ActivateWebPolicy(context.Background(), service, activation)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 123, result.PolicyId)
}

func TestWebPolicy_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/web/policy/edit"

	// The /edit endpoint returns a bare {success, id} envelope, with id as
	// an unquoted JSON number. The struct response models that exactly.
	server.On("PUT", path, common.SuccessResponse(map[string]interface{}{
		"success": "true",
		"id":      205241,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updatePolicy := &web_policy.WebPolicy{
		ID:     "123",
		Name:   "Updated Policy",
		Active: "true",
	}

	result, err := web_policy.UpdateWebPolicy(context.Background(), service, updatePolicy)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "true", result.Success)
	assert.Equal(t, "205241", result.ID.String())
}

func TestWebPolicy_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/web/policy/123/delete"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = web_policy.DeleteWebPolicy(context.Background(), service, 123)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestWebPolicy_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebPolicy JSON marshaling", func(t *testing.T) {
		policy := web_policy.WebPolicy{
			ID:          "123",
			Name:        "Enterprise Policy",
			Description: "Main enterprise web policy",
			Active:      "1",
			DeviceType:  3,
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"123"`)
		assert.Contains(t, string(data), `"name":"Enterprise Policy"`)
		assert.Contains(t, string(data), `"active":"1"`)
		assert.Contains(t, string(data), `"device_type":3`)
		// OS-specific blocks are pointer + omitempty: an empty payload
		// must not emit androidPolicy/iosPolicy/etc. or the API rejects
		// the request with a 400.
		assert.NotContains(t, string(data), `"androidPolicy"`)
		assert.NotContains(t, string(data), `"iosPolicy"`)
		assert.NotContains(t, string(data), `"linuxPolicy"`)
		assert.NotContains(t, string(data), `"macPolicy"`)
		assert.NotContains(t, string(data), `"windowsPolicy"`)
	})

	t.Run("WebPolicy JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "456",
			"name": "Branch Policy",
			"description": "Policy for branch offices",
			"active": "0",
			"device_type": 4
		}`

		var policy web_policy.WebPolicy
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.Equal(t, "456", policy.ID)
		assert.Equal(t, "Branch Policy", policy.Name)
		assert.Equal(t, "0", policy.Active)
		assert.Equal(t, 4, policy.DeviceType)
	})

	t.Run("WebPolicyActivation JSON marshaling", func(t *testing.T) {
		activation := web_policy.WebPolicyActivation{
			PolicyId:   123,
			DeviceType: 1,
		}

		data, err := json.Marshal(activation)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"policyId":123`)
		assert.Contains(t, string(data), `"deviceType":1`)
	})

	// The /edit endpoint silently rejects payloads that send numeric
	// scalars as quoted strings (the API responds with HTTP 200 +
	// {"success":"false","id":0}). Guard against accidental string-typed
	// fields by asserting the wire shape is always a JSON number.
	t.Run("WebPolicy numeric scalars marshal as JSON numbers", func(t *testing.T) {
		policy := web_policy.WebPolicy{
			Name:                      "Numeric Wire Format",
			Active:                    "1",
			DeviceType:                4,
			RuleOrder:                 1,
			LogMode:                   -1,
			LogLevel:                  0,
			LogFileSize:               100,
			ReactivateWebSecurityMins: 0,
			ReauthPeriod:              8,
			TunnelZappTraffic:         0,
			GroupAll:                  0,
			HighlightActiveControl:    0,
			SendDisableServiceReason:  0,
			EnableDeviceGroups:        0,
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)
		body := string(data)

		assert.Contains(t, body, `"ruleOrder":1`)
		assert.Contains(t, body, `"logMode":-1`)
		assert.Contains(t, body, `"logLevel":0`)
		assert.Contains(t, body, `"logFileSize":100`)
		assert.Contains(t, body, `"reactivateWebSecurityMinutes":0`)
		assert.Contains(t, body, `"reauth_period":8`)
		assert.Contains(t, body, `"tunnelZappTraffic":0`)
		assert.Contains(t, body, `"groupAll":0`)
		assert.Contains(t, body, `"highlightActiveControl":0`)
		assert.Contains(t, body, `"sendDisableServiceReason":0`)
		assert.Contains(t, body, `"enableDeviceGroups":0`)

		assert.NotContains(t, body, `"ruleOrder":"1"`)
		assert.NotContains(t, body, `"logMode":"-1"`)
		assert.NotContains(t, body, `"reauth_period":"8"`)
	})

	// The ZCC GET /listByCompany response is inconsistent: most of the
	// numeric fields come back as JSON numbers, but a couple of them
	// (reactivateWebSecurityMinutes, reauth_period) are returned quoted.
	// IntOrString must accept both forms, and also tolerate an empty
	// string for never-populated policies.
	t.Run("WebPolicy numeric scalars unmarshal from mixed wire forms", func(t *testing.T) {
		jsonData := `{
			"ruleOrder": 7,
			"logMode": -1,
			"logLevel": 0,
			"logFileSize": 100,
			"reactivateWebSecurityMinutes": "0",
			"reauth_period": "12",
			"tunnelZappTraffic": 0,
			"groupAll": "",
			"highlightActiveControl": 0,
			"sendDisableServiceReason": 0,
			"enableDeviceGroups": null
		}`

		var policy web_policy.WebPolicy
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.Equal(t, 7, int(policy.RuleOrder))
		assert.Equal(t, -1, int(policy.LogMode))
		assert.Equal(t, 0, int(policy.LogLevel))
		assert.Equal(t, 100, int(policy.LogFileSize))
		assert.Equal(t, 0, int(policy.ReactivateWebSecurityMins))
		assert.Equal(t, 12, int(policy.ReauthPeriod))
		assert.Equal(t, 0, int(policy.TunnelZappTraffic))
		assert.Equal(t, 0, int(policy.GroupAll))
		assert.Equal(t, 0, int(policy.HighlightActiveControl))
		assert.Equal(t, 0, int(policy.SendDisableServiceReason))
		assert.Equal(t, 0, int(policy.EnableDeviceGroups))
	})

	// The macPolicy block uses snake_case for the password and SSL cert
	// fields. Earlier we serialized them as disablePassword / installCerts
	// / logoutPassword / uninstallPassword (camelCase), which the parser
	// silently ignored, so the resulting policy looked empty on the wire.
	t.Run("MacPolicy snake_case wire shape", func(t *testing.T) {
		mac := web_policy.MacPolicy{
			DisablePassword:     "secret",
			InstallSslCerts:     1,
			LogoutPassword:      "lout",
			UninstallPassword:   "rm",
			BrowserAuthType:     -1,
			UseDefaultBrowser:   0,
			CaptivePortalConfig: `{"automaticCapture":1}`,
		}

		data, err := json.Marshal(mac)
		require.NoError(t, err)
		body := string(data)

		assert.Contains(t, body, `"disable_password":"secret"`)
		assert.Contains(t, body, `"install_ssl_certs":1`)
		assert.Contains(t, body, `"logout_password":"lout"`)
		assert.Contains(t, body, `"uninstall_password":"rm"`)
		assert.Contains(t, body, `"browserAuthType":-1`)
		assert.Contains(t, body, `"useDefaultBrowser":0`)
		assert.Contains(t, body, `"captivePortalConfig":"{\"automaticCapture\":1}"`)

		// Guard against regressing to the camelCase forms.
		assert.NotContains(t, body, `"disablePassword"`)
		assert.NotContains(t, body, `"installCerts"`)
		assert.NotContains(t, body, `"logoutPassword"`)
		assert.NotContains(t, body, `"uninstallPassword"`)
	})
}
