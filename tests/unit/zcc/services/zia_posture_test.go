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
	zp "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/zia_posture"
)

// crowdstrikePosture mirrors the fixture in
// openapi/json_payload/zia_posture.json so the test data stays close to a
// real API payload.
func crowdstrikePosture() zp.ZIAPosture {
	cn := []zp.TrustCriterion{
		{ID: "9911", Name: "CrowdStrike_ZPA_ZTA_40", UDID: "6e36dd2f-ce19-47b3-8f26-f1e4e8f6313e"},
		{ID: "9913", Name: "CrowdStrike_ZPA_ZTA_80", UDID: "fc73ffb2-3ad7-49d5-9bff-10480589d188"},
		{ID: "9915", Name: "CrowdStrike_ZPA_Pre-ZTA", UDID: "cfab2ee9-9bf4-4482-9dcc-dadf7311c49b"},
	}
	cs := []zp.TrustCriteriaSet{{Cn: cn}}
	return zp.ZIAPosture{
		ID:                  16345,
		Name:                "CrowdstrikePosture01",
		Platform:            3,
		HighTrustCriteria:   zp.HighTrustCriteria{Cs: cs},
		MediumTrustCriteria: zp.MediumTrustCriteria{Cs: []zp.TrustCriteriaSet{{Cn: cn[:2]}}},
		LowTrustCriteria:    zp.LowTrustCriteria{Cs: []zp.TrustCriteriaSet{{Cn: cn[:2]}}},
	}
}

// =====================================================
// SDK Function Tests
// =====================================================

func TestZIAPosture_Get_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/zia-posture-profiles/{id}"
	server.On("GET", itemPath, commontests.SuccessResponse(crowdstrikePosture()))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	got, err := zp.Get(context.Background(), service, 16345)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 16345, got.ID)
	assert.Equal(t, "CrowdstrikePosture01", got.Name)
	assert.Equal(t, 3, got.Platform)
	require.Len(t, got.HighTrustCriteria.Cs, 1)
	require.Len(t, got.HighTrustCriteria.Cs[0].Cn, 3)
	assert.Equal(t, "9911", got.HighTrustCriteria.Cs[0].Cn[0].ID)
	assert.Equal(t, "CrowdStrike_ZPA_ZTA_40", got.HighTrustCriteria.Cs[0].Cn[0].Name)
}

func TestZIAPosture_Create_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/zia-posture-profiles"
	server.On("POST", listPath, commontests.SuccessResponse(crowdstrikePosture()))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	in := crowdstrikePosture()
	in.ID = 0
	created, _, err := zp.Create(context.Background(), service, &in)
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, 16345, created.ID)
	assert.Equal(t, "CrowdstrikePosture01", created.Name)
}

func TestZIAPosture_Create_NilBody_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = zp.Create(context.Background(), service, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zia posture is required")
}

// PUT echoes the payload but the API does not include `id`. SDK must
// pin the path id onto the returned object.
func TestZIAPosture_Update_PinsPathIDWhenAPIOmitsID_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/zia-posture-profiles/{id}"
	posture := crowdstrikePosture()
	posture.ID = 0
	server.On("PUT", itemPath, commontests.SuccessResponse(posture))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updated, _, err := zp.Update(context.Background(), service, 16345, &posture)
	require.NoError(t, err)
	assert.Equal(t, 16345, updated.ID, "SDK must pin the path id when the API response omits id")
	assert.Equal(t, "CrowdstrikePosture01", updated.Name)
}

func TestZIAPosture_PartialUpdate_PinsPathID_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/zia-posture-profiles/{id}"
	posture := crowdstrikePosture()
	posture.ID = 0
	server.On("PATCH", itemPath, commontests.SuccessResponse(posture))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	patched, _, err := zp.PartialUpdate(context.Background(), service, 16345, &posture)
	require.NoError(t, err)
	assert.Equal(t, 16345, patched.ID)
}

func TestZIAPosture_PartialUpdate_NilBody_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = zp.PartialUpdate(context.Background(), service, 16345, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zia posture is required")
}

func TestZIAPosture_Delete_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/zia-posture-profiles/{id}"
	server.On("DELETE", itemPath, commontests.NoContentResponse())

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = zp.Delete(context.Background(), service, 16345)
	require.NoError(t, err)
}

func TestZIAPosture_GetAll_SinglePage_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/zia-posture-profiles"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[zp.ZIAPosture]{
		Items: []zp.ZIAPosture{
			{ID: 1, Name: "alpha", Platform: 3},
			{ID: 2, Name: "beta", Platform: 4},
		},
		Total: 2, Offset: 0, Limit: 50, Count: 2,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	all, err := zp.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	require.Len(t, all, 2)
	assert.Equal(t, "alpha", all[0].Name)
}

// Regression guard: zia-posture-profiles takes `keyword`, `platformType`,
// `skip`, and `perPage`. v1 names must never leak.
func TestZIAPosture_GetAll_UsesDocumentedQueryParams_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/zia-posture-profiles"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[zp.ZIAPosture]{
		Items: []zp.ZIAPosture{{ID: 1, Name: "alpha", Platform: 3}},
		Total: 1, Offset: 0, Limit: 50, Count: 1,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = zp.GetAll(context.Background(), service, &zp.GetAllFilterOptions{
		Keyword:      "alpha",
		PlatformType: 3,
	})
	require.NoError(t, err)

	last := server.LastRequest()
	require.NotNil(t, last)
	assert.Contains(t, last.Query, "perPage=50")
	assert.Contains(t, last.Query, "keyword=alpha")
	assert.Contains(t, last.Query, "platformType=3")
	assert.NotContains(t, last.Query, "pageSize=")
	assert.NotContains(t, last.Query, "page=")
	assert.NotContains(t, last.Query, "type=")
}

// PlatformType=0 means "all platforms" per the docs and must be omitted
// from the wire (omitempty) — not sent as `platformType=0`.
func TestZIAPosture_GetAll_OmitsZeroPlatformType_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/zia-posture-profiles"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[zp.ZIAPosture]{
		Items: []zp.ZIAPosture{}, Total: 0, Offset: 0, Limit: 50, Count: 0,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = zp.GetAll(context.Background(), service, &zp.GetAllFilterOptions{PlatformType: 0})
	require.NoError(t, err)

	last := server.LastRequest()
	require.NotNil(t, last)
	assert.NotContains(t, last.Query, "platformType=")
}

func TestZIAPosture_GetAll_MultiPage_SDK(t *testing.T) {
	var calls int32
	queries := make([]string, 0, 2)

	first := common.PaginatedResponseV2[zp.ZIAPosture]{
		Items:  buildZPItems(1, 50),
		Total:  55,
		Offset: 0,
		Limit:  50,
		Count:  50,
	}
	second := common.PaginatedResponseV2[zp.ZIAPosture]{
		Items:  buildZPItems(51, 5),
		Total:  55,
		Offset: 50,
		Limit:  50,
		Count:  5,
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

	all, err := zp.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	assert.Len(t, all, 55)

	require.Len(t, queries, 2)
	assert.NotContains(t, queries[0], "skip=")
	assert.Contains(t, queries[1], "skip=50")
}

func TestZIAPosture_GetByName_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/zia-posture-profiles"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[zp.ZIAPosture]{
		Items: []zp.ZIAPosture{
			{ID: 16345, Name: "CrowdstrikePosture01"},
			{ID: 16346, Name: "CrowdstrikePosture02"},
		},
		Total: 2, Offset: 0, Limit: 50, Count: 2,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	got, err := zp.GetByName(context.Background(), service, "crowdstrikeposture01")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 16345, got.ID)

	last := server.LastRequest()
	require.NotNil(t, last)
	assert.Contains(t, last.Query, "keyword=crowdstrikeposture01")
}

func TestZIAPosture_GetByName_NotFound_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/zia-posture-profiles"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[zp.ZIAPosture]{
		Items: []zp.ZIAPosture{{ID: 1, Name: "alpha"}},
		Total: 1, Offset: 0, Limit: 50, Count: 1,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = zp.GetByName(context.Background(), service, "missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no zia posture found with name: missing")
}

// =====================================================
// Structure tests
// =====================================================

func TestZIAPosture_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ZIAPosture round-trips the documented payload shape", func(t *testing.T) {
		const payload = `{
			"id": 16345,
			"name": "CrowdstrikePosture01",
			"platform": 3,
			"highTrustCriteria": {
				"cs": [
					{
						"cn": [
							{"id": "9911", "name": "CrowdStrike_ZPA_ZTA_40", "udid": "6e36dd2f-ce19-47b3-8f26-f1e4e8f6313e"},
							{"id": "9913", "name": "CrowdStrike_ZPA_ZTA_80", "udid": "fc73ffb2-3ad7-49d5-9bff-10480589d188"},
							{"id": "9915", "name": "CrowdStrike_ZPA_Pre-ZTA", "udid": "cfab2ee9-9bf4-4482-9dcc-dadf7311c49b"}
						]
					}
				]
			},
			"mediumTrustCriteria": {
				"cs": [
					{
						"cn": [
							{"id": "9911", "name": "CrowdStrike_ZPA_ZTA_40", "udid": "6e36dd2f-ce19-47b3-8f26-f1e4e8f6313e"},
							{"id": "9913", "name": "CrowdStrike_ZPA_ZTA_80", "udid": "fc73ffb2-3ad7-49d5-9bff-10480589d188"}
						]
					}
				]
			},
			"lowTrustCriteria": {
				"cs": [
					{
						"cn": [
							{"id": "9911", "name": "CrowdStrike_ZPA_ZTA_40", "udid": "6e36dd2f-ce19-47b3-8f26-f1e4e8f6313e"},
							{"id": "9913", "name": "CrowdStrike_ZPA_ZTA_80", "udid": "fc73ffb2-3ad7-49d5-9bff-10480589d188"}
						]
					}
				]
			}
		}`

		var got zp.ZIAPosture
		err := json.Unmarshal([]byte(payload), &got)
		require.NoError(t, err)
		assert.Equal(t, 16345, got.ID)
		assert.Equal(t, "CrowdstrikePosture01", got.Name)
		assert.Equal(t, 3, got.Platform)
		require.Len(t, got.HighTrustCriteria.Cs, 1)
		require.Len(t, got.HighTrustCriteria.Cs[0].Cn, 3)
		require.Len(t, got.MediumTrustCriteria.Cs, 1)
		require.Len(t, got.MediumTrustCriteria.Cs[0].Cn, 2)
		require.Len(t, got.LowTrustCriteria.Cs, 1)
		require.Len(t, got.LowTrustCriteria.Cs[0].Cn, 2)
		assert.Equal(t, "9915", got.HighTrustCriteria.Cs[0].Cn[2].ID)
		assert.Equal(t, "cfab2ee9-9bf4-4482-9dcc-dadf7311c49b", got.HighTrustCriteria.Cs[0].Cn[2].UDID)
	})

	t.Run("ZIAPosture re-marshal preserves nested arrays", func(t *testing.T) {
		original := crowdstrikePosture()

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var round zp.ZIAPosture
		require.NoError(t, json.Unmarshal(data, &round))

		assert.Equal(t, original.ID, round.ID)
		assert.Equal(t, original.Name, round.Name)
		assert.Equal(t, original.Platform, round.Platform)
		assert.Equal(t, original.HighTrustCriteria, round.HighTrustCriteria)
		assert.Equal(t, original.MediumTrustCriteria, round.MediumTrustCriteria)
		assert.Equal(t, original.LowTrustCriteria, round.LowTrustCriteria)
	})
}

// =====================================================
// Helpers
// =====================================================

func buildZPItems(startID, n int) []zp.ZIAPosture {
	items := make([]zp.ZIAPosture, n)
	for i := 0; i < n; i++ {
		id := startID + i
		items[i] = zp.ZIAPosture{ID: id, Platform: 3}
	}
	return items
}
