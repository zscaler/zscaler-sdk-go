// Package services provides unit tests for ZIdentity services
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcommon "github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	apiclients "github.com/zscaler/zscaler-sdk-go/v3/zscaler/ziam/services/api_clients"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ziam/services/common"
)

func TestAPIClients_Structure(t *testing.T) {
	t.Parallel()

	t.Run("APIClients JSON marshaling", func(t *testing.T) {
		status := true
		apiClient := apiclients.APIClients{
			ID:                  "client-12345",
			Name:                "terraform-automation",
			Description:         "Client used by CI",
			Status:              &status,
			AccessTokenLifeTime: 3600,
			ClientAuthentication: &apiclients.ClientAuthentication{
				AuthType: apiclients.AuthTypeSecret,
			},
			ClientResources: []apiclients.ClientResource{
				{
					ID: "jhlm44rd107q7",
					SelectedScopes: []apiclients.SelectedScope{
						{ID: "hhlm44rae07ib:mplm44rqi07jb:hplm44rqvg7n5"},
					},
				},
			},
		}

		data, err := json.Marshal(apiClient)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"client-12345"`)
		assert.Contains(t, string(data), `"name":"terraform-automation"`)
		assert.Contains(t, string(data), `"status":true`)
		assert.Contains(t, string(data), `"accessTokenLifeTime":3600`)
		assert.Contains(t, string(data), `"authType":"SECRET"`)
		assert.Contains(t, string(data), `"selectedScopes"`)
	})

	// Status is a *bool for the same reason users.Users.Status is: as a plain
	// bool with `omitempty`, a configured false was dropped from the request
	// body, so a client could be enabled but never disabled. These three cases
	// pin all the states.
	t.Run("Status pointer expresses disabled, enabled, and unset", func(t *testing.T) {
		disabled := false
		data, err := json.Marshal(apiclients.APIClients{Name: "c", Status: &disabled})
		require.NoError(t, err)
		assert.Contains(t, string(data), `"status":false`,
			"an explicit false must be transmitted, or a client can never be disabled")

		enabled := true
		data, err = json.Marshal(apiclients.APIClients{Name: "c", Status: &enabled})
		require.NoError(t, err)
		assert.Contains(t, string(data), `"status":true`)

		data, err = json.Marshal(apiclients.APIClients{Name: "c"})
		require.NoError(t, err)
		assert.NotContains(t, string(data), `"status"`,
			"an unset status must be omitted so callers that do not manage it are unaffected")
	})

	// accessTokenLifeTime is required and carries no `omitempty`, so a zero
	// value reaches the API and is rejected with a range error rather than
	// vanishing from the body and letting the server apply a default.
	t.Run("accessTokenLifeTime is always transmitted", func(t *testing.T) {
		data, err := json.Marshal(apiclients.APIClients{Name: "c"})
		require.NoError(t, err)
		assert.Contains(t, string(data), `"accessTokenLifeTime":0`)
	})

	// clientAuthentication is a pointer because `omitempty` does nothing on a
	// struct value: as a value type an omitted object would unmarshal to an
	// empty struct, indistinguishable from one whose fields are blank.
	t.Run("clientAuthentication distinguishes absent from empty", func(t *testing.T) {
		data, err := json.Marshal(apiclients.APIClients{Name: "c"})
		require.NoError(t, err)
		assert.NotContains(t, string(data), `"clientAuthentication"`)

		var parsed apiclients.APIClients
		require.NoError(t, json.Unmarshal([]byte(`{"name":"c"}`), &parsed))
		assert.Nil(t, parsed.ClientAuthentication)
	})

	t.Run("AuthType values", func(t *testing.T) {
		assert.Equal(t, apiclients.AuthType("SECRET"), apiclients.AuthTypeSecret)
		assert.Equal(t, apiclients.AuthType("PUBKEYCERT"), apiclients.AuthTypePubKeyCert)
		assert.Equal(t, apiclients.AuthType("JWKS"), apiclients.AuthTypeJWKS)
	})

	t.Run("AddSecretRequest transmits a zero expiry", func(t *testing.T) {
		data, err := json.Marshal(apiclients.AddSecretRequest{ExpiresAt: 0})
		require.NoError(t, err)
		assert.JSONEq(t, `{"expiresAt":0}`, string(data),
			"expiresAt is required, so `omitempty` would produce a body the endpoint rejects")
	})
}

func TestAPIClients_ResponseParsing(t *testing.T) {
	t.Parallel()

	// The payload here is written against the documented wire format rather
	// than produced by marshaling the struct, so the test cannot pass merely
	// because the struct round-trips itself.
	t.Run("Parse an api client as reported by the API", func(t *testing.T) {
		jsonResponse := `{
			"id": "client-12345",
			"name": "terraform-automation",
			"description": "Client used by CI",
			"status": true,
			"accessTokenLifeTime": 3600,
			"clientAuthentication": {
				"clientJWKsUrl": "https://example.com/.well-known/jwks.json",
				"publicKeys": [
					{"keyName": "primary", "keyValue": "MIIBIjANBg"}
				],
				"clientCertificates": [
					{"certContent": "-----BEGIN CERTIFICATE-----"}
				],
				"authType": "JWKS"
			},
			"clientResources": [
				{
					"id": "jhlm44rd107q7",
					"name": "Zscaler APIs",
					"defaultApi": true,
					"selectedScopes": [
						{"id": "hhlm44rae07ib:mplm44rqi07jb:hplm44rqvg7n5", "name": "Scope:Default: Role:API Full Access"},
						{"id": "hhlm44rd307qf::9h6p7ebv903k4", "name": "Role:Super Admin"}
					]
				}
			]
		}`

		var apiClient apiclients.APIClients
		require.NoError(t, json.Unmarshal([]byte(jsonResponse), &apiClient))

		assert.Equal(t, "client-12345", apiClient.ID)
		assert.Equal(t, "terraform-automation", apiClient.Name)
		assert.Equal(t, "Client used by CI", apiClient.Description)
		require.NotNil(t, apiClient.Status)
		assert.True(t, *apiClient.Status)
		assert.Equal(t, int32(3600), apiClient.AccessTokenLifeTime)

		require.NotNil(t, apiClient.ClientAuthentication)
		auth := apiClient.ClientAuthentication
		assert.Equal(t, apiclients.AuthTypeJWKS, auth.AuthType)
		assert.Equal(t, "https://example.com/.well-known/jwks.json", auth.ClientJWKsURL)
		require.Len(t, auth.PublicKeys, 1)
		assert.Equal(t, "primary", auth.PublicKeys[0].KeyName)
		assert.Equal(t, "MIIBIjANBg", auth.PublicKeys[0].KeyValue)
		require.Len(t, auth.ClientCertificates, 1)
		assert.Equal(t, "-----BEGIN CERTIFICATE-----", auth.ClientCertificates[0].CertContent)

		// name and defaultApi are reported on read but not accepted on write;
		// carrying them costs nothing and omitting them would drop data.
		require.Len(t, apiClient.ClientResources, 1)
		resource := apiClient.ClientResources[0]
		assert.Equal(t, "jhlm44rd107q7", resource.ID)
		assert.Equal(t, "Zscaler APIs", resource.Name)
		assert.True(t, resource.DefaultApi)
		require.Len(t, resource.SelectedScopes, 2)
		assert.Equal(t, "hhlm44rae07ib:mplm44rqi07jb:hplm44rqvg7n5", resource.SelectedScopes[0].ID)
		assert.Equal(t, "Scope:Default: Role:API Full Access", resource.SelectedScopes[0].Name)
	})

	t.Run("Parse a disabled api client", func(t *testing.T) {
		var apiClient apiclients.APIClients
		require.NoError(t, json.Unmarshal([]byte(`{"id":"c","status":false}`), &apiClient))

		require.NotNil(t, apiClient.Status, "a reported false must not be indistinguishable from absent")
		assert.False(t, *apiClient.Status)
	})

	t.Run("Parse the list envelope", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 2,
			"pageOffset": 0,
			"pageSize": 100,
			"next_link": "",
			"prev_link": "",
			"records": [
				{"id": "client-1", "name": "ci-automation", "accessTokenLifeTime": 3600},
				{"id": "client-2", "name": "reporting", "accessTokenLifeTime": 60}
			]
		}`

		var response common.PaginationResponse[apiclients.APIClients]
		require.NoError(t, json.Unmarshal([]byte(jsonResponse), &response))

		assert.Equal(t, 2, response.ResultsTotal)
		assert.Equal(t, 100, response.PageSize)
		require.Len(t, response.Records, 2)
		assert.Equal(t, "ci-automation", response.Records[0].Name)
		assert.Equal(t, int32(60), response.Records[1].AccessTokenLifeTime)
	})

	t.Run("Parse the secrets list", func(t *testing.T) {
		// A bare array rather than the paginated envelope, and the values are
		// masked.
		jsonResponse := `[
			{"expiresAt": 1767225600, "id": "secret-1", "createdAt": 1735689600, "value": "****"},
			{"expiresAt": 1798761600, "id": "secret-2", "createdAt": 1735689600, "value": "****"}
		]`

		var secrets []apiclients.APIClientSecret
		require.NoError(t, json.Unmarshal([]byte(jsonResponse), &secrets))

		require.Len(t, secrets, 2)
		assert.Equal(t, "secret-1", secrets[0].ID)
		assert.Equal(t, int64(1767225600), secrets[0].ExpiresAt)
		assert.Equal(t, int64(1735689600), secrets[0].CreatedAt)
		assert.Equal(t, "****", secrets[0].Value)
	})
}

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestAPIClients_Get_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID

	server.On("GET", path, testcommon.SuccessResponse(map[string]interface{}{
		"id":                  clientID,
		"name":                "terraform-automation",
		"status":              true,
		"accessTokenLifeTime": 3600,
		"clientAuthentication": map[string]interface{}{
			"authType": "SECRET",
		},
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := apiclients.Get(context.Background(), service, clientID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, clientID, result.ID)
	assert.Equal(t, "terraform-automation", result.Name)
	require.NotNil(t, result.Status)
	assert.True(t, *result.Status)
	require.NotNil(t, result.ClientAuthentication)
	assert.Equal(t, apiclients.AuthTypeSecret, result.ClientAuthentication.AuthType)
}

func TestAPIClients_GetAll_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	path := "/ziam/admin/api/v1/api-clients"

	server.On("GET", path, testcommon.SuccessResponse(common.PaginationResponse[apiclients.APIClients]{
		ResultsTotal: 2,
		PageOffset:   0,
		PageSize:     100,
		Records: []apiclients.APIClients{
			{ID: "client-1", Name: "ci-automation", AccessTokenLifeTime: 3600},
			{ID: "client-2", Name: "reporting", AccessTokenLifeTime: 60},
		},
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	results, err := apiclients.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, "ci-automation", results[0].Name)
}

// GetPage exposes the pagination metadata that GetAll discards.
func TestAPIClients_GetPage_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	path := "/ziam/admin/api/v1/api-clients"

	server.On("GET", path, testcommon.SuccessResponse(common.PaginationResponse[apiclients.APIClients]{
		ResultsTotal: 7,
		PageOffset:   0,
		PageSize:     2,
		NextLink:     "/ziam/admin/api/v1/api-clients?offset=2&limit=2",
		Records: []apiclients.APIClients{
			{ID: "client-1", Name: "ci-automation"},
			{ID: "client-2", Name: "reporting"},
		},
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	page, err := apiclients.GetPage(context.Background(), service, nil)
	require.NoError(t, err)
	require.NotNil(t, page)
	assert.Equal(t, 7, page.ResultsTotal)
	assert.NotEmpty(t, page.NextLink)
	assert.Len(t, page.Records, 2)
}

// GetByName is a substring match, so a search for "api" returns every client
// whose name contains it. Callers needing one client must narrow further.
func TestAPIClients_GetByName_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	path := "/ziam/admin/api/v1/api-clients"

	server.On("GET", path, testcommon.SuccessResponse(common.PaginationResponse[apiclients.APIClients]{
		ResultsTotal: 3,
		PageOffset:   0,
		PageSize:     100,
		Records: []apiclients.APIClients{
			{ID: "client-1", Name: "prod-api"},
			{ID: "client-2", Name: "reporting"},
			{ID: "client-3", Name: "prod-api-secondary"},
		},
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	results, err := apiclients.GetByName(context.Background(), service, "prod-api")
	require.NoError(t, err)
	require.Len(t, results, 2, "a substring search must return every candidate, not just the first")
	assert.Equal(t, "prod-api", results[0].Name)
	assert.Equal(t, "prod-api-secondary", results[1].Name)
}

func TestAPIClients_Create_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	path := "/ziam/admin/api/v1/api-clients"

	server.On("POST", path, testcommon.SuccessResponse(map[string]interface{}{
		"id":                  "client-12345",
		"name":                "terraform-automation",
		"accessTokenLifeTime": 3600,
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	status := true
	result, _, err := apiclients.Create(context.Background(), service, &apiclients.APIClients{
		Name:                "terraform-automation",
		Description:         "Client used by CI",
		Status:              &status,
		AccessTokenLifeTime: 3600,
		ClientAuthentication: &apiclients.ClientAuthentication{
			AuthType: apiclients.AuthTypeSecret,
		},
		ClientResources: []apiclients.ClientResource{
			{
				ID:             "jhlm44rd107q7",
				SelectedScopes: []apiclients.SelectedScope{{ID: "hhlm44rd307qf::9h6p7ebv903k4"}},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "client-12345", result.ID)

	req := server.LastRequest()
	require.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, path, req.Path)
	assert.JSONEq(t, `{
		"name": "terraform-automation",
		"description": "Client used by CI",
		"status": true,
		"accessTokenLifeTime": 3600,
		"clientAuthentication": {"authType": "SECRET"},
		"clientResources": [
			{"id": "jhlm44rd107q7", "selectedScopes": [{"id": "hhlm44rd307qf::9h6p7ebv903k4"}]}
		]
	}`, string(req.Body))
}

func TestAPIClients_Create_NilPayload_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := apiclients.Create(context.Background(), service, nil)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestAPIClients_Update_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID

	server.On("PUT", path, testcommon.SuccessResponse(map[string]interface{}{
		"id":                  clientID,
		"name":                "terraform-automation-renamed",
		"accessTokenLifeTime": 7200,
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := apiclients.Update(context.Background(), service, clientID, &apiclients.APIClients{
		ID:                  clientID,
		Name:                "terraform-automation-renamed",
		AccessTokenLifeTime: 7200,
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "terraform-automation-renamed", result.Name)

	req := server.LastRequest()
	require.NotNil(t, req)
	assert.Equal(t, "PUT", req.Method)
	assert.Equal(t, path, req.Path)
}

// Disabling a client must reach the wire as an explicit false.
func TestAPIClients_Update_DisableClient_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID

	server.On("PUT", path, testcommon.SuccessResponse(map[string]interface{}{"id": clientID}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	disabled := false
	_, _, err = apiclients.Update(context.Background(), service, clientID, &apiclients.APIClients{
		ID:                  clientID,
		Name:                "terraform-automation",
		Status:              &disabled,
		AccessTokenLifeTime: 3600,
	})
	require.NoError(t, err)

	req := server.LastRequest()
	require.NotNil(t, req)
	assert.Contains(t, string(req.Body), `"status":false`,
		"a client that cannot be disabled is the bug this pointer exists to prevent")
}

// A 204 answer to the PUT must not panic. The shared client returns a nil
// object for an empty body, so the response is nil-checked before use.
func TestAPIClients_Update_NoContentResponse_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID

	server.On("PUT", path, testcommon.SuccessResponseWithStatus(http.StatusNoContent, nil))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := apiclients.Update(context.Background(), service, clientID, &apiclients.APIClients{
		ID:                  clientID,
		Name:                "terraform-automation",
		AccessTokenLifeTime: 3600,
	})
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestAPIClients_Update_NilPayload_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := apiclients.Update(context.Background(), service, "client-12345", nil)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestAPIClients_Delete_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID

	server.On("DELETE", path, testcommon.NoContentResponse())

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = apiclients.Delete(context.Background(), service, clientID)
	require.NoError(t, err)

	req := server.LastRequest()
	require.NotNil(t, req)
	assert.Equal(t, "DELETE", req.Method)
	assert.Equal(t, path, req.Path)
}

// =====================================================
// Client secrets
// =====================================================

func TestAPIClients_GetSecrets_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID + "/secrets"

	// A bare array, not the paginated envelope the rest of the package uses.
	server.On("GET", path, testcommon.SuccessResponse([]map[string]interface{}{
		{"expiresAt": 1767225600, "id": "secret-1", "createdAt": 1735689600, "value": "****"},
		{"expiresAt": 1798761600, "id": "secret-2", "createdAt": 1735689600, "value": "****"},
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	secrets, err := apiclients.GetSecrets(context.Background(), service, clientID)
	require.NoError(t, err)
	require.Len(t, secrets, 2)
	assert.Equal(t, "secret-1", secrets[0].ID)
	assert.Equal(t, int64(1767225600), secrets[0].ExpiresAt)
}

// This is the load-bearing test of the file.
//
// AddSecret deliberately bypasses Client.Create, which unmarshals the response
// into a value of the *request* struct's type — AddSecretRequest, a struct
// holding nothing but ExpiresAt. Routed through that helper, the id, the
// creation time, and above all the secret value would be silently discarded,
// and the value is unrecoverable afterwards because GetSecrets masks it.
//
// If someone "simplifies" AddSecret onto Client.Create, this test fails.
func TestAPIClients_AddSecret_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID + "/secrets"

	server.On("POST", path, testcommon.SuccessResponseWithStatus(http.StatusCreated, map[string]interface{}{
		"expiresAt": 1767225600,
		"id":        "secret-1",
		"createdAt": 1735689600,
		"value":     "the-only-time-this-is-ever-returned",
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	secret, _, err := apiclients.AddSecret(context.Background(), service, clientID, 1767225600)
	require.NoError(t, err)
	require.NotNil(t, secret)

	assert.Equal(t, "secret-1", secret.ID)
	assert.Equal(t, "the-only-time-this-is-ever-returned", secret.Value,
		"the secret value must survive; it cannot be read back later")
	assert.Equal(t, int64(1735689600), secret.CreatedAt)
	assert.Equal(t, int64(1767225600), secret.ExpiresAt)

	req := server.LastRequest()
	require.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, path, req.Path)
	assert.JSONEq(t, `{"expiresAt":1767225600}`, string(req.Body))
}

// A zero expiry must still be transmitted: expiresAt is required, and
// `omitempty` would drop it and produce a body the endpoint rejects.
func TestAPIClients_AddSecret_ZeroExpiry_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID + "/secrets"

	server.On("POST", path, testcommon.SuccessResponseWithStatus(http.StatusCreated, map[string]interface{}{
		"id": "secret-1", "value": "v", "expiresAt": 0, "createdAt": 1735689600,
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = apiclients.AddSecret(context.Background(), service, clientID, 0)
	require.NoError(t, err)

	req := server.LastRequest()
	require.NotNil(t, req)
	assert.JSONEq(t, `{"expiresAt":0}`, string(req.Body))
}

// An empty body is success with nothing to report rather than an error, so a
// tenant answering 204 does not surface as a failed call.
func TestAPIClients_AddSecret_NoContentResponse_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	path := "/ziam/admin/api/v1/api-clients/" + clientID + "/secrets"

	server.On("POST", path, testcommon.SuccessResponseWithStatus(http.StatusNoContent, nil))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	secret, _, err := apiclients.AddSecret(context.Background(), service, clientID, 1767225600)
	require.NoError(t, err)
	assert.Nil(t, secret)
}

func TestAPIClients_DeleteSecret_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	clientID := "client-12345"
	secretID := "secret-1"
	path := "/ziam/admin/api/v1/api-clients/" + clientID + "/secrets/" + secretID

	server.On("DELETE", path, testcommon.NoContentResponse())

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = apiclients.DeleteSecret(context.Background(), service, clientID, secretID)
	require.NoError(t, err)

	req := server.LastRequest()
	require.NotNil(t, req)
	assert.Equal(t, "DELETE", req.Method)
	assert.Equal(t, path, req.Path)
}

// TestAPIClients_SelectedScopeIDIsOpaque pins that a scope id crosses the wire
// exactly as resource_servers reported it.
//
// The id looks decomposable — "<serviceId>:<zpaScopeId>:<id>", with the middle
// segment empty for every service but ZPA — and it is tempting to split it into
// the tserviceId / zpaScopeId fields the ZIdentity admin console sends. The
// public API rejects that with "400 Invalid Input for clientResources", and
// both accepts and returns the joined form. Verified against a tenant; do not
// reintroduce a split.
func TestAPIClients_SelectedScopeIDIsOpaque(t *testing.T) {
	t.Parallel()

	for _, id := range []string{
		"hhlm44rae07ib:mplm44rqi07jb:hplm44rqvg7n5", // ZPA
		"hhlm44raf07ps::hplm45bceg7r0",              // ZIA
		"hhlm44rd307qf::9h6p7ebv903k4",              // ZIAM
	} {
		data, err := json.Marshal(apiclients.SelectedScope{ID: id})
		require.NoError(t, err)

		assert.JSONEq(t, `{"id":"`+id+`"}`, string(data))
		assert.NotContains(t, string(data), "tserviceId")
		assert.NotContains(t, string(data), "zpaScopeId")
	}

	// The read side returns the same joined id plus a resolved name.
	var scope apiclients.SelectedScope
	require.NoError(t, json.Unmarshal(
		[]byte(`{"id":"hhlm44rae07ib:mplm44rqi07jb:hplm44rqvg7n5","name":"Scope:Default: Role:API Full Access"}`),
		&scope))
	assert.Equal(t, "hhlm44rae07ib:mplm44rqi07jb:hplm44rqvg7n5", scope.ID)
	assert.Equal(t, "Scope:Default: Role:API Full Access", scope.Name)
}
