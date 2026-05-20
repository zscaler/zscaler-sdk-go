// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	commontests "github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
	nt "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/notification_template"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestNotificationTemplate_Get_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/notification-templates/{id}"
	server.On("GET", itemPath, commontests.SuccessResponse(nt.NotificationTemplate{
		ID:                  21,
		Name:                "Default",
		IsDefaultTemplate:   true,
		EnableClient:        true,
		EnableZia:           true,
		EnableServiceStatus: true,
		DurationInSeconds:   30,
		ZIANotificationTemplate: nt.ZIANotificationTemplate{
			EnableZiaFirewall: true,
			EnableZiaDNS:      true,
			EnableZiaIPS:      true,
		},
		ZPANotificationTemplate: nt.ZPANotificationTemplate{
			EnableZpaReauth:            true,
			ZpaReauthIntervalInMinutes: 10,
		},
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	got, err := nt.Get(context.Background(), service, 21)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 21, got.ID)
	assert.Equal(t, "Default", got.Name)
	assert.True(t, got.IsDefaultTemplate)
	assert.True(t, got.ZIANotificationTemplate.EnableZiaFirewall)
	assert.Equal(t, 10, got.ZPANotificationTemplate.ZpaReauthIntervalInMinutes)
}

func TestNotificationTemplate_Create_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/notification-templates"
	server.On("POST", listPath, commontests.SuccessResponse(nt.NotificationTemplate{
		ID:           77,
		Name:         "New Template",
		EnableClient: true,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	created, _, err := nt.Create(context.Background(), service, &nt.NotificationTemplate{
		Name:         "New Template",
		EnableClient: true,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, 77, created.ID)
	assert.Equal(t, "New Template", created.Name)
	assert.True(t, created.EnableClient)
}

func TestNotificationTemplate_Create_NilBody_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = nt.Create(context.Background(), service, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "notification template is required")
}

// PUT echoes the template body but omits `id`. The SDK pins the path id
// onto the result so callers don't see ID == 0.
func TestNotificationTemplate_Update_PinsPathIDWhenEcho_OmitsID_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/notification-templates/{id}"
	server.On("PUT", itemPath, commontests.SuccessResponse(nt.NotificationTemplate{
		Name:                "Updated",
		EnableClient:        true,
		EnableServiceStatus: true,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updated, _, err := nt.Update(context.Background(), service, 21, &nt.NotificationTemplate{
		Name:                "Updated",
		EnableClient:        true,
		EnableServiceStatus: true,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, 21, updated.ID, "SDK must pin the path ID when the API response omits id")
	assert.Equal(t, "Updated", updated.Name)
}

// PUT honors the id when the API does echo it.
func TestNotificationTemplate_Update_RespectsEchoedID_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/notification-templates/{id}"
	server.On("PUT", itemPath, commontests.SuccessResponse(nt.NotificationTemplate{
		ID:   21,
		Name: "Updated",
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updated, _, err := nt.Update(context.Background(), service, 21, &nt.NotificationTemplate{Name: "Updated"})
	require.NoError(t, err)
	assert.Equal(t, 21, updated.ID)
}

func TestNotificationTemplate_PartialUpdate_PinsPathID_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/notification-templates/{id}"
	server.On("PATCH", itemPath, commontests.SuccessResponse(nt.NotificationTemplate{
		Name: "Patched",
		ZPANotificationTemplate: nt.ZPANotificationTemplate{
			EnableZpaReauth:            true,
			ZpaReauthIntervalInMinutes: 15,
		},
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	patched, _, err := nt.PartialUpdate(context.Background(), service, 21, &nt.NotificationTemplate{
		Name: "Patched",
		ZPANotificationTemplate: nt.ZPANotificationTemplate{
			EnableZpaReauth:            true,
			ZpaReauthIntervalInMinutes: 15,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 21, patched.ID, "SDK must pin path ID for PATCH too")
	assert.Equal(t, "Patched", patched.Name)
	assert.Equal(t, 15, patched.ZPANotificationTemplate.ZpaReauthIntervalInMinutes)
}

func TestNotificationTemplate_PartialUpdate_NilBody_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = nt.PartialUpdate(context.Background(), service, 21, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "notification template is required")
}

func TestNotificationTemplate_Delete_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/notification-templates/{id}"
	server.On("DELETE", itemPath, commontests.NoContentResponse())

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = nt.Delete(context.Background(), service, 21)
	require.NoError(t, err)
}

func TestNotificationTemplate_GetAll_SinglePage_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/notification-templates"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[nt.NotificationTemplate]{
		Items: []nt.NotificationTemplate{
			{ID: 1, Name: "alpha"},
			{ID: 2, Name: "beta"},
		},
		Total: 2, Offset: 0, Limit: 50, Count: 2,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	all, err := nt.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	require.Len(t, all, 2)
	assert.Equal(t, "alpha", all[0].Name)
	assert.Equal(t, "beta", all[1].Name)
}

// Regression guard: the v2 endpoint takes `keyword` and `perPage`, never
// `search` or `pageSize`.
func TestNotificationTemplate_GetAll_UsesDocumentedQueryParams_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/notification-templates"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[nt.NotificationTemplate]{
		Items: []nt.NotificationTemplate{{ID: 1, Name: "alpha"}},
		Total: 1, Offset: 0, Limit: 50, Count: 1,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = nt.GetAll(context.Background(), service, &nt.GetAllFilterOptions{Keyword: "alpha"})
	require.NoError(t, err)

	last := server.LastRequest()
	require.NotNil(t, last)
	assert.Contains(t, last.Query, "perPage=50")
	assert.Contains(t, last.Query, "keyword=alpha")
	assert.NotContains(t, last.Query, "pageSize=")
	assert.NotContains(t, last.Query, "search=")
	// notification-templates has no documented `type` filter; the SDK
	// must not accidentally include it.
	assert.NotContains(t, last.Query, "type=")
	assert.NotContains(t, last.Query, "platformType=")
}

func TestNotificationTemplate_GetAll_MultiPage_SDK(t *testing.T) {
	var calls int32
	queries := make([]string, 0, 2)

	first := common.PaginatedResponseV2[nt.NotificationTemplate]{
		Items:  buildNTItems(1, 50),
		Total:  53,
		Offset: 0,
		Limit:  50,
		Count:  50,
	}
	second := common.PaginatedResponseV2[nt.NotificationTemplate]{
		Items:  buildNTItems(51, 3),
		Total:  53,
		Offset: 50,
		Limit:  50,
		Count:  3,
	}

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if atomic.AddInt32(&calls, 1) == 1 {
			_ = json.NewEncoder(w).Encode(first)
			return
		}
		_ = json.NewEncoder(w).Encode(second)
	}))
	defer upstream.Close()

	service, err := commontests.CreateTestService(context.Background(),
		&commontests.TestServer{Server: upstream, Handler: commontests.NewMockHandler()},
		"123456")
	require.NoError(t, err)

	all, err := nt.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	assert.Len(t, all, 53)

	require.Len(t, queries, 2)
	assert.NotContains(t, queries[0], "skip=")
	assert.Contains(t, queries[1], "skip=50")
}

func TestNotificationTemplate_GetByName_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/notification-templates"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[nt.NotificationTemplate]{
		Items: []nt.NotificationTemplate{
			{ID: 10, Name: "Default"},
			{ID: 11, Name: "Custom"},
		},
		Total: 2, Offset: 0, Limit: 50, Count: 2,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	got, err := nt.GetByName(context.Background(), service, "default")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 10, got.ID)
	assert.Equal(t, "Default", got.Name)

	last := server.LastRequest()
	require.NotNil(t, last)
	assert.Contains(t, last.Query, "keyword=default")
}

func TestNotificationTemplate_GetByName_NotFound_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/notification-templates"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[nt.NotificationTemplate]{
		Items: []nt.NotificationTemplate{{ID: 1, Name: "alpha"}},
		Total: 1, Offset: 0, Limit: 50, Count: 1,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = nt.GetByName(context.Background(), service, "missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no notification template found with name: missing")
}

// =====================================================
// Structure tests
// =====================================================

func TestNotificationTemplate_Structure(t *testing.T) {
	t.Parallel()

	// All `enable*` booleans must be serialized even when false — the API
	// rejects a payload that omits them entirely. This locks the
	// no-omitempty contract in.
	t.Run("Booleans always present on the wire even when false", func(t *testing.T) {
		template := nt.NotificationTemplate{
			Name: "Strict",
			ZIANotificationTemplate: nt.ZIANotificationTemplate{},
			ZPANotificationTemplate: nt.ZPANotificationTemplate{},
		}

		data, err := json.Marshal(template)
		require.NoError(t, err)
		body := string(data)

		assert.Contains(t, body, `"isDefaultTemplate":false`)
		assert.Contains(t, body, `"enableClient":false`)
		assert.Contains(t, body, `"enableZia":false`)
		assert.Contains(t, body, `"enableAppUpdates":false`)
		assert.Contains(t, body, `"enableServiceStatus":false`)
		assert.Contains(t, body, `"enablePersistent":false`)
		assert.Contains(t, body, `"enableDoNotDisturb":false`)

		assert.Contains(t, body, `"enableZiaFirewall":false`)
		assert.Contains(t, body, `"enableZiaDNS":false`)
		assert.Contains(t, body, `"enableZiaIPS":false`)
		assert.Contains(t, body, `"enableZiaPersistent":false`)

		assert.Contains(t, body, `"enableDevicePostureFailure":false`)
		assert.Contains(t, body, `"enableZpaReauth":false`)
		assert.Contains(t, body, `"delayPostureFailureSeconds":0`)
	})

	t.Run("NotificationTemplate round-trip preserves nested structs", func(t *testing.T) {
		jsonData := `{
			"id": 42,
			"name": "Mixed",
			"isDefaultTemplate": false,
			"enableClient": true,
			"enableZia": true,
			"enableAppUpdates": false,
			"enableServiceStatus": true,
			"durationInSeconds": 30,
			"enablePersistent": false,
			"enableDoNotDisturb": false,
			"ziaNotificationTemplate": {
				"enableZiaFirewall": true,
				"enableZiaFirewallPopup": false,
				"enableZiaDNS": true,
				"enableZiaDNSPopup": false,
				"enableZiaIPS": false,
				"enableZiaIPSPopup": false,
				"enableZiaPersistent": false
			},
			"zpaNotificationTemplate": {
				"enableDevicePostureFailure": true,
				"enableZpaReauth": true,
				"zpaReauthIntervalInMinutes": 15,
				"delayPostureFailureSeconds": 0
			}
		}`

		var got nt.NotificationTemplate
		err := json.Unmarshal([]byte(jsonData), &got)
		require.NoError(t, err)
		assert.Equal(t, 42, got.ID)
		assert.True(t, got.EnableClient)
		assert.True(t, got.ZIANotificationTemplate.EnableZiaDNS)
		assert.False(t, got.ZIANotificationTemplate.EnableZiaIPS)
		assert.True(t, got.ZPANotificationTemplate.EnableZpaReauth)
		assert.Equal(t, 15, got.ZPANotificationTemplate.ZpaReauthIntervalInMinutes)
	})
}

// =====================================================
// Helpers
// =====================================================

func buildNTItems(startID, n int) []nt.NotificationTemplate {
	items := make([]nt.NotificationTemplate, n)
	for i := 0; i < n; i++ {
		id := startID + i
		items[i] = nt.NotificationTemplate{ID: id}
	}
	return items
}
