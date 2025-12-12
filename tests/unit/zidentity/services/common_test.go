// Package services provides unit tests for ZIdentity services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/common"
)

func TestCommon_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IDNameDisplayName JSON marshaling", func(t *testing.T) {
		item := common.IDNameDisplayName{
			ID:          "item-123",
			Name:        "TestItem",
			DisplayName: "Test Item Display Name",
		}

		data, err := json.Marshal(item)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"item-123"`)
		assert.Contains(t, string(data), `"name":"TestItem"`)
		assert.Contains(t, string(data), `"displayName":"Test Item Display Name"`)
	})

	t.Run("IDNameDisplayName JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "item-456",
			"name": "AnotherItem",
			"displayName": "Another Item Display"
		}`

		var item common.IDNameDisplayName
		err := json.Unmarshal([]byte(jsonData), &item)
		require.NoError(t, err)

		assert.Equal(t, "item-456", item.ID)
		assert.Equal(t, "AnotherItem", item.Name)
		assert.Equal(t, "Another Item Display", item.DisplayName)
	})

	t.Run("PaginationResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"results_total": 150,
			"pageOffset": 50,
			"pageSize": 50,
			"next_link": "/api/v1/items?offset=100&limit=50",
			"prev_link": "/api/v1/items?offset=0&limit=50",
			"records": [
				{"id": "1", "name": "Item1"},
				{"id": "2", "name": "Item2"}
			]
		}`

		var response common.PaginationResponse[common.IDNameDisplayName]
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 150, response.ResultsTotal)
		assert.Equal(t, 50, response.PageOffset)
		assert.Equal(t, 50, response.PageSize)
		assert.NotEmpty(t, response.NextLink)
		assert.NotEmpty(t, response.PrevLink)
		assert.Len(t, response.Records, 2)
	})

	t.Run("PaginationQueryParams creation", func(t *testing.T) {
		params := common.NewPaginationQueryParams(50)
		
		assert.Equal(t, 50, params.Limit)
		assert.Equal(t, 0, params.Offset)
	})

	t.Run("PaginationQueryParams with default", func(t *testing.T) {
		params := common.NewPaginationQueryParams(0)
		
		assert.Equal(t, common.DefaultPaginationOptions.DefaultPageSize, params.Limit)
	})

	t.Run("PaginationQueryParams with max limit", func(t *testing.T) {
		params := common.NewPaginationQueryParams(5000)
		
		assert.Equal(t, common.DefaultPaginationOptions.MaxPageSize, params.Limit)
	})

	t.Run("PaginationQueryParams builder methods", func(t *testing.T) {
		params := common.NewPaginationQueryParams(100)
		params.WithOffset(200).
			WithLimit(50).
			WithNameFilter("test").
			WithExcludeDynamicGroups(true)
		
		assert.Equal(t, 200, params.Offset)
		assert.Equal(t, 50, params.Limit)
		assert.Equal(t, "test", params.NameLike)
		assert.True(t, params.ExcludeDynamicGroups)
	})

	t.Run("PaginationQueryParams user filter methods", func(t *testing.T) {
		params := common.NewPaginationQueryParams(100)
		params.WithLoginName([]string{"user1", "user2"}).
			WithLoginNameLike("john").
			WithDisplayNameLike("John Doe").
			WithPrimaryEmailLike("john@example.com").
			WithDomainName([]string{"example.com", "test.com"}).
			WithIDPName([]string{"Okta", "Azure"})
		
		assert.Equal(t, []string{"user1", "user2"}, params.LoginName)
		assert.Equal(t, "john", params.LoginNameLike)
		assert.Equal(t, "John Doe", params.DisplayNameLike)
		assert.Equal(t, "john@example.com", params.PrimaryEmailLike)
		assert.Equal(t, []string{"example.com", "test.com"}, params.DomainName)
		assert.Equal(t, []string{"Okta", "Azure"}, params.IDPName)
	})

	t.Run("PaginationQueryParams ToURLValues", func(t *testing.T) {
		params := common.NewPaginationQueryParams(100)
		params.WithOffset(50).
			WithNameFilter("test").
			WithExcludeDynamicGroups(true)
		
		values := params.ToURLValues()
		
		assert.Equal(t, "50", values.Get("offset"))
		assert.Equal(t, "100", values.Get("limit"))
		assert.Equal(t, "test", values.Get("name[like]"))
		assert.Equal(t, "true", values.Get("excludedynamicgroups"))
	})
}

func TestCommon_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse first page response", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 250,
			"pageOffset": 0,
			"pageSize": 100,
			"next_link": "/api/v1/items?offset=100&limit=100",
			"prev_link": "",
			"records": [
				{"id": "item-1", "name": "Item 1", "displayName": "First Item"}
			]
		}`

		var response common.PaginationResponse[common.IDNameDisplayName]
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.PageOffset)
		assert.NotEmpty(t, response.NextLink)
		assert.Empty(t, response.PrevLink)
	})

	t.Run("Parse last page response", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 250,
			"pageOffset": 200,
			"pageSize": 100,
			"next_link": "",
			"prev_link": "/api/v1/items?offset=100&limit=100",
			"records": [
				{"id": "item-201", "name": "Item 201"}
			]
		}`

		var response common.PaginationResponse[common.IDNameDisplayName]
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 200, response.PageOffset)
		assert.Empty(t, response.NextLink)
		assert.NotEmpty(t, response.PrevLink)
	})

	t.Run("Parse empty response", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 0,
			"pageOffset": 0,
			"pageSize": 100,
			"next_link": "",
			"prev_link": "",
			"records": []
		}`

		var response common.PaginationResponse[common.IDNameDisplayName]
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.ResultsTotal)
		assert.Empty(t, response.Records)
	})

	t.Run("DefaultPaginationOptions values", func(t *testing.T) {
		assert.Equal(t, 100, common.DefaultPaginationOptions.DefaultPageSize)
		assert.Equal(t, 1000, common.DefaultPaginationOptions.MaxPageSize)
		assert.False(t, common.DefaultPaginationOptions.UseCursor)
	})
}

